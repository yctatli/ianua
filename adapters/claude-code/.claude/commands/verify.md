---
description: Verify acceptance criteria against concrete evidence (QA role)
argument-hint: spec number
---
Delegate to the `builder` subagent (adapters/claude-code/README.md, "Model routing" — evidence
auditing is execution-tier, not planning-tier): assume the QA role (docs/roles/qa.md) for spec
$ARGUMENTS using prompts/verify.md, and produce the criterion ↔ evidence table. Run tests to
capture real output; a claim without evidence is a gap, and gaps are findings. Do not write or
modify any production code — evidence gathering only.
