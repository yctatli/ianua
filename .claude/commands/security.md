---
description: Independent security review of the current change set (dedicated pass, separate from /review)
argument-hint: spec number or branch/diff reference
---
Delegate to the `security` subagent: review the change set for $ARGUMENTS against docs/security.md
and its spec in specs/active/, using prompts/security-review.md. You (the main session) must not
review it yourself. Return the subagent's findings verbatim for my triage (real / noise /
investigate), with no softening or commentary.

Use this alongside `/review`, not instead of it — they cover different dimensions
(docs/roles/security.md vs docs/roles/reviewer.md). In strict mode both are required for every
feature/fix REVIEW step; in lite mode this is optional, worth running whenever the change touches
auth, input handling, data access, or dependencies.
