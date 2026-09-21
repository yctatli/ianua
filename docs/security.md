# Security

> **Template — filled during bootstrap.** Baseline rules agents must honor in every plan and review.

## Secrets
- Secrets never enter the repo, specs, prompts, or chat. `.env` is gitignored; provide `.env.example`.
- Agents never print secret values, even when debugging.

## Input & output
<!-- Validation strategy at boundaries; what is escaped/encoded; what never leaks in error messages. -->

## AuthN / AuthZ
<!-- Who may call what; where authorization is enforced; default-deny statements. -->

## Dependencies
<!-- Policy for adding new dependencies: who approves, what gets checked (license, maintenance, CVEs). -->

## Automated checks (tooling, not prose)
<!-- Bootstrap fills this with real commands, stack-detected or chosen at interview time, and adds
them as a `security:` step in scripts/check.conf — so every `scripts/check` run is also a security
gate, not just something a reviewer eyeballs. Two kinds, add what applies:
- Dependency/CVE audit: `npm audit`, `pip-audit`, `dotnet list package --vulnerable`, `bundle audit`...
- Static analysis (SAST): `semgrep --config auto`, `bandit -r .`, language-specific linters with
  security rule sets enabled.
A step that isn't wired into scripts/check.conf is advice an agent can silently skip; wire it in. -->

## Review lens
Security is a mandatory dimension of every independent review (see `prompts/review.md`), not a
separate afterthought phase — automated checks catch known patterns, the reviewer catches
design-level issues (authz gaps, trust-boundary mistakes) no scanner will find. For a systematic
pass beyond the reviewer's own read, Claude Code users can run the built-in `security-review`
skill against the change set — see `adapters/claude-code/README.md`.
