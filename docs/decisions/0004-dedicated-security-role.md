# ADR 0004 — Add a dedicated Security role, mode-graduated

- Status: Accepted
- Date: 2026-09-21
- Supersedes: the "no dedicated Security role" alternative in
  `docs/decisions/0003-security-enforcement.md` (that ADR's other decisions — the `security:`
  check.conf step and the skill invocation inside `reviewer` — still stand and are incorporated here).

## Context
ADR 0003 considered and rejected a dedicated Security role, reasoning that `docs/security.md`
already committed to "not a separate afterthought phase" and that a bolt-on phase is the first
thing dropped under deadline pressure. On reflection (prompted by the human questioning that call),
that reasoning conflated two different things: *security shouldn't be a phase that only happens at
the end*, which is true and still holds, and *security can't be its own role*, which doesn't follow
— Reviewer, QA, and Developer are all "separate phases" too, and nobody considers that a problem,
because `docs/roles/README.md` already has the actual answer to "won't this get skipped under
pressure": **mode-graduated mandatory-ness**. Strict mode makes role separation mandatory; lite
mode requires only the independent review and lets the rest collapse into one session. That
existing mechanism is precisely the tool needed here and wasn't applied.

A second, concrete signal: `reviewer` already needed to invoke a dedicated tool (the built-in
`security-review` skill) just to cover its security bullet properly. A review dimension that needs
its own tool to do well is already halfway to deserving its own role.

## Decision
Add `docs/roles/security.md` (fifth role card, same structure as the other four) and
`prompts/security-review.md` (dedicated prompt, parallel to `prompts/verify.md` for QA). Wire it
mode-graduated, exactly like the rest of role separation:
- **Strict mode**: Security runs as its own independent session/pass, alongside Reviewer, for
  every REVIEW step in `feature-development.md`, and for security-sensitive changes in
  `bug-fix.md`/`refactor.md`. Both must return "clean" or triaged findings.
- **Lite mode**: unchanged from ADR 0003 — Security collapses into Reviewer's own security
  dimension by default; running it separately is optional, recommended for anything touching auth,
  input handling, data access, or dependencies.

Claude Code: new `security` subagent, read-only (no Edit/Write, same guarantee as `reviewer`),
pinned to **opus / high effort** — same tier as `reviewer`, same reasoning (see
`docs/decisions/0002-claude-code-model-effort-routing.md`'s table: low call volume, high cost of a
miss). New `/security` command. `reviewer` keeps its `security-review` skill invocation as the
lite-mode fallback, now explicitly framed as a cross-check (not the only pass) when `security` also
ran.

Codex: new `security` skill mirroring `review`, same mode-graduated framing.

`docs/roles/README.md` and `docs/security.md` updated to state the mode split explicitly, so it's
one documented rule instead of something each workflow file re-derives.

## Consequences
**Buys us:** security gets the same structural guarantee correctness already has in strict mode —
an independent session that can't be talked out of a finding by the builder's rationale — instead
of living or dying by how much attention one reviewer happens to give one bullet in a six-item list.
Lite-mode projects pay nothing extra by default; the option is there without being forced. The
"security-review needs its own tool" signal from ADR 0003 is now honored structurally (a role with
its own prompt and subagent) instead of awkwardly bolted onto Reviewer's instructions.

**Costs us:** a fifth role card and prompt to keep in sync with the other four if the process
changes; in strict mode, one more independent session per feature (real time/token cost — mirrors
the cost QA and Reviewer already add, not a new category of cost). Three workflow files now
mention Security explicitly, so future workflow edits need to remember it exists, same maintenance
tax role separation already carries for Reviewer/QA.

## Alternatives considered
- **Keep ADR 0003's position** (security as a lens only, no role): rejected on reflection — see
  Context. The deadline-pressure risk this was meant to guard against is better handled by
  lite-mode collapse than by refusing to let security be a role at all.
- **Make Security mandatory in lite mode too**: rejected — would raise the cost floor for every
  small lite-mode change regardless of risk, contradicting lite mode's whole premise ("minimum
  viable ceremony"). Recommending it for security-sensitive lite changes, without forcing it for
  everything, keeps the graduation meaningful.
- **Merge Security into QA instead of Reviewer's peer**: rejected — QA verifies against already-
  written acceptance criteria with evidence; Security is an adversarial audit like Reviewer's, not
  an evidence-mapping exercise. It belongs next to Reviewer, not QA.

## Revisit triggers
- A strict-mode feature ships with a real security finding that Security's pass missed but a
  generalist reviewer would plausibly have caught anyway — re-examine whether the split pass is
  pulling its weight versus just being an extra gate.
- Lite-mode projects never opt into the optional Security pass even for auth/data-touching
  changes — the "recommended but not mandatory" framing may need to become a stronger nudge
  (e.g. a `scripts/doctor` warning when a diff touches auth-looking paths without a Security note).
- Same drift risk as ADR 0002/0003: if Anthropic changes subagent effort support or the
  `security-review` skill's behavior, re-verify before assuming this still holds.
