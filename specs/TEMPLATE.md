# Spec NNNN — <short name>

- Status: Draft | Approved | In progress | Shipped
- Mode: lite | strict (from AGENTS.md at creation time)
- Plan: `specs/plans/NNNN-plan.md`
- Epic: `specs/epics/NNNN-<name>.md` (optional — only if this is one of several specs toward one shared outcome)

## Intent
<!-- 3–5 sentences, BUSINESS language: who wants this and why; what success looks like;
what is deliberately NOT being done. Decisions from clarifying questions are folded in here. -->

## Requirements
<!-- Behavior, not technical solutions. No endpoints, tables, or caches here — the plan owns "how". -->

## Constraints & out of scope
<!-- Hard limits (performance targets, compliance) and conscious V1 exclusions. -->

## Acceptance criteria
<!-- Each line independently testable. Cover the happy path, boundaries, empties, and error cases. -->
- [ ] AC-1 —
- [ ] AC-2 —
- [ ] AC-3 —

## Definition of Done
- [ ] Every acceptance criterion mapped to proof (test or reproducible observation)
- [ ] `scripts/check` green (this includes the security step — see docs/security.md)
- [ ] Independent review done; real findings fixed, noise rejected with written rationale
- [ ] Security dimension of review explicitly addressed: clean, or findings triaged by severity
      (`docs/security.md`) — no open Critical/High finding without a fix or a dedicated
      risk-acceptance ADR
- [ ] Docs / ADRs updated if behavior or architecture changed
- [ ] Spec moved to `specs/done/` (it becomes immutable there)

## Scorecard (fill at ship — honest numbers make the process improvable)
| Metric | Value |
|---|---|
| Spec revisions | |
| Fix rounds | |
| Review findings: real / noise | |
| Regressions introduced | |
| Bugs escaped to production | |

## Lesson (only if the scorecard shows real friction)
<!-- Spec revisions > 1, fix rounds > 2, or a real finding that shouldn't have made it past BUILD?
Name ONE concrete thing that would have prevented it — a new test category (docs/testing.md), a
new docs/conventions.md rule, or an ADR. Clean numbers → skip this section entirely; it exists to
close a real gap, not to manufacture ceremony for work that went fine. -->
