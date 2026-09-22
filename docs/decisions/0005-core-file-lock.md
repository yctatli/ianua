# ADR 0005 — Core-file-lock: protect the workspace's own process files during ordinary dev

- Status: Accepted
- Date: 2026-09-22

## Context
The human asked directly: while doing normal development (a feature, a bugfix), how can they be
sure this workspace's own governance files — `AGENTS.md`, `workflows/`, `prompts/`, `scripts/`,
`adapters/`, `docs/roles/`, `docs/decisions/`, spec/plan templates — never get overwritten as a
side effect? Before this ADR, the only existing technical protection was `specs/done/` immutability
(rule 7, enforced by `protect-shipped.sh` / `policy.py`). Everything else was prose: role cards say
Developer "MAY NOT... exceed the plan's file scope silently," `docs/decisions/README.md` says ADRs
aren't edited after acceptance — real rules, zero technical enforcement. `AGENTS.md`'s own founding
principle again: "Rules it cannot bypass... Prose is advice; tooling is law."

A real risk this closes: a builder-tier agent under pressure to reach green could plausibly "fix"
`scripts/check.conf` by removing a failing step, or touch `AGENTS.md`/`workflows/` as an incidental
side effect of an unrelated task, with nothing structurally stopping it.

The hard part: unlike `specs/done/` (truly permanent, no legitimate reason to ever change), these
process files are supposed to evolve — just not as a side effect of feature/bugfix BUILD work. They
change during **bootstrap** (`AGENTS.md`, `scripts/check.conf`) and during deliberate **meta/ADR
work** (exactly what this session has been doing for the last several turns — editing `workflows/`,
`prompts/`, `docs/roles/`, both adapters). A hook that blocks these paths unconditionally, with no
escape hatch, would have made this entire session's work impossible.

Also found and fixed along the way: the existing hook only matched Claude Code's `Edit`/`Write`
tools — a `Bash` command doing `sed -i`/redirect/`cp`/`rm` against a protected path sailed through
untouched. Same gap existed conceptually on the Codex side (it already checked Bash, this was
Claude-Code-specific).

## Decision
Extend the existing immutability hook (`protect-shipped.sh` for Claude Code, `policy.py` for
Codex) with a second, gated tier, keeping the two tiers clearly distinct:

1. **`specs/done/`** — unchanged: absolute, no override, ever.
2. **Core set** (`AGENTS.md`, `CLAUDE.md`, `workflows/`, `prompts/`, `scripts/check`,
   `scripts/check.conf`, `scripts/doctor`, `scripts/init`, `adapters/`, `docs/roles/`,
   `docs/decisions/`, `specs/TEMPLATE.md`, `specs/plans/TEMPLATE.md`, and the installed
   `.claude/`, `.codex/`, `.agents/` directories themselves) — blocked by default, allowed only
   when the human sets `ANEW_ALLOW_CORE_EDIT=1|true|yes` for that session. No silent bypass: the
   block message always states which file/path and why.

Both hooks now also match `Bash`-tool-issued writes (redirect/`sed -i`/`cp`/`mv`/`rm`/`tee`/
`git mv`/`git rm`), heuristically, not just the dedicated Edit/Write/`apply_patch` tools.

`AGENTS.md` rule 7 reworded to state both tiers. `workflows/bootstrap.md` and
`prompts/bootstrap.md` now state the `ANEW_ALLOW_CORE_EDIT=1` prerequisite explicitly, since
bootstrap is the one common flow that legitimately needs it from a cold start.

Deliberately excluded from the core set: `docs/architecture.md`, `docs/domain.md`,
`docs/conventions.md`, `docs/testing.md`, `docs/security.md`, `docs/git.md` — these are product
knowledge that's *supposed* to evolve as features land (Definition of Done: "Docs / ADRs updated
if behavior or architecture changed"). Locking them would fight the framework's own DoD.

## Consequences
**Buys us:** a technical answer to "how do I make sure a coding agent never touches the process
files while it's just doing feature work" — the actual question asked. Symmetric between Claude
Code and Codex. Catches Bash-based bypasses that the original `specs/done/` hook missed.

**Costs us:** first-run friction — bootstrap now visibly requires a step (`ANEW_ALLOW_CORE_EDIT=1`)
it didn't need before; undocumented, this would look like a broken tool on first use, hence the
explicit prerequisite notes added to both bootstrap files. The Bash-matching is still a substring
heuristic (documented as "a courtesy layer, not a security boundary" from the start) — a command
containing an escaped quote or indirection can still slip past it; found this directly while
testing (a naive test command with embedded `\"\"` broke the crude grep-based JSON field
extraction). A determined agent could also just set the env var itself if nothing stops it from
running arbitrary shell — this protects against silent/incidental edits, not against an agent that
decides to defeat its own guardrail on purpose.

## Alternatives considered
- **Subagent-scoped hooks** (Claude Code's per-subagent `hooks:` frontmatter field) restricting
  only `builder`/`test-writer`, leaving the main/planner session unrestricted: rejected as the
  primary mechanism — Codex has no equivalent per-skill hook scoping, so this would have made the
  two adapters asymmetric, and it doesn't protect a main session that bypasses the designed routing
  and does BUILD-type work directly. Kept as a possible future *addition*, not a replacement.
- **Unconditional block, no override**: rejected — would have made bootstrap and all of this
  session's own meta-work impossible, and there's no such thing as a repo that never needs its own
  process files touched (see Context).
- **Real JSON parsing (Python) instead of grep/sed for the Claude Code hook**: considered, for
  robustness against the escaped-quote issue found during testing; not done here to keep the two
  hooks in their existing idiomatic languages (`sh` for Claude Code matching the original,
  `python3` for Codex matching its own official example) — flagged as a revisit trigger instead of
  a silent scope increase.

## Revisit triggers
- The escaped-quote / crude-parsing limitation actually causes a missed block in practice (not
  just in a test) — rewrite `protect-shipped.sh`'s extraction with a real parser (e.g. shell out to
  `python3 -c` the way `policy.py` already does) instead of grep/sed.
- `ANEW_ALLOW_CORE_EDIT` gets left set in a persistent shell profile and stops meaning anything —
  consider `scripts/doctor` warning if it's set in the *current* environment when doctor runs, as a
  gentle reminder it's meant to be session-scoped, not permanent.
- A real incident where the core set was too narrow (something not listed here got silently
  clobbered) or too broad (legitimate work kept tripping the block) — adjust the CORE_PATTERNS
  list in both hook scripts together, never just one.
