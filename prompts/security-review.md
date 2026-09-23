# Prompt: Security Review

```
Role: Security (docs/roles/security.md) — fresh session or read-only subagent, independent from
both the builder's session and the general Reviewer's session. You may read files and run
./scripts/check (including its `security:` step); you may NOT write or modify anything.

Review this change set against docs/security.md and the spec in specs/active/<NNNN>-<name>.md:
<diff or branch reference>

Confirm first: did scripts/check's `security:` step run and pass? If missing or skipped, that is
itself a finding, not something to quietly work around.

For each category give concrete findings WITH EVIDENCE (file:line), or explicitly say "clean":
- AuthN/AuthZ: who may call what, is authorization enforced server-side, default-deny where it
  matters.
- Input & output: validation at boundaries, injection (SQL/command/template/deserialization),
  encoding, what leaks in error messages or logs.
- Secrets: none in code, logs, specs, or prompts; .env handling correct (docs/security.md).
- Dependencies: any newly introduced dependency with known CVEs, or added without the review
  docs/security.md calls for.
- Trust boundaries: does data cross one (network, process, tenant, privilege level) without being
  re-validated on the other side?

Assign each finding a severity from docs/security.md's "Severity & gating" table — Critical / High
/ Medium / Low, using its criteria, not a free-form guess — and a short ID (S-1, S-2, ...). Order
findings Critical first. Each finding carries a recommended action (proposal rule). If a finding is
Critical, say so plainly and first: work should stop here for immediate human attention (recovery
ramp R-13) — don't bury it at the end of a long report.
Do not invent findings to appear useful — "clean" is a valid verdict. I will triage:
real / noise / investigate — except Critical/High, which cannot be closed as "noise" per
docs/security.md — and record each one in the plan's Findings log; you're read-only, so that
recording is my job, not yours.
```
