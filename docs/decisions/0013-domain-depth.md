# ADR 0013 — Bootstrap asks "domain depth" before adopting an existing codebase

- Status: Accepted
- Date: 2026-09-23

## Context
The human flagged, from real adoption use (`adplus-dsp-dashboard`), that bootstrap's business-logic
docs (`docs/domain.md`, `docs/architecture.md`) were coming out shallower than expected. The cause:
`workflows/bootstrap.md`'s INSPECT step already reads the code for **stack/convention** detection
("detect stack(s), build system, test setup, and existing conventions *from the code*"), but the
INTERVIEW step fills `docs/domain.md`'s business rules and ubiquitous-language table purely from
what the human says out loud — the agent never commits to actually reading the domain/service/
handler code itself unless asked to. For a large or unfamiliar codebase, the human doesn't know or
won't think to narrate every rule, so the doc ends up only as complete as the conversation was.

The fix isn't "always read deeply" — a genuine deep read (opening and reasoning about real
business-logic files, not just detecting a stack) costs real time and tokens, and for a small or
low-stakes adoption the interview alone is plenty. It also isn't "never" — that's the status quo
that prompted this ADR. The right lever, consistent with every other bootstrap question, is to ask,
with a recommendation (Proposal rule), and let the human decide the cost/thoroughness trade for
their specific project.

## Decision
Add a **domain depth** question to INSPECT, adoption flow only (skipped for greenfield — there's no
existing business logic to read): **deep** (the agent reads the actual business-logic code and
drafts `docs/domain.md`/`docs/architecture.md` from it, for the human to correct against their own
answers rather than write from scratch) vs **shallow** (those two docs come from the INTERVIEW
alone, today's behavior, unchanged). The agent recommends deep for a non-trivial codebase, a
`strict`-leaning mode, or anything revenue/compliance-sensitive flagged during INSPECT; shallow for
something small or low-stakes — same recommend-then-let-the-human-decide shape as every other
bootstrap question.

Deep mode is explicitly routed as Clarify/Plan-tier judgment work: Claude Code delegates it to the
`planner` subagent (opus/high) rather than doing it inline in whichever session is running
bootstrap, matching the existing model-routing table's own reasoning (`adapters/claude-code/README.md`).
Codex has no equivalent subagent tier, so its skill just states the expectation in prose.

`docs/domain.md` and `docs/architecture.md` each get a one-line provenance note recording which
path was taken ("drafted from code, corrected against interview" vs "from interview only"), so a
future session — or a human three months later — isn't silently guessing how much to trust them.

## Consequences
**Buys us:** business-logic docs that can actually be thorough when the project needs it, without
making every bootstrap (including small/greenfield-adjacent adoptions) pay the cost of a deep code
read by default. The provenance note also means a shallow doc's limits are documented, not hidden —
someone deciding whether to trust `docs/domain.md` for a high-stakes change can check which mode
produced it instead of assuming.

**Costs us:** one more INTERVIEW-adjacent question (though only for adoption, not greenfield), and
a genuine token/time cost when deep is chosen — that's the point, not a side effect, but it means
bootstrap on a large existing codebase is no longer uniformly fast.

## Alternatives considered
- **Always deep for adoption flow:** rejected — makes every existing-codebase bootstrap pay full
  cost regardless of project size or stakes; the human asked specifically because the current
  behavior felt too shallow for *some* cases, not because they wanted maximum depth unconditionally
  every time.
- **Leave it as an unstated agent judgment call** ("use your best judgment on how deep to go"):
  rejected — violates the Proposal rule (AGENTS.md invariant 6) and the same reasoning ADR 0011
  already used for chat language: a project-wide, cost-relevant default belongs to an explicit
  human answer recorded somewhere findable, not a silent per-session guess that can vary run to run.
- **A third "medium" tier** (skim key files, don't do a full read): rejected for v1 as unnecessary
  granularity — a binary choice with a clear recommendation is easier to reason about than three;
  revisit if deep/shallow turns out too coarse in practice (see below).

## Revisit triggers
- The deep/shallow binary turns out too coarse in real use (people consistently want "deep, but
  only these three modules") — consider letting the human scope which paths to read deeply, instead
  of an all-or-nothing choice.
- Codex CLI grows a subagent/sub-session delegation primitive comparable to Claude Code's — revisit
  routing its deep-mode read through that instead of leaving it as same-session prose instruction.
