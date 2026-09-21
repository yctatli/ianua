---
description: Run the bug-fix workflow (reproduction before fix)
argument-hint: what is broken, expected vs actual
---
Read AGENTS.md and workflows/bug-fix.md. Run the workflow for: $ARGUMENTS

Route by tier (adapters/claude-code/README.md, "Model routing"):
- REPRODUCE → `test-writer` subagent: the minimal failing reproduction test (R-05 pattern). Iron
  rule: no fix before this test exists and is red for the right reason.
- DIAGNOSE → `planner` subagent: root cause, not symptom (R-01). Diagnosis before any change.
- FIX → `builder` subagent: smallest change that turns the reproduction test green.
- REVIEW → `reviewer` subagent, narrow scope: the fix diff only. Strict mode, or the fix touches
  auth/input/data handling: also `security` subagent (prompts/security-review.md).

Stop at every human gate. The reproduction test is permanent — never weakened or removed to reach
green.
