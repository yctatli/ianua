# ADR 0002 — Per-role model and effort routing for Claude Code subagents

- Status: Accepted
- Date: 2026-09-21

## Context
The user's stated goal: more quality work per token spent, not the cheapest possible setup. Running
every workflow step in whatever model the top-level session happens to be on wastes the expensive
model on mechanical work (writing tests from an already-specified criterion) and risks under-using
it on judgment-heavy work (spec ambiguity, root-cause diagnosis, plan risk assessment) if the human
forgot to switch up before those steps.

Claude Code subagents (`.claude/agents/*.md`) each pin their own `model` and, independently,
`effort` in frontmatter — `effort` (low/medium/high/xhigh/max) controls how many tokens a model
spends *within its own capability ceiling*, separate from which model is used. Confirmed during
research: Haiku does not support the `effort` parameter at all (absent from Anthropic's supported-
model list as of 2026-09) — its economy is the smaller model itself, not a dialed-down effort knob.

Before this decision, only one subagent existed (`reviewer`), with no `model:` pin at all — it
silently inherited whatever the session was on.

## Decision
Split workflow execution across four subagents, each with an explicit `model` + (where supported)
`effort`:
- `planner` — opus, high. Clarify/Spec/Plan/ADR/diagnosis/judgment-heavy recovery ramps.
- `builder` — sonnet, medium. BUILD/FIX/REFACTOR execution and VERIFY evidence-auditing, once a
  plan or diagnosis exists.
- `test-writer` — haiku, no effort field. Test authorship only, from an already-specified
  criterion/bug report/characterization target — never decides what "correct" means.
- `reviewer` — opus, high (newly pinned; previously unset). Independent review — kept at the same
  tier as `planner` since a missed real finding is as costly as a bad plan.

Every workflow slash command (`/new-feature`, `/fix-bug`, `/refactor`, `/verify`, `/adr`,
`/recover`) now explicitly routes each step to the matching subagent, so routing happens by default
through normal use of the commands rather than requiring the human to remember to delegate.
Each "cheap" subagent (`builder`, `test-writer`) is instructed to flag it in plain text and suggest
stepping up a tier if it's visibly struggling, rather than silently retrying or guessing.

## Consequences
**Buys us:** the expensive model is spent where research says it matters most (Opus 5's own
guidance: start at `high`, reserve `xhigh`/`max` for demonstrated need) and nowhere else by default;
test-writing and mechanical execution run on cheaper models without a quality cliff, because the
judgment call was already made upstream by `planner`; the routing survives whatever model the human
picks for their own top-level session, since it's pinned per-subagent, not inherited.

**Costs us:** four files to keep in sync instead of one; the `builder` effort choice (`medium`) is
the most debatable call in this set — Anthropic's own docs generally file "coding" under `high`,
and this bets that an already-detailed, already-approved plan lowers the bar enough for `medium` to
hold quality. Unverified in practice yet — no real feature has been run through this routing.
`test-writer` on Haiku carries the most capability risk of the four; its own hard rule (escalate
instead of guessing on ambiguity) is the only guardrail against a small model quietly asserting the
wrong behavior into a test.

## Alternatives considered
- **One subagent for everything, model chosen per session by the human**: rejected — relies on the
  human remembering to switch before every step, which is exactly the discipline this workspace
  tries to encode as tooling rather than habit (`AGENTS.md`'s own "rules it cannot bypass" principle).
- **`builder` at `high` effort instead of `medium`**: considered the safer default; rejected for now
  in favor of `medium` because the cost-saving case is strongest for the "hard thinking already
  happened" execution step — but flagged as the first thing to revisit if quality regresses.
- **Merge `test-writer` into `builder`** (one execution-tier subagent that also writes tests):
  rejected — collapses the one place a genuinely cheaper model (not just lower effort) is safe to
  use, since test authorship from an already-specified criterion is lower-ambiguity than general
  implementation.

## Revisit triggers
- A real feature run through this routing shows `builder`'s `medium` effort producing rework
  (fix rounds exceeding the plan's own estimate, or the "step up to `high`" self-report firing
  often) — raise its floor to `high`.
- `test-writer` on Haiku produces tests that pass for the wrong reason, or that assert the wrong
  behavior with unwarranted confidence instead of escalating — move it to `sonnet` (still without
  `effort`, or reconsider once/if Haiku gains `effort` support).
- Anthropic changes which models support `effort`, or changes the level names/semantics — re-check
  against `platform.claude.com/docs/en/build-with-claude/effort` before assuming this table still
  holds.
