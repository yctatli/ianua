# Adapter: Claude Code

Install: `./scripts/init claude-code` — copies `CLAUDE.md` and `.claude/` to the repo root.

What it adds on top of the core:

- **Slash commands** (`.claude/commands/`): `/bootstrap`, `/new-feature`, `/fix-bug`, `/refactor`,
  `/review`, `/security`, `/verify`, `/adr`, `/recover` — each one loads the matching workflow,
  honors the operating mode, and routes each step to the subagent below. Thin by design: they point
  to `workflows/` and `prompts/`, they don't restate them.
- **Five subagents** (`.claude/agents/`), each pinned to a model tier — see "Model routing" below:
  `planner` (opus), `builder` (sonnet), `test-writer` (haiku), `reviewer` (opus, read-only),
  `security` (opus, read-only).
- **Read-only reviewer and security subagents**: the "producer never verifies its own work" rule
  made *impossible to break* — neither has Edit/Write tools. `security` is a genuinely independent
  pass (docs/roles/security.md), not a bullet inside `reviewer` — mandatory as its own session in
  strict mode, collapsible into `reviewer` in lite mode (docs/security.md, "Review lens";
  `docs/decisions/0004-dedicated-security-role.md`).
- **Permission denies** (`.claude/settings.json`): force push, hard reset, `rm -rf` blocked by
  the tool, not by politeness.
- **Core-file-lock hook** (`.claude/hooks/protect-shipped.sh`, runs on Edit/Write/Bash): two tiers.
  `specs/done/` is always rejected, no override. This workspace's own process files (`AGENTS.md`,
  `workflows/`, `prompts/`, `scripts/`, `adapters/`, `docs/roles/`, `docs/decisions/`, spec/plan
  templates, and the installed `.claude/`/`.codex/`/`.agents/` themselves) are rejected too, unless
  `ANEW_ALLOW_CORE_EDIT=1` is set — deliberately, for bootstrap or a dedicated ADR, never as a side
  effect of ordinary feature/bugfix work. See `docs/decisions/0005-core-file-lock.md`.
- **Security, backed by a real tool, not just a checklist line**: both `reviewer` (lite-mode
  fallback) and `security` (dedicated pass) are instructed to invoke the built-in `security-review`
  skill as part of their pass. For an ad hoc deep dive outside any workflow gate — before a risky
  PR, or periodically — invoke `security-review` yourself, or run `/security`; neither needs a
  workflow step to be useful on demand.

All defaults are adjustable — see "Adapting it" in the root README.

## Model routing (cost/quality policy)

Two independent knobs, both set in subagent frontmatter, both worth setting deliberately rather
than leaving at session defaults:

- **`model`** — which model (opus/sonnet/haiku): the ceiling on capability and the dominant cost
  driver. Pick by *how wrong an answer can be* and *how often this step runs*.
- **`effort`** — `low`/`medium`/`high`/`xhigh`/`max` on models that support it: how many tokens
  that model spends thinking and responding *within* its own ceiling. Pick by *how much this
  specific step benefits from deeper reasoning at that model's own level*. **Haiku does not
  support `effort` at all** (not in Anthropic's supported-model list as of 2026-09) — its economy
  comes purely from being the smaller model, so `test-writer` has no `effort` line.

The overall goal: spend opus+high where a wrong call is costly and the output is short (specs,
plans, diagnoses, reviews), and spend a cheap model/lower effort where output is high-volume and
the judgment call was already made upstream (tests, implementation). Each subagent pins both
knobs in its own frontmatter, so this holds regardless of what model or effort *you* picked for
the top-level session — set your own session to `sonnet` at `medium` effort (`/model sonnet`,
`/effort medium`) for cheap orchestration, since the commands above do the actual thinking and
typing in subagents, not in this session.

| Subagent | Model | Effort | Used for | Why this tier |
|---|---|---|---|---|
| `planner` | opus | high | Clarify, Spec, Plan, ADR discussion, diagnosis (bug DIAGNOSE / incident ROOT CAUSE), judgment-heavy recovery ramps (R-01, R-07, R-10) | Low volume, high stakes — a bad plan or a wrong root cause is expensive several steps later. `high` is Opus's own recommended default; the agent is instructed to ask you to bump to `xhigh` for one turn on genuinely hard/ambiguous cases instead of that being the default spend. |
| `builder` | sonnet | medium | BUILD (implementation), FIX, REFACTOR steps, VERIFY (evidence auditing) | The hard judgment call already happened in `planner`; this is disciplined execution against a spec. `medium` is Sonnet 5's cost-saving step-down (~Sonnet 4.6 at `high`) — the agent is instructed to flag it and suggest `high` if it's visibly struggling on a step, rather than silently burning retries. |
| `test-writer` | haiku | — (unsupported) | Tests from the criterion↔test map, bug REPRODUCE, refactor characterization tests | Highest volume, lowest ambiguity — transcribes an already-specified case into a failing test. Escalates to you in plain text if the criterion is genuinely ambiguous rather than guessing. |
| `reviewer` | opus | high | Independent REVIEW of any change set | Quality-critical, adversarial audit — the one place a cheap model or a low effort setting catching fewer real bugs is a bad trade. Same tier as `planner`: comparable call volume (one review per feature), comparable cost of being wrong. |
| `security` | opus | high | Independent SECURITY pass (strict mode: mandatory separate session; lite mode: optional, or covered inside `reviewer`) | Same reasoning as `reviewer` — a missed real vulnerability is at least as costly as a missed correctness bug, often more. See `docs/decisions/0004-dedicated-security-role.md`. |

Reassign a tier by editing the `model:`/`effort:` lines in the matching `.claude/agents/*.md` file
— nothing else needs to change, the commands reference subagents by name, not by model. Don't
change top-level effort mid-session if you're relying on prompt caching for a long-running
subagent conversation (it invalidates the cached prefix) — these subagents are short-lived
one-shots per invocation, so that mostly doesn't bite here, but it's the reason the "bump to
xhigh/high" suggestions above are framed as a deliberate one-turn call, not a standing change.

## Recovery ramp ownership

`prompts/recovery/README.md` lists all 13 ramps tool-agnostically; this is which subagent should
run each one here, so "who handles this" isn't reinvented per incident. Judgment-flavored ramps go
to `planner`; mechanical/technical ones stay with `builder`; QA-flavored ones run under `builder`
assuming the QA role (same as `/verify`, since there's no dedicated QA subagent — QA's job is
evidence-auditing, not a distinct model tier from `builder`'s).

| Ramp | Owner | Why |
|---|---|---|
| R-01 Compile/runtime error | `planner` | Diagnosis is the judgment call; the fix itself may hand off to `builder` |
| R-02 Red test | `planner` | "Is code, test, or spec wrong" is a judgment call, not something to decide mid-implementation — don't let `builder` talk itself into "the test must be wrong" |
| R-03 Flaky test | `builder` | Finding the nondeterminism source (timing/ordering/state/I-O) is technical execution, not a judgment call about correct behavior |
| R-04 Regression | `planner` (diagnose narrower fix) → `builder` (revert + fix) → `test-writer` (permanent regression test) | Spans all three by design — don't run it in one subagent end to end |
| R-05 Investigate finding | `builder` (as QA) | QA's reproduction work — same tier `/verify` already uses |
| R-06 Fix-loop-limit exceeded | `planner` | "Analysis only, no code" per the ramp itself |
| R-07 Plan drift | `planner` | Already the plan's own author |
| R-08 Spec changed mid-work | `planner` | Spec/plan work, Analyst-flavored |
| R-09 Context fog / handoff | whichever session hit it | Tier-agnostic — it's a state-file/handoff problem, not a reasoning-depth one |
| R-10 Ambiguity | `planner` | Options-with-costs is exactly `planner`'s job |
| R-11 Safe rollback | `builder` | `git revert` + a risk report is execution, not a new judgment call |
| R-12 Performance target missed | `builder` | Measure → one optimization → re-measure is disciplined execution |
| R-13 Critical security finding | `security` (detects/escalates) → `builder` (fix) → `test-writer` (regression test) | Same multi-agent shape as R-04 — detection, fix, and proof are different kinds of work |
