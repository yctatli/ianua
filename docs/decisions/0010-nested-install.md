# ADR 0010 — Nested install: `.ianua/` as a nested clone for existing repos

- Status: Accepted
- Date: 2026-09-22

## Context
Quickstart said "drop the files into an existing repo," meaning: copy every Ianua file directly
into the target repo's root, alongside its own code. In practice, adopting Ianua into a real,
already-tracked repo this way clutters the diff with a whole second project's worth of files
(`docs/`, `specs/`, `workflows/`, `prompts/`, `scripts/`, `adapters/`) mixed into the consumer's
own tree, with no boundary between "this is Ianua" and "this is the actual product code" — and no
clean way to `git pull` Ianua updates later without re-copying by hand.

The instinct to isolate it under a single `.ianua/` folder is right, but it collides with two real
mechanical constraints, not just style preferences:

1. Claude Code only auto-loads `CLAUDE.md` and only auto-discovers `.claude/` (commands, agents,
   hooks) at the **project root** — not inside an arbitrary subfolder. Same story for Codex CLI and
   `AGENTS.md`/`.agents/`/`.codex/`. A `.ianua/`-only install is invisible to every adapter.
2. `scripts/check` resolves its own script directory as the execution root and `cd`s there before
   running `check.conf`'s build/test commands. If the script lives at `.ianua/scripts/check`, that
   root is `.ianua/`, not the consumer project — every build/test command fails immediately since
   the actual code (e.g. `go.mod`) isn't there.

Discovered concretely while adopting Ianua into `adplus-dsp-api`: bootstrap had already completed
correctly (there is no bug in bootstrap itself), but the intended `.ianua/`-only layout was never
actually reachable — nothing was nesting anything, the whole core had spread across the consumer
repo's root because that's the only place the adapters, and `scripts/check`'s command execution,
actually work today.

## Decision
Add a second, symlink-based install mode to `scripts/init`, auto-detected (no new flag): if the
script's own root directory is literally named `.ianua`, it treats itself as nested and installs
each adapter's discovery files as **symlinks** at the parent directory (the real project root)
instead of copies — `AGENTS.md`, `CLAUDE.md`, `.claude/`, `.agents/`, `.codex/`, `.cursor/rules/`,
`.github/copilot-instructions.md`, whichever a given adapter needs. `.ianua/` itself stays the only
place with real content; the project root only gets thin pointers into it, so there is nothing to
keep in sync by hand and nothing meaningful to review in a diff.

`scripts/check` and `scripts/doctor` split "core root" (where `check.conf`/specs/docs live — always
their own script directory) from "exec root" (where build/test commands actually run and where
adapter files land) — exec root is the parent directory when nested, otherwise the same as core
root. `scripts/init` also appends `.ianua` to the consumer repo's `.gitignore` if one exists and
doesn't already have it.

The core-file-lock hook (`protect-shipped.sh`, ADR 0005) gets the same path list a second time with
a `.ianua/` prefix, so files are protected under either layout without runtime detection — a plain
repo never has a `.ianua/` directory, so the extra patterns are inert there.

The classic (non-nested, "use as a GitHub template") install path is untouched — same copy-based
behavior as before. This is strictly additive.

## Consequences
**Buys us:** a real answer to "how do I use this in an existing repo" that doesn't scatter a
second project's files through the consumer's tree — one folder, one `.gitignore` line, symlinks
instead of a copy that silently drifts out of sync with the next `git pull` inside `.ianua/`.

**Costs us:** symlinks are a platform assumption (works on macOS/Linux; Windows needs Developer
Mode or admin rights for `ln -s`/`mklink` — not handled here, flagged as a revisit trigger).
Two install shapes now exist instead of one, so `scripts/init`, `scripts/check`, and `scripts/doctor`
each carry a small branch, and the hook's path list is now duplicated rather than computed —
deliberate, to keep the hook static and auditable instead of adding shell logic to a security-ish
script, at the cost of the two lists needing to be kept in sync by hand if the core-path set ever
changes (already true of the two per-tool hook scripts before this ADR, per ADR 0005's own revisit
trigger).

## Alternatives considered
- **Copy into `.ianua/` too (mirror the classic path 1:1, just nested):** rejected — doesn't fix
  either mechanical constraint above; adapters still wouldn't be discovered at the real project
  root, and `scripts/check` would still `cd` into `.ianua/` before running build commands.
- **`CLAUDE.md`/`AGENTS.md` content-templated per mode (rewrite the pointer text at install time
  instead of symlinking):** rejected — needs `scripts/init` to know and maintain two content
  variants per adapter file, and drifts the moment either variant is hand-edited later; a symlink
  needs no template at all and can never drift from its target.
- **Runtime nested-detection inside the hook** (check for a `.ianua/` directory at hook-run time
  and prepend the prefix dynamically) instead of a static doubled pattern list: considered simpler
  to read at a glance, rejected because it adds branching shell logic to the one script whose whole
  job is being a small, auditable gate — a static list stays literal and diffable.

## Revisit triggers
- A Windows-first user actually hits this (symlink creation fails or requires elevation) — add a
  copy-based fallback for nested mode specifically on that platform, gated on `ln -s` failing.
- The core-path list changes in one of the four places it now lives (either hook script, prefixed
  or not) without the other three being updated — same risk ADR 0005 already flagged, now doubled;
  if it causes a real miss, generate the four lists from one source instead of hand-copying.
- Someone wants Codex CLI and Claude Code adapters installed side by side in nested mode and hits a
  symlink collision on a shared target (e.g. `AGENTS.md`) — `link()` currently just overwrites
  existing symlinks silently on re-run, which is fine for re-running the *same* adapter but would
  silently repoint `AGENTS.md` if a second adapter's init also targets it; not currently a problem
  since both adapters' `AGENTS.md` symlink points at the same `.ianua/AGENTS.md`, but worth
  re-checking if that stops being true.
