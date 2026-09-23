# ADR 0012 — Optional TUI dashboard (`tui/`), outside core, a deliberate dependency exception

- Status: Accepted
- Date: 2026-09-23

## Context
The human asked for something in the spirit of [lazyskills.sh](https://lazyskills.sh/) — a "mission
control" terminal UI for managing AI-tool skills across many tools, with install/update/remove,
per-tool visibility, and broken-config detection. Ianua's own equivalent information already exists
as three separate plain-text commands (`scripts/doctor`, `scripts/status`, `scripts/init`), which is
consistent with the stated design principle "No frameworks, no dependencies, no code generators — a
working system, written in Markdown plus a handful of small scripts" (root README).

A real TUI is a real dependency: a compiled binary, a language toolchain (Go), a third-party
library (Bubble Tea) for anyone who wants to build it themselves. That directly conflicts with the
principle above if it's treated as part of core. The question was whether to fold dashboard
behavior INTO the POSIX-shell scripts (keeping zero dependencies, losing interactivity — a shell
script can't easily do tabs/live re-render/color without reaching for something like `gum` anyway,
which is itself a dependency) or accept the dependency but keep it strictly optional and separate.

## Decision
Add `tui/` as its own Go module (`go.mod` inside `tui/`, not at the repo root), entirely outside
the core-file-lock's protected set (`is_core_path()` in `protect-shipped.sh` was not extended to
cover it — it's ordinary, unprotected application code, edited the same way any other feature
would be). It re-implements `scripts/doctor` and `scripts/status`'s checks natively in Go (reading
the same files, not shelling out and parsing text output — see `tui/README.md`, "Why re-implement
the checks") across three read-only views (Health, Status, Adapters), plus a `--check` flag for a
non-interactive/CI-friendly mode. It never writes a file — `scripts/init` remains the only thing
that installs an adapter; the Adapters tab only reports what it finds, including nested-install
symlink health (a broken symlink shows in red — this is the same nested-install failure mode ADR
0010 and the wrong-directory guard fix both already deal with, just surfaced visually here instead
of via a workflow warning).

Both READMEs mention it as one optional line in "Monitoring progress," not a required step in
Quickstart — finding out about it should never be a precondition for using Ianua at all.

## Consequences
**Buys us:** the LazySkills-style experience the human asked for (one screen instead of three
commands, live symlink-health for nested installs) without touching core's own guarantee that it
stays plain-text and dependency-free — someone who never installs Go still has a fully working
Ianua.

**Costs us:** a second implementation of the doctor/status logic that must be kept in sync with the
shell scripts by hand (flagged as a revisit trigger below, same shape as ADR 0005's own "two hooks,
kept in sync by hand" tradeoff) — a real risk, not a hypothetical one, since it already happened
once with the two per-tool core-file-lock hooks. `tui/` also needs someone with Go on their machine
to build it; that's an acceptable cost for an explicitly optional tool, not for anything in core.

## Alternatives considered
- **Fold dashboard behavior into `scripts/doctor`/`scripts/status` themselves** (e.g. a `--watch`
  flag, ANSI color): rejected — gets you re-running text output, not tabs/navigation/live state;
  the "one screen" value LazySkills demonstrates needs an actual TUI framework, and POSIX shell has
  no good one without reaching for another dependency anyway (`gum`, `dialog`, ...).
- **Shell out from `tui/` to `scripts/doctor`/`scripts/status` and parse their text output**:
  rejected — their output is meant for a human terminal, not a stable machine format; parsing it
  back out is fragile and would silently drift the moment either script's wording changes without
  the parser being updated in lockstep. Reading the same underlying files natively is no more work
  and doesn't have this failure mode.
- **A different language with a smaller footprint** (e.g. a POSIX-shell `dialog`/`whiptail`-based
  UI): rejected — Go was picked to match LazySkills' own likely stack (single static binary, easy
  cross-platform `go build`, no runtime to install separately), and Bubble Tea is the de facto
  standard for exactly this kind of terminal dashboard.

## Revisit triggers
- The Go implementation and the shell scripts' logic actually diverge in practice (a check passes
  in one and fails in the other) — same signal ADR 0005 already named for its own two hooks; when
  it happens for real, consider generating one from the other instead of hand-syncing three places.
- Someone wants write actions from the dashboard (install an adapter, delete a stale symlink,
  bootstrap from inside it) — that's a deliberate scope change away from "read-only viewer," not a
  small addition; it would need its own discussion, not a quiet feature creep into `tui/`.
- Go stops being a reasonable assumption for this project's actual users (e.g. adoption skews
  heavily toward teams with no Go tooling at all) — reconsider distributing prebuilt binaries via
  `scripts/init`-style install instead of "build it yourself."
