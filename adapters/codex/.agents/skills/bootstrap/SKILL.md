---
name: bootstrap
description: Adapt this Ianua workspace to the current project (new or existing) through a short interview. Use when AGENTS.md still says "NOT CONFIGURED", or the user asks to bootstrap, set up, or initialize the workspace.
---
If the current working directory is itself named `.ianua`, STOP — this is Ianua's own template
(likely a nested install), not the actual project. Tell the human to `cd ..` and re-launch here
from the real project root instead.

Read AGENTS.md, then execute workflows/bootstrap.md using the prompt in prompts/bootstrap.md.
Interview one topic at a time, always with your recommendation (proposal rule). Stop at every
human gate. If the working tree has application code, ask domain depth (deep code read vs
interview-only) before the interview, with your recommendation — deep mode means you actually read
the business-logic files yourself before drafting docs/domain.md and docs/architecture.md, not
inferring them from the interview alone. When generating, keep the AGENTS.md invariant rules
verbatim and the file ≤ 40 lines. Finish by running ./scripts/doctor and ./scripts/check and
reporting results and open decisions.
