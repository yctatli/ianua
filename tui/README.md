# ianua-tui

An optional, read-only terminal dashboard for an Ianua workspace — the same information
`scripts/doctor`, `scripts/status`, and `scripts/init` already expose as plain text, in one
interactive view. Built in the spirit of [lazyskills.sh](https://lazyskills.sh/): "mission control"
for what's installed and what's in flight, without leaving the terminal.

This is **not part of core.** Ianua's core (`docs/`, `specs/`, `workflows/`, `prompts/`, `scripts/`,
`AGENTS.md`) stays plain Markdown and POSIX shell, on purpose (see the root README, "Design
principles"). `tui/` is a separate Go module that *consumes* the core's file conventions — it
never writes anything, and nothing in core depends on it existing.

## Install & run

Requires Go 1.24+.

```bash
cd tui
go build -o ianua-tui .
./ianua-tui              # from a classic install root, or the directory containing .ianua/
```

Or run without building a binary first: `go run . ` from inside `tui/`.

In a **nested install**, run it from the project root (the directory that *contains* `.ianua/`),
same rule as launching Claude Code or Codex — never from inside `.ianua/` itself. It walks up from
the current directory looking for `AGENTS.md` + `workflows/bootstrap.md`, checking a nested
`.ianua/` at each level first, so it also works from a subdirectory of the project.

`--check` prints the Health view as plain text and exits non-zero on any FAIL — no TUI, no
terminal required, useful in a script or CI step alongside (not instead of) `scripts/check`:

```bash
./ianua-tui --check
```

## Views

| Tab | Mirrors | Shows |
|---|---|---|
| **Health** | `scripts/doctor` | Structure completeness, recovery ramp count, AGENTS.md configuration state (status/mode/chat language/invariant block), `check.conf` presence, adapter presence, hook executable bits |
| **Status** | `scripts/status` | Every spec in `specs/active/`: status, mode, plan written?, Definition-of-Done progress, epic membership — plus any epics |
| **Adapters** | `scripts/init` | Mode + chat language (from `AGENTS.md`), and which adapter discovery files exist at the project root, with symlink health (a nested install's `AGENTS.md`/`CLAUDE.md`/`.claude`/etc. should all resolve — a broken one shows in red) |

Keys: `tab` / `←` `→` (or `1`/`2`/`3`) to switch tabs, `r` to refresh, `q`/`esc`/`ctrl+c` to quit.

## Why re-implement the checks instead of shelling out to the scripts

`scripts/doctor` and `scripts/status` print human-readable text meant for a terminal, not a stable
machine format — parsing that back out would be fragile and would silently drift the moment either
script's wording changes. `tui/health.go` and `tui/specs.go` re-implement the same checks natively
in Go, reading the same files. If the two ever disagree, that's a bug in one of them — file it
the same way as any other correctness bug (`workflows/bug-fix.md`), the check logic isn't special.

## Scope

Read-only by design, matching the "producer never verifies its own work" spirit elsewhere in this
project: this is a status *viewer*, not a workflow driver. It doesn't run `/bootstrap`, doesn't
write specs, doesn't install adapters for you — `scripts/init` still does that. If the Adapters tab
shows something missing, the fix is still `./scripts/init <adapter>` (or `./.ianua/scripts/init
<adapter>` when nested), run yourself.
