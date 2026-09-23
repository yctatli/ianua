# AGENTS.md — Project Rules

> **STATUS: NOT CONFIGURED.** This workspace has not been adapted to a project yet.
> The only correct first action is the bootstrap workflow (`workflows/bootstrap.md`).
> Until bootstrap completes and rewrites this file, do not write application code.
>
> **Wrong directory?** If your current working directory is itself named `.ianua` (or you were
> launched from inside one), STOP — this is Ianua's own template, not the actual project. Tell the
> human: `cd ..` (go up one level, out of `.ianua/`) and re-launch your AI tool from the real
> project root instead — that's where `./.ianua/scripts/init <adapter>` installs `AGENTS.md` /
> `CLAUDE.md` / `.claude/` for it. Do not bootstrap `.ianua/` itself as if it were the project.

## Operating mode

**Mode: unset** — bootstrap sets this to `lite` or `strict` (see `workflows/README.md`).
Every workflow honors the gates of the current mode.

**Chat language: unset** — bootstrap asks and sets this (e.g. `Turkish`; code/docs/commits stay
English unless stated otherwise). Until set, respond in whatever language the human writes in.

## Invariant rules (these survive bootstrap — never delete or weaken them)

*Nested install: if a `.ianua/` folder exists at this project's root, only `AGENTS.md`, `CLAUDE.md`,
and `.claude/`/`.agents/`/`.codex/` are actually there — every other bare path below or in any
workflow/prompt/command it points to (`workflows/x`, `docs/y`, `specs/z`, `prompts/w`,
`scripts/check`, ...) lives under `.ianua/` instead. Resolve `workflows/x` as `.ianua/workflows/x`,
etc., for every such path you hit — don't conclude you need to run from inside `.ianua/`.*

1. **No spec, no code.** Every piece of work starts as a spec in `specs/active/` (from `specs/TEMPLATE.md`) — except the narrow trivial-change exception in `workflows/README.md`.
2. **Plan before build.** A human approves the plan before any code is written.
3. **The producer never verifies its own work.** Review and QA run in a separate session or a read-only subagent, working from files (diff + spec), never from the builder's chat.
4. **Evidence over claims.** "Done" requires `scripts/check` green and every acceptance criterion mapped to proof. Never claim completion without showing evidence.
5. **Tests are protected.** Weakening asserts, deleting or skipping tests to get to green is forbidden — always.
6. **Proposal rule.** Every question, option, or finding comes with your own recommendation and rationale. The human decides; nothing is applied without approval.
7. **Shipped specs are immutable, and so is this workspace's own process.** Files under `specs/done/` are never edited. `AGENTS.md`, `workflows/`, `prompts/`, `scripts/`, `adapters/`, `docs/roles/`, `docs/decisions/`, and the spec/plan templates change only via bootstrap or a dedicated ADR, gated by `IANUA_ALLOW_CORE_EDIT=1` — never as a side effect of feature/bugfix work.
8. **Uncertainty is surfaced, not assumed.** On ambiguity or a docs/code conflict: stop and use the matching recovery ramp (`prompts/recovery/`).

## Where things live

| What | Where |
|---|---|
| Architecture & boundaries | `docs/architecture.md` |
| Domain language & business rules | `docs/domain.md` |
| Coding conventions | `docs/conventions.md` |
| Testing rules | `docs/testing.md` |
| Security rules | `docs/security.md` |
| Git & branching rules | `docs/git.md` |
| Decisions with rationale (ADRs) | `docs/decisions/` |
| Roles (who may do what) | `docs/roles/` |
| Specs & plans | `specs/active/` · `specs/plans/` · shipped → `specs/done/` |
| Processes & gates | `workflows/` |
| Reusable prompts & recovery ramps | `prompts/` |
| The single verification command | `scripts/check` |
