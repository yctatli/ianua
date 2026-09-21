---
description: Discuss an architecture decision, then record it as an ADR
argument-hint: the decision question
---
Delegate to the `planner` subagent (adapters/claude-code/README.md, "Model routing" — a wrong call
here is expensive downstream, and the output is a page, not a codebase): use prompts/adr.md for:
$ARGUMENTS
First discuss options with a recommendation and revisit triggers; wait for my decision; only then
write docs/decisions/NNNN-<slug>.md from the template, with honest consequences (benefits AND costs).
