# Workflows

A workflow is the *process*: steps, roles, gates, and evidence. Each step points to the prompt it
uses (`prompts/…`) — workflows never copy prompt text, so there is exactly one source of truth.

| Workflow | Use when |
|---|---|
| `bootstrap.md` | First run: adapt this workspace to your project (new or existing). |
| `feature-development.md` | Any new behavior. |
| `bug-fix.md` | Something works incorrectly. Reproduction comes before the fix. |
| `refactor.md` | Structure changes, behavior doesn't. |
| `incident.md` | Production is on fire. Stabilize first, learn after. |

## Legend

- **[GATE: human]** — a human decision. Never automated, in any mode.
- **[GATE: human — strict]** — required in strict mode; in lite mode the agent states its
  assumption and proceeds unless you object.
- **Role:** who executes the step (see `docs/roles/`). In strict mode, role = separate session.
- **Evidence:** what must exist before the step counts as done.

## Modes

Set at bootstrap, recorded in `AGENTS.md`:

- **lite** — Spec → Plan → **[GATE: human]** → Build → Independent Review → Verify.
  For solo developers and low-risk work. Independent review is still non-negotiable — it's cheap
  (a subagent or second window) and catches what the producer can't.
- **strict** — Intent → Clarify → Spec → **[GATE]** → Plan → **[GATE]** → Build → Independent
  Review → **[GATE: triage]** → Verify → **[GATE: ship]**. For teams and critical systems.

If things go wrong at any step, don't improvise: `prompts/recovery/README.md` has the ramp.

## Trivial changes (the one spec exception)

A change may skip the spec entirely — no `specs/active/` entry, no plan, no role assumption; just
make the edit, referencing this line in the commit — only if **all** of these hold:

- Single file, no behavior change: a typo, a comment, formatting, a broken link, a version bump
  with no code impact.
- Zero risk: nothing that touches auth, data, money, or a forbidden dependency
  (`docs/architecture.md`).
- `scripts/check` still runs and is green afterward.

If any of these is unclear or debatable, it is **not** trivial — write the spec. Default to the
heavier path when in doubt (same instinct as R-10: ambiguity is never resolved by assumption).

## Epics — grouping specs toward one shared outcome

When one coherent outcome needs several specs (not one spec that's grown too big — several
genuinely separate specs that only make sense together), open `specs/epics/NNNN-<name>.md` from
`specs/epics/TEMPLATE.md`: a one-line shared goal and a table of the specs in it. It is not a PRD
or an architecture doc — if its Notes section starts accumulating real design decisions, that
content belongs in a spec or an ADR, not the epic file. Each spec still gets its own full spine
(Spec → Plan → Build → Review → Verify → Ship) independently; the epic is a pointer, not a
shortcut past any of that. `./scripts/status` rolls up every active spec (and any epic it belongs
to) into one view.
