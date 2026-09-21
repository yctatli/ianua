# ADR 0001 — Add a Codex CLI adapter

- Status: Accepted
- Date: 2026-09-21

## Context
This workspace ships adapters for Claude Code, GitHub Copilot, and Cursor, plus a generic fallback.
The user also drives this project from OpenAI's Codex CLI and wants the same rules (AGENTS.md,
workflows, prompts) honored there, with tool-specific enforcement where Codex actually supports it
— not just a pointer file.

Codex CLI turns out to read `AGENTS.md` natively already (root-down, plus the user's own
`~/.codex/AGENTS.md`) — that convention is where `AGENTS.md` originated — so no pointer file is
needed, unlike the Copilot/Cursor adapters. Codex also has two real, project-scoped mechanisms we
can use instead of prose: **Skills** (`.agents/skills/<name>/SKILL.md`, repo-shareable, replacing
the now-deprecated user-only "custom prompts") and a **PreToolUse hook** contract
(`.codex/hooks.json` + a script, JSON on stdin, `permissionDecision: deny` to block) that closely
mirrors Claude Code's own hook contract.

Both of those only load when the project is marked **trusted** in the user's own
`~/.codex/config.toml` — an untrusted project silently skips `.codex/` entirely, no error.

## Decision
Add `adapters/codex/` with: 8 skills under `.agents/skills/` (one per workflow, mirroring the
Claude Code slash commands), a `.codex/hooks.json` + `.codex/hooks/policy.py` PreToolUse hook that
blocks writes under `specs/done/` and a handful of destructive git ops (force push, hard reset,
rebase, `rm -rf`), and a suggested `.codex/config.toml` (`sandbox_mode = "workspace-write"`,
`approval_policy = "on-request"`). Wire it into `scripts/init codex` and `scripts/doctor`.

## Consequences
**Buys us:** the same invariant rules (specs/done/ immutability, forbidden git ops) enforced
mechanically in Codex, not just in Claude Code; a workflow entry point (`$bootstrap`, `$new-feature`,
...) that matches the Claude Code slash commands one-to-one; zero duplication of core rules — the
skills are thin pointers to `workflows/` and `prompts/`, same as the Claude Code commands.

**Costs us:** the policy hook is a courtesy layer, not a security boundary — `policy.py` does
substring/regex matching on shell commands, the same class of heuristic as the Claude Code hook it
was ported from, and both can be worked around by a sufficiently adversarial or careless command.
The hook and skills are inert until the human manually trusts the project — a step this workspace
cannot perform for them, and a silent-failure mode if they forget (no error, just nothing firing).
Codex's Skills/hooks surface is newer and moved fast enough during research (multiple conflicting
third-party sources on the skills directory path) that it may drift from what's documented here;
`adapters/codex/README.md` flags this explicitly.

## Alternatives considered
- **Pointer-file-only adapter** (like Copilot/Cursor): rejected — Codex's hook and skills
  mechanisms are real and closely mirror Claude Code's, so matching that depth was worth the
  extra surface, especially for the `specs/done/` immutability guarantee.
- **One shared cross-tool hook script** invoked by both Claude Code's and Codex's hook configs:
  rejected for v1 — the stdin/stdout JSON contracts differ enough (field names, tool names) that a
  single script would need tool-detection branching for little real reuse; two small, readable
  scripts beat one branchy one.

## Revisit triggers
- Codex CLI changes the Skills directory convention or the PreToolUse JSON schema (watch for this
  specifically — it was the least stable fact found during adoption).
- Codex ships a native declarative command allow/deny list (currently absent) — would let
  `policy.py`'s git-guardrail portion move out of the hook and into `.codex/config.toml`.
