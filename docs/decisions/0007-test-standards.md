# ADR 0007 — Test standards: risk-tiered levels, mandatory categories, named anti-patterns

- Status: Accepted
- Date: 2026-09-22

## Context
The human asked the testing analog of the question behind ADR 0006: what should test standards be
based on, what do we actually have, what's missing. Tracing `docs/testing.md` found the same shape
of gap ADR 0006 found in security:

1. **"What must be tested" was empty** — a comment hint (`<!-- e.g. every business rule... -->`),
   never filled with an actual list, even though most of it is genuinely stack-agnostic and could
   be written now (unlike "Frameworks & layout," which really is stack-specific and correctly stays
   a bootstrap-filled template).
2. **No risk-tiered rigor.** Every acceptance criterion cleared the same bar — "at least one test
   asserting behavior" — regardless of whether it's a pure calculation or a money-moving,
   irreversible action. A mocked unit test was implicitly "enough" proof for anything.
3. **No named anti-patterns.** "Assert behavior, not implementation" is a good principle with
   nothing concrete under it — `docs/roles/reviewer.md`'s Tests dimension had no list to check
   findings against, so what counts as a violation depended on whoever was reviewing that day.
4. **Regression-test requirement inconsistent between ramps.** R-04 requires a permanent regression
   test explicitly; R-13 (new in ADR 0006) didn't repeat this for a fixed Critical/High security
   finding, even though the same logic applies — a fix without a test proving it stays fixed isn't
   actually done.
5. **`specs/plans/TEMPLATE.md`'s criterion↔test map had no test-level field** — "Proven by" alone
   doesn't say whether unit-level proof is sufficient, so that call had nowhere to be made
   explicitly at plan time (where it belongs, per this workspace's own "judgment upstream, execution
   downstream" pattern from ADR 0002) and would otherwise default to whatever the test-writing pass
   happened to pick.
6. **Recovery-ramp ownership for the Claude Code adapter was incomplete.** ADR 0002's model-routing
   table only assigned R-01/R-07/R-10 to `planner` explicitly; R-02 through R-13 were unassigned,
   leaving genuinely judgment-flavored ramps (R-02: is code, test, or spec wrong?) at risk of being
   decided by `builder` — the same execution-tier subagent whose own work might be in question.

## Decision
Fill `docs/testing.md`'s "What must be tested" with five concrete, stack-agnostic categories
(every acceptance criterion including boundaries/errors, every business rule, every forbidden-
dependency/architecture boundary as an automated architecture test where the stack supports one,
at least one E2E/smoke test per critical flow, and a permanent regression test for every fixed
R-04/R-13 finding). Add two new sections:
- **"Test level by risk"** — a table mapping what a criterion touches (pure logic / cross-module /
  money-auth-irreversible-data-loss / a fixed Critical-High security finding) to a minimum test
  level (Unit / Integration / E2E), explicit that a mocked unit test is not sufficient proof for
  the highest-risk row.
- **"Test smells to reject"** — five named anti-patterns (assertion-free tests, mock-tests-the-mock,
  order-dependent/shared-state tests, timing band-aids, testing private internals) for Reviewer's
  Tests dimension to check against by name.

`specs/plans/TEMPLATE.md`'s criterion↔test map gets a **Level** column — the test-level call is
made at plan time (`planner`), not left implicit for whoever writes the test.

Claude Code adapter: a new "Recovery ramp ownership" table assigns all 13 ramps to a subagent
(`adapters/claude-code/README.md`). Notably: **R-02 moves to `planner`** (deciding code-vs-test-vs-
spec-wrong is a judgment call, not `builder`'s to make about its own work); R-04 and R-13 are
explicitly multi-agent (diagnosis → fix → regression test, three different kinds of work, not one
subagent doing all three); R-05 runs under `builder` assuming the QA role, matching `/verify`'s
existing tier since there's no dedicated QA subagent.

## Consequences
**Buys us:** the same benefit ADR 0006 bought security — test rigor now depends on what a criterion
actually is, not on habit. Reviewer has a named checklist instead of a vibe. The test-level call
happens once, at plan time, instead of being silently re-decided (or not decided at all) by
whichever pass writes the test. Recovery-ramp ownership is no longer a gap a real incident would
have had to discover the hard way.

**Costs us:** the plan step gains one more explicit judgment call (test level per criterion) —
deliberate friction, same trade-off ADR 0006 made for severity. `docs/testing.md` is now
opinionated pre-bootstrap in a way the original template design left to bootstrap entirely; the
five categories added are believed stack-agnostic, but a genuinely unusual project might need to
prune one (e.g. no meaningful "architecture boundary" test is possible in every stack) — that's an
explicit revisit trigger, not assumed away.

## Alternatives considered
- **Leave "What must be tested" for bootstrap to fill, like "Frameworks & layout"**: rejected —
  unlike test framework/file-layout choices, which really are stack-specific, the five categories
  added here don't depend on the stack; deferring them gains nothing and loses the ability to
  reference them from other core, stack-agnostic files (this ADR, the ramp table) before bootstrap.
- **A numeric/formal test-coverage threshold (e.g. "80% line coverage")**: rejected — coverage
  percentage measures volume, not whether the highest-risk criteria got proportionate rigor; the
  risk-tiered table targets the actual risk instead of a number that can be gamed with low-value
  tests on easy code paths.
- **One subagent (`builder`) owns every recovery ramp**: rejected for the same reason ADR 0002
  split `planner`/`builder` in the first place — R-02/R-06/R-07/R-08/R-10 are diagnostic judgment
  calls, and letting the subagent whose own work triggered the ramp also decide "is my code, the
  test, or the spec wrong" reintroduces exactly the self-verification problem this workspace's
  Reviewer/Security separation exists to prevent.

## Revisit triggers
- A project bootstrapped from this template has a stack where the "architecture boundary as an
  automated test" category genuinely doesn't apply (no tooling exists for it) — note the exception
  in that project's own `docs/testing.md`, don't silently drop the category from the template.
- The Unit/Integration/E2E three-level scale proves too coarse for a real project (e.g. contract
  tests, mutation testing) — extend the table, keep the same risk-tiering structure.
- `builder`-as-QA (R-05) starts feeling under-powered in practice — reconsider a dedicated QA
  subagent tier rather than folding it into `builder`.
