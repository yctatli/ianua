---
name: new-feature
description: Run the full feature-development workflow (intent → spec → plan → build → review → verify → ship) for a new piece of behavior. Use when the user wants to add or build a new feature.
---
Read AGENTS.md (note the operating mode) and workflows/feature-development.md. If I haven't
already described the feature, ask me for it, then run the workflow.

Honor every gate of the current mode — stop and wait at each [GATE: human]. Use the prompts the
workflow references (prompts/clarify.md, spec.md, plan.md, build.md) as written, filling
placeholders. For the REVIEW step, tell me to open a **new Codex session** (ideally started in the
read-only sandbox) that receives only the diff and the spec path, and run it with
prompts/review.md — you must not review your own change set in this session
(docs/roles/README.md: "the mind that produces cannot audit itself"). Bring its findings back to
me for triage before any fixes.
