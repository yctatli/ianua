# Adapter: Codex CLI

Install: `./scripts/init codex` — copies `.agents/skills/`, `.codex/hooks.json`, `.codex/hooks/`,
and `.codex/config.toml` to the repo root.

Codex CLI reads **AGENTS.md natively** (repo root down to cwd, plus your own `~/.codex/AGENTS.md`)
— no pointer file needed, unlike the Copilot/Cursor adapters. What this adapter adds on top:

- **Skills** (`.agents/skills/`): `bootstrap`, `new-feature`, `fix-bug`, `refactor`, `review`,
  `security`, `verify`, `adr`, `recover` — one per workflow. Invoke explicitly with `$<name>`, or
  let Codex match one from your prompt (implicit invocation, via each skill's `description`). Thin
  by design: they point to `workflows/` and `prompts/`, they don't restate them. `security` is a
  genuinely independent pass (`docs/roles/security.md`) — mandatory in strict mode alongside
  `review`, optional/on-demand in lite mode — not a duplicate of `review`'s own security bullet.
- **PreToolUse hook** (`.codex/hooks/policy.py`, wired in `.codex/hooks.json`): blocks writes
  under `specs/done/` and a handful of destructive git ops (force push, hard reset, rebase,
  `rm -rf`) — the same rules the Claude Code adapter enforces, ported to Codex's hook contract.
- **Suggested project config** (`.codex/config.toml`): `sandbox_mode = "workspace-write"`,
  `approval_policy = "on-request"`.

## Trust — read this or the hook will silently not run

Codex only loads project-scoped `.codex/` files (config, hooks, skills) once the project is marked
**trusted**. The first time you run `codex` in this repo it will ask; or set it yourself in your
**user-level** `~/.codex/config.toml`:

```toml
[projects."/absolute/path/to/this/repo"]
trust_level = "trusted"
```

An untrusted project skips `.codex/config.toml`, `.codex/hooks.json`, and project-scoped rules
entirely — no error, no warning, the policy hook just never fires.

## Independent review, mechanically

Skills have no hard read-only mode the way Claude Code's `reviewer` subagent (`tools:` allowlist)
does — the `review` skill's "don't write" instruction is prose the model follows, not something
the tool enforces. For a real guarantee, start the REVIEW step in a **separate Codex session**
(new terminal / `codex` invocation, never the session that built the change — see
`docs/roles/README.md`, "the mind that produces cannot audit itself") launched with the read-only
sandbox: `sandbox_mode = "read-only"`, set via CLI flag, `config.toml`, or the in-session
`/permissions` command.

## A note on drift

Skills, and the PreToolUse hook contract they're wired through, are a newer and fast-moving part
of Codex CLI. If `$review` doesn't trigger, or the policy hook doesn't fire, check the directory
path (`.agents/skills/<name>/SKILL.md`) and the JSON shape in `.codex/hooks/policy.py` against
your installed Codex CLI version — OpenAI's own docs are the source of truth; this adapter is a
snapshot.
