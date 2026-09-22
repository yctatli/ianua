# Testing

> **Template — filled during bootstrap.**

## The contract
- Every acceptance criterion maps to at least one test (criterion ↔ test map lives in the plan).
- Tests assert **behavior**, not implementation details or mere status codes.
- The whole suite runs inside `scripts/check` — one command, everywhere.

## Frameworks & layout
<!-- Test framework(s), where tests live, naming pattern. -->

## What must be tested
- Every acceptance criterion (the contract, above) — including the boundary/empty/error cases
  `specs/TEMPLATE.md`'s own Acceptance Criteria section already requires, not just the happy path.
- Every business rule (BR-n) referenced in `docs/domain.md`.
- Every forbidden dependency / architectural boundary in `docs/architecture.md` — as an automated
  architecture test where the stack supports one (e.g. dependency-cruiser, ArchUnit, import-linter),
  not left to a human noticing it in review.
- At least one critical/happy-path user flow, end-to-end or smoke — unit coverage of the pieces
  isn't evidence the flow actually works wired together.
- Every fixed regression (R-04) and every fixed Critical/High-severity security finding (R-13,
  `docs/security.md`) — a permanent regression test, at the level that would actually have caught
  it. This is not optional: a fix without a regression test is not "done" per R-04/R-13's own
  protection clauses.

## Test level by risk
Not every acceptance criterion deserves the same rigor — the level is a plan-time decision
(`specs/plans/TEMPLATE.md`'s criterion↔test map "Level" column), not left to whoever writes the
test:

| Criterion touches... | Minimum level |
|---|---|
| Pure logic / calculation, no external interaction | Unit |
| Cross-module or cross-service interaction | Integration |
| Money, auth, an irreversible action, or data loss | Integration/E2E — a mocked unit test alone is not sufficient proof |
| A fixed Critical/High security finding | Whatever level would actually have caught it (often integration/E2E, since the miss was usually at a boundary) |

## Test smells to reject
Concrete instances of "asserts implementation, not behavior" (the contract, above) — Reviewer's
Tests dimension checks for these by name, not just in spirit:
- **Assertion-free tests.** Runs code, checks it didn't throw, asserts nothing about the actual
  result. Green forever, proves nothing.
- **Mock-tests-the-mock.** The assertion only confirms a mock was called with certain arguments,
  never exercises real behavior. If the mock is wrong, the test still passes.
- **Order-dependent tests / shared mutable state.** Passes in the suite, fails alone (or vice versa)
  — a correctness bug in the test, not a flake to route around (see Determinism, below).
- **Timing band-aids.** Arbitrary `sleep()` instead of waiting on a real condition — forbidden by
  R-03, called out again here because it's a test-writing habit, not just a fix-time shortcut.
- **Testing private/internal details directly** instead of the public behavior contract — breaks on
  refactors that don't change behavior, which is exactly what `workflows/refactor.md`'s "zero test
  edits" rule is designed to catch and reject.

## Protected-tests rule
Weakening asserts, deleting, or skipping tests to reach green is forbidden. A red test triggers
`prompts/recovery/red-test.md` (R-02) — first decide what is wrong: code, test, or spec.

## Determinism
Flaky tests are fixed, not retried or skipped — see R-03. Evidence of a fix: 5 consecutive green runs.
