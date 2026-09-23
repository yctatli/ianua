# Plan NNNN — <spec short name>

- Spec: `specs/active/NNNN-<name>.md`
- Status: Awaiting approval | Approved | Superseded
- Approved by / on: <!-- human name + date — a plan without this line is not approved -->

## Files to change or add
<!-- Path by path. This list is the blast radius; silent additions are plan drift (R-07). -->

## Ordered steps
<!-- Which layer/part first and why. Keep steps commit-sized. -->

## Risks & open questions
<!-- Each with YOUR recommendation and rationale (proposal rule). Unresolved questions block build. -->

## Criterion ↔ test map
<!-- Level: Unit | Integration | E2E — see docs/testing.md "Test level by risk". This is a plan-time
call (judgment), not something the test-writing pass decides for itself. -->
| Acceptance criterion | Level | Proven by |
|---|---|---|
| AC-1 | | |
| AC-2 | | |

## Findings log
<!-- One row per review finding, across every round — appended, never deleted, only status-updated.
The reviewer/security subagent is read-only and reports findings in chat; it cannot write here.
Recording each finding + its triage disposition is the orchestrating session's job, done as part of
TRIAGE (workflows/feature-development.md step 7, workflows/refactor.md's REVIEW step) — this is
what makes "did review happen, what did it find, is it actually fixed" answerable from a file
instead of only from chat history. scripts/status (and tui/) roll this table up per spec. -->
| ID | Round | Source | Severity | Finding | Status | Resolution |
|---|---|---|---|---|---|---|
<!-- Status: Open · Fixed · Noise (rejected — say why in Resolution) · Deferred (say to which later
step/spec, and why, in Resolution). Security Critical/High may never be Noise (docs/security.md,
"Severity & gating") — Fixed or a dedicated risk-acceptance ADR referenced in Resolution instead. -->
