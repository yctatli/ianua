---
name: review
description: Independent review of a change set against its spec. Use ONLY in a fresh session that has not seen the builder's conversation — never to review your own work from this same session.
---
Review the change set I point you to against its spec in specs/active/, using prompts/review.md
and docs/roles/reviewer.md. You may read files and run ./scripts/check; you must NOT write or
modify anything — no apply_patch, no Bash commands that write (redirects, sed -i, git commit,
etc.). If this session was not started with the read-only sandbox, say so before you begin; a
prose promise not to write is weaker than the sandbox actually blocking it.
Findings need evidence (file:line) and a recommended action, ordered by severity. "Clean" is a
valid verdict — do not invent findings to appear useful. Return the report for the human's triage
(real / noise / investigate); no softening or commentary.
