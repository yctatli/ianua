# Prompt: Bootstrap

```
Before anything else: if the current working directory is itself named `.ianua`, STOP — this is
Ianua's own template (likely cloned here as a nested install), not my actual project. Tell me to
`cd ..` and re-launch you from the real project root instead; do not bootstrap `.ianua/` itself.

This workspace's own process files (AGENTS.md, docs/roles/, docs/decisions/, workflows/, prompts/,
scripts/, adapters/, and any installed .claude/.codex/.agents/) are protected by a core-file-lock
hook (AGENTS.md rule 7). If IANUA_ALLOW_CORE_EDIT is not already set, tell me and wait — don't try
to route around it.

Read AGENTS.md and workflows/bootstrap.md. We are adapting this workspace to a real project.

1. INSPECT the working tree. If it contains application code, detect stack, build/test tooling and
   existing conventions from the code, and present what you found for confirmation. If empty,
   we start from my answers.

   If it contains application code, also ask me — with your own recommendation — about DOMAIN
   DEPTH: **deep** (you read the actual business-logic code — services, handlers, domain/model
   layers — yourself, and draft docs/domain.md + docs/architecture.md from it, for me to correct)
   vs **shallow** (those two docs come only from what I tell you in the INTERVIEW below). Recommend
   deep for a non-trivial codebase, a `strict`-leaning mode, or anything revenue/compliance-
   sensitive; shallow is fine for something small/low-stakes. This is Clarify/Plan-tier judgment,
   not transcription — if you're the `planner` subagent, that's already you; if you're the main
   session, delegate the deep read to `planner` rather than doing it inline.

2. INTERVIEW me one topic at a time — with your own recommendation and rationale for every
   question (proposal rule):
   product & goal · domain terms and business rules · architecture style & module boundaries ·
   forbidden dependencies · conventions (data rules, error handling, naming) · testing
   expectations · security posture (INCLUDING which dependency/CVE audit and SAST tool fits this
   stack — recommend one, don't leave docs/security.md's "Automated checks" unfilled) · git rules ·
   operating mode (lite or strict — recommend one based on team size and risk) · chat language
   (what language you should respond in — recommend matching the language I'm writing this
   interview in; code/docs/commits stay English by default unless I say otherwise).
3. GENERATE from my answers: fill every template in docs/ (architecture, domain, conventions,
   testing, security, git) — docs/domain.md and docs/architecture.md per the domain-depth answer
   from step 1, each starting with a one-line provenance note ("drafted from code, corrected
   against interview" or "from interview only"); write scripts/check.conf with real build/test/lint
   commands for this stack PLUS a `security:` step (dependency audit and/or SAST — see
   scripts/check.conf.example); set the Mode and Chat language lines in AGENTS.md; rewrite the
   AGENTS.md project summary. HARD RULES: the "Invariant rules" block in AGENTS.md is kept
   verbatim; AGENTS.md stays a signpost ≤ 40 lines — it POINTS to docs, it never copies them.
4. VERIFY: run ./scripts/doctor and ./scripts/check and show me the output.
5. REPORT: what you generated (including which domain depth you used and why), what you assumed,
   and every open question that still needs my decision — with your recommendations.
```
