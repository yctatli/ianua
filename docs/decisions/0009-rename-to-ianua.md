# ADR 0009 — Rename the project: ANEW → Ianua

- Status: Accepted
- Date: 2026-09-22

## Context
The workspace was checked out under a folder name (`roombook-main`) that had nothing to do with
what it actually is, and the project itself had never chosen a real identity beyond the backronym
"ANEW" (AI Native Engineering Workspace) used in the README title. The human asked for a real,
deliberate name, explicitly wanting something creative with a Latin root rather than another
generic tech-sounding word.

Brainstormed against the system's own actual vocabulary — the single most repeated structural
concept across every workflow is the **gate** (`[GATE: human]`, `[GATE: human — strict]`, a gate at
nearly every step of every workflow). **Ianua** is Latin for "door, gateway, threshold" — the root
of Janus, the Roman god of doorways and transitions, depicted looking simultaneously backward and
forward. That double-facing is an apt, unforced metaphor for what this workspace actually does:
ADRs and `specs/done/` look backward (why a decision was made, what shipped and how it was proven);
specs, plans, and epics look forward (what's being built and why). Other candidates considered —
Custos (guardian), Indicium (evidence), Pactum (covenant/contract), Gradus (a step) — each mapped to
one real mechanism in the system; Ianua was chosen because "gate" is the one concept that recurs at
literally every level (workflow gates, the core-file-lock hook, review/security as gates on ship).

## Decision
Rename the project's branding from "ANEW" to "Ianua" everywhere it's forward-facing:
- README.md title, prose ("Ianua gives you...", "Every file in Ianua...", "Ianua distills the
  methodology...").
- Script header comments and echoed banners (`scripts/check`, `doctor`, `init`, `status`).
- Both adapters' README prose and the `bootstrap` skill's description.
- The `adapters/cursor/rules/anew.mdc` pointer file, renamed to `ianua.mdc` (and `scripts/init`'s
  cursor case updated to match).
- The core-file-lock environment variable, `ANEW_ALLOW_CORE_EDIT` → `IANUA_ALLOW_CORE_EDIT`, in
  both hook scripts' actual logic (not just comments) and every current doc that instructs a human
  to set it.

Deliberately **not** touched: the body text of already-accepted ADRs (0005, 0006, 0008) that
reference "ANEW" or `ANEW_ALLOW_CORE_EDIT` — per `docs/decisions/README.md`, decisions aren't
edited after acceptance except for a status update, and a rebrand is not a reason to make an
exception. Those ADRs are historically accurate: that was the project's name and that was the
variable's name when those decisions were made. A reader hitting `ANEW_ALLOW_CORE_EDIT` in ADR
0005's body and `IANUA_ALLOW_CORE_EDIT` in the live hook script should read that as history, not
as a bug — this ADR is the explanation for why the two differ.

The actual on-disk folder name and any GitHub repository name are outside this ADR's reach — a
running session cannot safely rename its own working directory or a remote repo; the human did the
folder rename directly (`roombook-main` → `ianua`) alongside this request.

## Consequences
**Buys us:** a name with real, checkable meaning instead of an arbitrary placeholder or a strained
backronym, consistently applied everywhere a human or agent actually reads it going forward.

**Costs us:** the two-name split in history (old ADRs say ANEW/`ANEW_ALLOW_CORE_EDIT`, everything
current says Ianua/`IANUA_ALLOW_CORE_EDIT`) requires this ADR to exist and be found — a reader who
never sees it could be confused by ADR 0005 instructing `ANEW_ALLOW_CORE_EDIT` when the hook no
longer checks that name. Mitigated by this ADR's own explicit callout, and by the fact that ADRs
are read in numeric order alongside the decisions log table in README.md, where 0009 sits right
after the ones it explains.

## Alternatives considered
- **Also rewrite the old ADRs' body text to say Ianua**: rejected — violates this workspace's own
  ADR-immutability convention for no real gain; the "history says the old name" outcome is correct,
  not a defect to fix.
- **Keep `ANEW_ALLOW_CORE_EDIT` as a legacy identifier, rename only prose**: considered, then
  rejected — leaving the env var name orphaned from the new brand for no functional reason would be
  a worse, permanent inconsistency than the one-time, well-documented history split this ADR
  accepts instead.
- **A non-Latin, more literal name (e.g. "Ledger", "Spine" — both discussed with the human)**:
  rejected per explicit request for a Latin-rooted, creative name; both remain fine alternatives if
  Ianua is revisited later.

## Revisit triggers
- Confusion in practice from the ADR-history name split (someone sets `ANEW_ALLOW_CORE_EDIT` from
  reading an old ADR and it silently does nothing) — consider having the hook scripts recognize
  both names for one transition period, logging a deprecation note on the old one, rather than
  relying solely on this ADR being read.
- The project is renamed again — write a new ADR the same way, and apply the same rule: rewrite
  forward-facing content, leave prior ADR bodies (including this one) as the historical record.
