---
name: bootstrap
description: Adapt this Ianua workspace to the current project (new or existing) through a short interview. Use when AGENTS.md still says "NOT CONFIGURED", or the user asks to bootstrap, set up, or initialize the workspace.
---
Read AGENTS.md, then execute workflows/bootstrap.md using the prompt in prompts/bootstrap.md.
Interview one topic at a time, always with your recommendation (proposal rule). Stop at every
human gate. When generating, keep the AGENTS.md invariant rules verbatim and the file ≤ 40 lines.
Finish by running ./scripts/doctor and ./scripts/check and reporting results and open decisions.
