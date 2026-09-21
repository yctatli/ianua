---
name: recover
description: Something went wrong in the workflow — pick and run the matching recovery ramp (R-01...R-12). Use for red tests, flaky tests, plan drift, ambiguity, regressions, or any stuck state.
---
Read prompts/recovery/README.md. Based on what I describe, identify the matching ramp (R-01…R-12),
tell me which one and why, and after my confirmation execute its prompt from
prompts/recovery/<file>.md exactly, including its protection rules. If nothing matches, treat it
as R-10 (ambiguity): options with costs, your recommendation, my decision.
