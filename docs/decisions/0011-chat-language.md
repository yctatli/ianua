# ADR 0011 — Bootstrap asks for a chat language, recorded in AGENTS.md

- Status: Accepted
- Date: 2026-09-22

## Context
The human asked for AI responses to default to a chosen language (Turkish, in the request that
triggered this) instead of Ianua silently assuming English everywhere. Since `AGENTS.md` is the one
file every adapter auto-loads at session start (see README, "Why this exists" §1), it's already the
right place to record a per-project default that every future session — any tool, any model —
should pick up without being told again each time.

Two shapes were on the table: hardcode a language somewhere in the template, or treat it like every
other project-specific fact Ianua already asks about (domain terms, conventions, mode) — interview
it, with a recommendation, and record the answer. Hardcoding would fight Ianua's own stated design
(general-purpose, not tuned to one team or one language) and the Proposal rule (every question
comes with a recommendation, the human still decides).

Worth keeping separate: *chat* language (what the agent speaks to the human in) is not the same
axis as *code/docs/commit* language, which most teams keep in English regardless of what language
they talk to their tools in, for the same reason most non-English-speaking engineering orgs already
write code comments in English — searchability, hiring, tooling. Conflating the two would force an
awkward choice neither this ADR nor the interview needs to force.

## Decision
Add "chat language" as one more INTERVIEW topic in `workflows/bootstrap.md` and `prompts/bootstrap.md`,
recommended default = whatever language the human is writing the interview in, explicit only if they
want something else. GENERATE writes the answer as a second line under `AGENTS.md`'s existing
"Operating mode" section (next to the `Mode:` line it already sets), scoped explicitly to chat only
— code/docs/commits stay English unless the human says otherwise in that same answer.
`scripts/doctor` gets a matching "chat language not set" warning, mirroring the existing
mode-not-set check, so an unconfigured workspace surfaces this the same way it already surfaces an
unset mode.

## Consequences
**Buys us:** one interview question, one AGENTS.md line, and every future session (regardless of
which AI tool or model) in that project picks it up automatically — no per-session instruction, no
per-user memory hack, and it travels with the repo the way mode already does.

**Costs us:** two more lines against the "AGENTS.md ≤ 40 lines" budget (already a soft target, not
a hard limit — `scripts/doctor` only warns past 60). Negligible.

## Alternatives considered
- **Hardcode a default language in the template:** rejected — contradicts "general-purpose,
  technology- and team-agnostic," the same reasoning that already keeps stack packs out of core.
- **Put it in assistant-side memory instead of AGENTS.md:** rejected — README's own "Memory" section
  already draws this line: assistant memory is per-user/per-tool and invisible to a teammate or a
  different AI tool reading the same repo; a project-wide default belongs in the repo, not there.
- **One combined "language" setting for chat AND code/docs/commits:** rejected — most non-English
  teams want English code/docs for the same searchability/hiring reasons everyone else does, while
  still wanting to talk to the agent in their own language; forcing one setting for both would make
  the common case awkward.

## Revisit triggers
- A team wants per-file-type language control beyond the chat/code split here (e.g. Turkish docs,
  English chat) — extend to a small table instead of one line, only if this is actually asked for.
- The chat-language line drifts out of sync with what the human actually wants mid-project (they
  start writing in a different language than the recorded default) — should probably just be a
  quick `/adr`-free AGENTS.md edit under `IANUA_ALLOW_CORE_EDIT=1`, not a process problem to solve.
