# Workflow: Bootstrap

**Goal:** adapt Ianua to *your* project — new or existing. Run once (rerunnable to revise).
**Prompt:** `prompts/bootstrap.md` · **Role:** Analyst (interview) → Developer (generation)

**Before starting:** `export IANUA_ALLOW_CORE_EDIT=1` (or the session equivalent) — bootstrap writes
`AGENTS.md`, `scripts/check.conf`, and (if adapters are installed) `.claude/`/`.codex/`/`.agents/`,
all of which are protected by the core-file-lock hook (AGENTS.md rule 7) outside this flag.

**Wrong-directory check, first:** if the current working directory is itself named `.ianua`, stop —
don't bootstrap Ianua's own template as if it were the project. This happens when someone clones
Ianua into `.ianua/` (a nested install, see `docs/decisions/0010-nested-install.md`) and then runs
their AI tool from inside that folder instead of from the real project root one level up. Tell them
to `cd ..` and re-launch there instead.

```
INSPECT → INTERVIEW → GENERATE → VERIFY → REPORT
```

1. **INSPECT.** The agent examines the working tree. Empty (only Ianua's own files) → greenfield flow.
   Contains code → adoption flow: detect stack(s), build system, test setup, and existing
   conventions *from the code*, to be confirmed rather than asked from scratch.

2. **INTERVIEW.** One topic at a time, every question carrying the agent's recommendation
   (proposal rule): product & domain terms, architecture style and module boundaries, forbidden
   dependencies, conventions that matter (data rules, error handling), testing expectations,
   security posture, git rules, **operating mode (lite/strict)**, and **chat language** (what
   language the agent should respond in day to day — recommend matching the language the human is
   writing this interview in; code, docs, and commits stay English by default regardless, unless
   the human says otherwise).
   **[GATE: human]** — your answers are the input; nothing is assumed.

3. **GENERATE.** The agent fills `docs/*.md` from the interview, writes `scripts/check.conf`
   (build/test/lint commands for your stack), sets the mode and chat-language lines in `AGENTS.md`,
   and rewrites `AGENTS.md`'s project summary — **keeping the invariant rules block verbatim** and
   keeping the file ≤ 40 lines. For existing repos it may also propose toolchain steps for
   `.github/workflows/check.yml`.

4. **VERIFY.** Run `./scripts/doctor` (structure + configuration) and `./scripts/check`
   (must pass; in an empty greenfield it may be a no-op with a note). Context quiz: open a *fresh*
   session and ask a project question (e.g. "what type do money fields use?") — the agent must
   answer from files. If it can't, the docs aren't teaching; fix them.

5. **REPORT.** What was generated, what was assumed, what still needs a human decision.
   **[GATE: human]** — you approve the workspace before the first feature starts.
