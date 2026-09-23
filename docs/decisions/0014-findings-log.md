# ADR 0014 — Persistent Findings log on the plan, not just chat

- Status: Accepted
- Date: 2026-09-23

## Context
The human, looking at a real multi-round review session (`adplus-dsp-dashboard`, spec 0001 —
findings labeled F-A through F-J across three review rounds, background subagents launched and
re-launched to fix them), asked why there was no itemized, checklist-style record of which findings
were fixed, deferred, or rejected — only narrative chat (background-agent notifications, a final
prose recap).

Checking: `scripts/status` and the spec's own Definition-of-Done checkboxes already give a
persistent, file-derived answer to "did review happen, is the work done" — that part was never the
gap. The actual gap is one level down: `prompts/review.md` and `prompts/security-review.md` tell
the reviewer to report findings with evidence, and TRIAGE (`workflows/feature-development.md` step
7, `workflows/refactor.md`'s REVIEW step) tells the human to decide real/noise/investigate per
finding — but nowhere does any workflow or prompt say to **write** the finding + its disposition to
a file. The reviewer/security subagents are read-only by design (`docs/decisions/0005-core-file-lock.md`'s
spirit, `adapters/claude-code/README.md` tool allowlists) — they structurally cannot write it even
if asked. So today, the only trace of "F-A was a real bug, F-B was noise, F-C got deferred to a
later step" is the conversation itself, and the eventual Scorecard only records a *count* ("Review
findings: real / noise") at SHIP, not the itemized list. This is a straightforward instance of
Ianua's own stated principle *"Context lives in files, not in chat"* (README, "Design principles")
not actually being followed for this specific piece of context.

## Decision
Add a **Findings log** table to `specs/plans/TEMPLATE.md`: `ID | Round | Source | Severity |
Finding | Status | Resolution`, appended-only, never deleted. Recording it is explicitly the
orchestrating/main session's job (it has write access; the reviewer doesn't) — `prompts/review.md`
and `prompts/security-review.md` now ask the reviewer to give each finding a short ID (F-1, F-2, ...
/ S-1, S-2, ...) specifically so the human/orchestrator can reference it 1:1 when writing the row.
`workflows/feature-development.md` step 7 (TRIAGE) and step 8 (FIX ROUNDS), `workflows/refactor.md`'s
REVIEW step, and `workflows/bug-fix.md`'s REVIEW step (when a fix escalated to a mini-spec with a
plan) all now say to record/update it as part of the gate that was already happening — this is not
a new gate, just a new place the existing triage decision gets written down.

`scripts/status` gets a `findings_summary()` function that rolls each plan's table into one line
(`findings: 2 open, 5 fixed, 1 noise, 1 deferred`) per spec, the same derived-from-files spirit as
everything else it already reports. `tui/specs.go` re-implements the identical parsing (ADR 0012's
already-accepted cost of a second implementation kept in sync by hand) and surfaces it in the
Status tab, highlighted when anything is still Open.

Scoped to the **plan** file, not the spec: feature-development and refactor both always have a
plan with multi-round review cycles where this actually pays for itself; a plain bug-fix (no plan
file, one short review cycle) doesn't need the ceremony — its fix diff is already the record,
matching the existing proportional-ceremony precedent (trivial-change exception, lite/strict modes).

## Consequences
**Buys us:** the exact "madde madde" (itemized) checklist view the human was looking for, answerable
from a file at any point — including after a context compaction or a brand new session — instead of
only reconstructable by scrolling chat history. `scripts/status`/`tui/` now show at a glance whether
a spec still has open findings, not just whether DoD boxes are checked (a spec could have all DoD
boxes checked while an Open finding sits unresolved if nobody's watching — this makes that visible).

**Costs us:** one more thing the orchestrating session must remember to write during TRIAGE and FIX
ROUNDS — nothing stops it from being skipped, same prose-level-not-mechanical caveat the "Do the
rules actually get enforced?" README section already names for plan approval and other human gates.
`scripts/status` and `tui/specs.go` now both parse the same table shape — a third place the format
needs to be kept in sync if it ever changes (same tradeoff, now on a third file, as ADR 0012 already
flagged for the doctor/status logic specifically).

## Alternatives considered
- **A separate `specs/active/NNNN-findings.md` file instead of a plan section**: rejected — adds a
  new file class and a new required-structure check to `scripts/doctor`/`tui/health.go` for
  something that's naturally plan-scoped (findings are about reviewing a plan's implementation);
  keeping it as a section in a file that already exists is less ceremony for the same result.
- **Let the reviewer subagent write it directly**: rejected outright — reviewer/security are
  read-only by tool allowlist on purpose (ADR 0005's own reasoning: the producer/auditor boundary
  should be a structural guarantee, not a request); making an exception here to save one step would
  undermine the exact property that makes "independent review" a hard guarantee instead of a prose
  one (see the README's "Do the rules actually get enforced?").
- **Only track counts (real/noise), not itemized rows**: rejected — that's the status quo (the
  Scorecard already does this) and is precisely what prompted the question; a count answers "how
  many," not "which ones, and what happened to them."

## Revisit triggers
- `scripts/status`'s awk parsing and `tui/specs.go`'s Go parsing of the Findings log table actually
  diverge in practice — same signal ADR 0005 and ADR 0012 already named for their own two-
  implementation pairs; consider generating one from the other, or a shared fixture test, if it
  happens for real instead of staying hypothetical.
- A bug-fix workflow with no plan file turns out to also want itemized review-finding tracking in
  practice (multiple review rounds on a fix, not just one) — extend the Findings log convention to
  live in the spec itself for that case, rather than assuming a plan always exists.
