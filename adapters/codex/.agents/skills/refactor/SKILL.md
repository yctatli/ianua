---
name: refactor
description: Run the refactor workflow — structure changes, behavior provably preserved. Use when the user wants to restructure or clean up existing code without changing behavior.
---
Read AGENTS.md and workflows/refactor.md. If I haven't said what to refactor and why, ask, then
run the workflow.
./scripts/check must be green before starting. If the area lacks tests, characterization tests
come first. Zero test edits — a needed test edit means behavior changed: stop and tell me.
