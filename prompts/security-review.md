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

Order findings by severity. Each finding carries a recommended action (proposal rule).
Do not invent findings to appear useful — "clean" is a valid verdict. I will triage:
real / noise / investigate.
```
