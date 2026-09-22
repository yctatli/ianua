# ADR 0008 — Selective adoption from BMAD-METHOD: status rollup, trivial-change exception, epics, lesson

- Status: Accepted
- Date: 2026-09-22

## Context
The human pointed at BMAD-METHOD (bmad-code-org/BMAD-METHOD on GitHub) and asked what's genuinely
good in it, given they explicitly do not want a system that complex. Reviewing its current (v6)
docs found real, verified mechanisms: five named personas (Mary/Analyst, John/PM, Winston/Architect,
Amelia/Dev, Sally/UX) with numbered command menus; a three-tier planning-depth router (Quick/Epic/
Full Project, chosen by scope and risk, not a fixed process); story files sharded from a PRD and
architecture doc so "a developer can implement these epics without inventing decisions nothing
records"; a `sprint-status.yaml` rollup aggregating every story's status, risk flags, and open
retrospective action items; and a retrospective step that appends learnings back into that backlog.

Most of this doesn't fit ANEW: the personas are a stylistic device with no functional payoff here;
per-role checklist files would duplicate what role cards' "Powers & prohibitions" and the spec's
own DoD already do; a hand/agent-maintained YAML tracking file is a second source of truth that can
drift from the specs it's supposed to describe — exactly the kind of duplication ANEW's "one
verification contract, files not chat" design avoids elsewhere.

Four ideas survived that filter as genuinely valuable AND cheap to add without new roles, new file
types beyond what already exists in spirit, or a second source of truth:
1. A rollup view across in-flight work (the value behind `sprint-status.yaml`).
2. Adaptive planning depth for genuinely trivial work — BMAD's Quick Path names a real gap:
   ANEW's rule 1 ("no spec, no code") has no floor, so a one-line typo fix pays the same ceremony
   tax as a real feature.
3. A lightweight grouping mechanism for multi-spec initiatives (BMAD's Epic Path, without the
   PRD/architecture-doc apparatus Full Project brings with it).
4. Closing the loop from friction back into a concrete improvement (BMAD's retrospective, without
   a dedicated role or session).

## Decision
Four additions, each deliberately smaller than its BMAD counterpart:

1. **`scripts/status`** — a new shell script, same family as `scripts/check`/`scripts/doctor`.
   Scans `specs/active/*.md` and `specs/epics/*.md` and reports Status, Mode, DoD checked/total,
   whether the security DoD line is checked, and the linked plan's status — all *derived*, nothing
   hand-maintained, so it cannot drift out of sync with the specs it summarizes the way a parallel
   YAML ledger could.
2. **Trivial-change exception** (`workflows/README.md`) — a narrowly scoped carve-out from rule 1:
   single-file, zero-behavior-change, zero-risk edits (typo, comment, formatting) may skip the spec
   and plan entirely, no role assumption needed. Any ambiguity about whether something qualifies
   defaults to writing the spec, same instinct as R-10.
3. **`specs/epics/`** (`README.md` + `TEMPLATE.md`) — a one-line shared goal plus a table of the
   specs that belong to it. Each spec inside still runs its full independent spine; the epic is
   explicitly documented as "a pointer, not a shortcut," with its own Notes section warned against
   growing into a de facto PRD (that content belongs in a spec or an ADR instead).
4. **Lesson section** (`specs/TEMPLATE.md`, after the Scorecard) — one optional paragraph, filled
   only when the Scorecard shows real friction (>1 spec revision, >2 fix rounds, a finding that
   shouldn't have reached BUILD): name one concrete preventive change. Referenced from the SHIP
   steps of `feature-development.md` and `bug-fix.md`.

`scripts/status`, and the `TEMPLATE.md` files it depends on, are added to the core-file-lock
protected set (`docs/decisions/0005-core-file-lock.md`) — same tier as `scripts/check`/`doctor`/
`init`. `specs/epics/` itself is treated like `specs/active/`/`specs/plans/` (open, not locked) —
only its `TEMPLATE.md` is protected, matching how the existing spec/plan templates are handled.

## Consequences
**Buys us:** a real answer to "what's in flight and how healthy is it" without a second file to
keep in sync; a documented, narrow relief valve for ceremony-disproportionate trivial work instead
of that pressure building up as quiet non-compliance; a way to talk about several related specs as
one initiative without inventing a PRD layer; a mechanism that turns "that was harder than it
should've been" into something that actually changes the next spec, instead of evaporating.

**Costs us:** one more script to maintain in lockstep with `specs/TEMPLATE.md`'s exact field names
(`scripts/status`'s parsing is line-pattern-based, same fragility class as the core-file-lock
hooks' own grep/sed parsing — documented, not hidden); a small but real erosion surface on rule 1 —
"trivial" is a judgment call, and a bad-faith or careless "this is basically trivial" could smuggle
real work past the spine. The criteria are written narrow and default-to-heavier specifically to
bound that risk, but it isn't zero.

## Alternatives considered
- **BMAD's full three-tier system (Quick/Epic/Full Project) as a formal, named routing step**:
  rejected — the human was explicit about not wanting this complexity. The trivial-change exception
  covers the one end (Quick) that had no ANEW equivalent at all; Epic's grouping need is covered by
  `specs/epics/`; Full Project's PRD/architecture-doc apparatus has no gap to fill — ANEW's
  `docs/architecture.md`+`docs/domain.md`, filled once at bootstrap, already serve that role.
- **A hand-maintained tracking file (`status.yaml` or similar), matching BMAD's `sprint-status.yaml`
  more literally**: rejected — a second source of truth that must be kept in sync by discipline
  (agent or human) is exactly the class of thing this workspace avoids elsewhere ("evidence over
  claims," "one verification contract"). A derived report has no sync problem because it has
  nothing to sync.
- **Named personas, per-role checklists, an elicitation-menu mechanism**: rejected outright — no
  identified functional gap they'd close that isn't already covered (role cards' own prohibitions
  and DoD checklists for the first two; the Proposal Rule for the third), and each adds ceremony or
  files for their own sake, which is precisely what the human asked to avoid.

## Revisit triggers
- The trivial-change exception gets invoked for something that, in hindsight, wasn't trivial —
  tighten the criteria in `workflows/README.md` rather than removing the exception outright.
- `scripts/status`'s line-based parsing breaks because `specs/TEMPLATE.md`'s field names or DoD
  structure change — update both together; they're coupled by design, the same way the core-file-
  lock hooks' CORE_PATTERNS lists are coupled across the two adapters.
- An epic's Notes section is observed growing into real design content — that's the signal to split
  it into a spec or an ADR, not a signal the epic template needs a bigger Notes section.
- The Lesson section is skipped even when the Scorecard clearly shows friction — consider making it
  a DoD checklist item instead of a conditional note, at the cost of a small amount of ceremony on
  every spec instead of only the ones that need it.
