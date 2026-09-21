---
description: Run the refactor workflow (structure changes, behavior provably preserved)
argument-hint: what to refactor and why
---
Read AGENTS.md and workflows/refactor.md. Run the workflow for: $ARGUMENTS
./scripts/check must be green before starting.

Route by tier (adapters/claude-code/README.md, "Model routing"):
- BASELINE (missing tests) → `test-writer` subagent: characterization tests that pin current
  behavior exactly as it is, warts included.
- SCOPE & PLAN → `planner` subagent: boundaries, "no observable behavior change", prompts/plan.md.
- REFACTOR → `builder` subagent: small, committable steps, suite green after each.
- REVIEW → `reviewer` subagent: extra lens — did semantics sneak in, are names/layers now more
  aligned with docs/architecture.md. Strict mode, or security-sensitive code: also `security`
  subagent (prompts/security-review.md) — structure changes can move a trust boundary silently.

Zero test edits during REFACTOR — a needed test edit means behavior changed: stop and tell me.
