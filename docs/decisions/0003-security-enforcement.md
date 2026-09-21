# ADR 0003 — Security as an enforced check, not just a review lens

- Status: Accepted
- Date: 2026-09-21

## Context
`docs/security.md` already stated a principle: "Security is a mandatory dimension of every
independent review... not a separate afterthought phase." In practice that meant security was
**prose inside a review checklist** — real, but only as strong as the reviewer's own judgment in
that pass, with no automated backstop and no dedicated tool invoked. The user asked explicitly for
every piece of written code to get security testing/analysis done, not just a mention in a list of
six review dimensions.

`AGENTS.md`'s own founding principle is directly relevant here: "Rules it cannot bypass... Prose is
advice; tooling is law." A `security:` line in `docs/security.md` that never becomes a
`scripts/check.conf` step is exactly the kind of rule an agent can silently skip under time
pressure — the same failure mode the whole `scripts/check` design exists to close for build/test/lint.

Separately, this Claude Code environment has a built-in `security-review` skill (systematic
security review of a change set) available to invoke — not previously wired into this workspace's
review process at all.

## Decision
Two changes, kept separate on purpose (automated ≠ human judgment, per `docs/security.md`'s own
distinction):
1. **Automated, tooling-enforced**: `docs/security.md` gets a new "Automated checks" section
   instructing bootstrap to wire a stack-appropriate dependency/CVE audit and/or SAST tool into
   `scripts/check.conf` as a `security:` step — `scripts/check.conf.example` now shows this per
   stack (.NET/Node/Python). `prompts/bootstrap.md`'s INTERVIEW and GENERATE steps now explicitly
   ask about and generate this step, instead of leaving it to be inferred from "security posture."
   `specs/TEMPLATE.md`'s Definition of Done gets an explicit line so a spec can't reach "done"
   with security silently unaddressed.
2. **Systematic, skill-backed (Claude Code only)**: the `reviewer` subagent is instructed to invoke
   the built-in `security-review` skill as part of covering the security dimension, on top of its
   own judgment. Documented in `adapters/claude-code/README.md` as also available on demand,
   outside the mandatory gate, for ad hoc deep dives.

Deliberately **not** done: no new "Security" role, no separate security-only workflow phase. Kept
consistent with the existing, already-considered position in `docs/security.md` that security is a
lens applied throughout (build discipline + automated checks + review dimension), not a bolt-on
phase at the end that's easy to schedule away under deadline pressure.

## Consequences
**Buys us:** a repo can no longer have a "filled-in-looking" `docs/security.md` with zero actual
enforcement — bootstrap won't consider security posture answered without a real `security:` command;
Claude Code reviews get a systematic tool pass, not just a human/model eyeballing a diff; the DoD
checklist makes "security wasn't silently skipped" independently checkable, same as every other
acceptance criterion.

**Costs us:** `security:` steps (SAST, dependency audits) are typically slower and noisier
(false positives) than lint/typecheck — bootstrap needs to pick tools and thresholds that don't make
`scripts/check` a chore to run locally, or people route around it. The `security-review` skill
dependency is Claude-Code-specific; Codex/Copilot/Cursor users get the automated-tooling half of
this decision but not the skill-backed half, until/unless an equivalent exists for those tools.

## Alternatives considered
- **A dedicated "Security" role** (5th role alongside Analyst/Developer/Reviewer/QA), with its own
  workflow phase: rejected — `docs/security.md` already made this call before this ADR existed
  ("not a separate afterthought phase"); a bolt-on phase is also the first thing dropped under
  deadline pressure, which defeats the point. Kept it as everyone's job, backed by tooling instead.
- **Security only as a `scripts/check` step, drop the skill invocation**: rejected — automated
  scanners catch known patterns (CVEs, common injection shapes) but not domain-specific authz
  logic; the reviewer's own read plus a systematic tool pass covers more than either alone.
- **A dedicated recovery ramp for security findings**: rejected — a real security finding from
  review is just a "real" finding through the existing triage path (R-05 if unverified, straight to
  `builder` §fixes if confirmed); a 13th ramp for something the existing ramps already route
  correctly would be an abstraction with no new behavior behind it.

## Revisit triggers
- Bootstrap runs against a stack with no obvious SAST/dependency-audit tool — the "Automated
  checks" section needs a documented fallback (manual periodic audit cadence) rather than staying
  silently unfilled.
- The `security-review` skill's behavior or availability changes, or an equivalent becomes
  available for Codex CLI — extend the skill-backed half of this decision to `adapters/codex/`.
- `security:` steps prove too slow/noisy in practice and get routinely skipped or `--no-verify`'d
  — that's a signal to split "fast" (SAST on changed files) from "full" (dependency audit,
  nightly/CI-only) rather than drop enforcement.
