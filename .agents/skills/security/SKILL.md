---
name: security
description: Independent security review of a change set — a dedicated pass, separate from the general review skill. Use for the strict-mode SECURITY step of any workflow, or on demand for a security-sensitive change (auth, input handling, data access, dependencies) in lite mode.
---
Review the change set I point you to against docs/security.md and its spec in specs/active/, using
prompts/security-review.md and docs/roles/security.md. You may read files and run ./scripts/check
(confirm its `security:` step ran and was green — its absence is a finding); you must NOT write or
modify anything.

Findings need a category (authN/authZ, input/output, secrets, dependencies, trust boundaries),
evidence (file:line), and a recommended action, ordered by severity. "Clean" is a valid verdict —
do not invent findings to appear useful. Return the report for the human's triage (real / noise /
investigate); no softening or commentary.

Use this alongside the `review` skill, not instead of it — they cover different dimensions. In
strict mode both are required for every feature/fix REVIEW step; in lite mode this one is optional,
worth running whenever the change touches auth, input handling, data access, or dependencies.
