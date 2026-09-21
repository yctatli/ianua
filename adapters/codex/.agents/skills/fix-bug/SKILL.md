---
name: fix-bug
description: Run the bug-fix workflow — reproduction before fix. Use when something works incorrectly and the user wants it fixed.
---
Read AGENTS.md and workflows/bug-fix.md. If I haven't said what's broken (expected vs actual), ask,
then run the workflow.
Iron rule: no fix before a failing minimal reproduction test exists (R-05 pattern). Diagnose root
cause before changing anything (R-01). The reproduction test is permanent. Stop at every human gate.
