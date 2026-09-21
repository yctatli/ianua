---
description: Something went wrong — pick and run the right recovery ramp (R-01…R-12)
argument-hint: what happened
---
Read prompts/recovery/README.md. Based on what happened — "$ARGUMENTS" — identify the matching
ramp (R-01…R-12) yourself or via the `planner` subagent (diagnosis is planner-tier judgment, see
adapters/claude-code/README.md, "Model routing"). Tell me which ramp and why, and after my
confirmation execute its prompt from prompts/recovery/<file>.md exactly, including its protection
rules — hand purely mechanical steps (retry, revert) to `builder`. If nothing matches, treat it as
R-10 (ambiguity): options with costs, your recommendation, my decision.
