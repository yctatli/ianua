---
description: Run the feature-development workflow (intent → spec → plan → build → review → verify → ship)
argument-hint: short feature description
---
Read AGENTS.md (note the operating mode) and workflows/feature-development.md.
Run the workflow for: $ARGUMENTS

Honor every gate of the current mode — stop and wait at each [GATE: human]. Route each step to the
subagent it belongs to, by cost/quality tier (adapters/claude-code/README.md, "Model routing"):
- CLARIFY, SPEC, PLAN → `planner` subagent (prompts/clarify.md, spec.md, plan.md).
- BUILD → per plan step: `test-writer` first (red test from the criterion↔test map), then
  `builder` to turn it green (prompts/build.md). Builder never writes the tests itself.
- REVIEW → `reviewer` subagent, diff + spec path only. Bring its findings back to me for triage
  before any fixes — you (this session) never review your own work.
- Fix rounds (real findings only) → `builder`.
- VERIFY → `builder`, using prompts/verify.md and docs/roles/qa.md; no production-code writes.

You orchestrate — gate by gate — and hand each step's thinking or typing to its subagent rather
than doing it in this session.
