# Security

> **Template — filled during bootstrap.** Baseline rules agents must honor in every plan and review.

## Secrets
- Secrets never enter the repo, specs, prompts, or chat. `.env` is gitignored; provide `.env.example`.
- Agents never print secret values, even when debugging.

## Input & output
<!-- Validation strategy at boundaries; what is escaped/encoded; what never leaks in error messages. -->

## AuthN / AuthZ
<!-- Who may call what; where authorization is enforced; default-deny statements. -->

## Dependencies
<!-- Policy for adding new dependencies: who approves, what gets checked (license, maintenance, CVEs). -->

## Automated checks (tooling, not prose)
<!-- Bootstrap fills this with real commands, stack-detected or chosen at interview time, and adds
them as a `security:` step in scripts/check.conf — so every `scripts/check` run is also a security
gate, not just something a reviewer eyeballs. Two kinds, add what applies:
- Dependency/CVE audit: `npm audit`, `pip-audit`, `dotnet list package --vulnerable`, `bundle audit`...
- Static analysis (SAST): `semgrep --config auto`, `bandit -r .`, language-specific linters with
  security rule sets enabled.
A step that isn't wired into scripts/check.conf is advice an agent can silently skip; wire it in. -->

## Severity & gating
Every security finding (from Security's pass, Reviewer's security dimension, or the automated
`security:` check) gets one of four severities — the finder assigns it from this table, not a
free-form guess:

| Severity | Examples | Ship rule |
|---|---|---|
| **Critical** | RCE, auth bypass, a secret committed to the repo, unauthenticated access to sensitive data, SQL injection with data access | Triage as "noise" is **not allowed**. Fix before ship, or ship only behind a dedicated ADR that names the accepted risk, its owner, and a revisit trigger (`docs/decisions/README.md`) — never a one-line triage note. Recovery ramp R-13. |
| **High** | Privilege escalation, significant data exposure, stored XSS, missing authZ on a sensitive endpoint | Blocks ship in **strict** mode. In **lite** mode, shipping without a fix needs explicit written human sign-off — not a silent "noise." |
| **Medium** | Boundary input-validation gaps, stack traces/internal errors leaking to users, missing rate limiting | May ship — but only with a tracked follow-up: a new `specs/active/` entry, not just a triage note that evaporates. |
| **Low** | Best-practice deviations, defense-in-depth suggestions, a non-vulnerable outdated dependency | Normal triage — human's call, deferrable freely. |

A vulnerability found in code that's **already shipped** (an ad hoc `/security` pass, not a
pre-ship review) is not a new spec — it's `workflows/incident.md`. The Critical/High ship rules
above still apply there; Critical additionally triggers recovery ramp R-13
(`prompts/recovery/critical-security-finding.md`) immediately, before anything else continues.

## Review lens
Security is a mandatory dimension of every independent review, not a separate afterthought phase —
automated checks catch known patterns, a human/model read catches design-level issues (authz gaps,
trust-boundary mistakes) no scanner will find. **In strict mode** this is a dedicated, independent
Security pass (`docs/roles/security.md`, `prompts/security-review.md`) — its own session, separate
from Reviewer's. **In lite mode** it collapses into Reviewer's own security dimension
(`prompts/review.md`) — still mandatory, just not a separate session by default. For a systematic
pass beyond either read, Claude Code users can also run the built-in `security-review` skill
against the change set — see `adapters/claude-code/README.md`.
