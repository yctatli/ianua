# Prompt: Bootstrap

```
Before anything else: this workspace's own process files (AGENTS.md, docs/roles/, docs/decisions/,
workflows/, prompts/, scripts/, adapters/, and any installed .claude/.codex/.agents/) are protected
by a core-file-lock hook (AGENTS.md rule 7). If ANEW_ALLOW_CORE_EDIT is not already set, tell me and
wait — don't try to route around it.

Read AGENTS.md and workflows/bootstrap.md. We are adapting this workspace to a real project.

1. INSPECT the working tree. If it contains application code, detect stack, build/test tooling and
   existing conventions from the code, and present what you found for confirmation. If empty,
   we start from my answers.
2. INTERVIEW me one topic at a time — with your own recommendation and rationale for every
   question (proposal rule):
   product & goal · domain terms and business rules · architecture style & module boundaries ·
   forbidden dependencies · conventions (data rules, error handling, naming) · testing
   expectations · security posture (INCLUDING which dependency/CVE audit and SAST tool fits this
   stack — recommend one, don't leave docs/security.md's "Automated checks" unfilled) · git rules ·
   operating mode (lite or strict — recommend one based on team size and risk).
3. GENERATE from my answers: fill every template in docs/ (architecture, domain, conventions,
   testing, security, git); write scripts/check.conf with real build/test/lint commands for this
   stack PLUS a `security:` step (dependency audit and/or SAST — see scripts/check.conf.example);
   set the Mode line in AGENTS.md; rewrite the AGENTS.md project summary.
   HARD RULES: the "Invariant rules" block in AGENTS.md is kept verbatim; AGENTS.md stays a
   signpost ≤ 40 lines — it POINTS to docs, it never copies them.
4. VERIFY: run ./scripts/doctor and ./scripts/check and show me the output.
5. REPORT: what you generated, what you assumed, and every open question that still needs my
   decision — with your recommendations.
```
