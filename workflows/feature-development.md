# Workflow: Feature Development

**Goal:** new behavior, from intent to shipped, with evidence at every gate.

```
INTENT → CLARIFY → SPEC → PLAN → [APPROVAL] → BUILD → REVIEW → [TRIAGE] → VERIFY → SHIP
```

| # | Step | Role | Prompt | Gate / evidence |
|---|---|---|---|---|
| 1 | **INTENT & CLARIFY** *(strict; lite: fold into spec)* | Analyst | `prompts/clarify.md` | Intent in business language + every clarifying question answered by a human. **[GATE: human — strict]** |
| 2 | **SPEC** — create `specs/active/NNNN-<name>.md` | Analyst | `prompts/spec.md` | Atomic, testable criteria; no tech in Requirements. **[GATE: human — strict]** Self-critique pass included. |
| 3 | **PLAN** — create `specs/plans/NNNN-plan.md` | Developer | `prompts/plan.md` | Files + steps + risks (with recommendations) + criterion↔test map. **No code.** |
| 4 | **APPROVAL** | Human | — | **[GATE: human]** Plan touches every criterion? Blast radius sane? Risks honest? Approval recorded in the plan file. |
| 5 | **BUILD** | Developer | `prompts/build.md` | Steps match plan; `scripts/check` green. Deviation → R-07. |
| 6 | **INDEPENDENT REVIEW** — fresh session(s) / read-only subagent(s) | Reviewer (+ Security, strict mode) | `prompts/review.md` (+ `prompts/security-review.md` in strict) | Findings with evidence (file:line) across all six dimensions, or "clean". Strict mode: Security runs as its own independent pass, not a bullet inside Reviewer's — `docs/roles/security.md`. |
| 7 | **TRIAGE** | Human | — | **[GATE: human]** Each finding: real (fix) / noise (reject, write why) / investigate (→ QA, R-05). Security: Critical/High cannot close as noise (`docs/security.md`, "Severity & gating") — fix, or a dedicated risk-acceptance ADR. Critical → R-13. Record every finding + its disposition in the plan's **Findings log** — the reviewer is read-only and reported this in chat, but the persistent record is this file, not the transcript. |
| 8 | **FIX ROUNDS** | Developer | `prompts/build.md` §fixes | Only real findings. Re-review the fix diff (step 6, narrow scope). Rounds > 3 → R-06. Update each fixed finding's Status/Resolution in the plan's Findings log as it lands — don't wait until SHIP to reconcile it. |
| 9 | **VERIFY** | QA | `prompts/verify.md` | Criterion ↔ evidence table complete. UI criteria: screenshot = evidence. |
| 10 | **SHIP** | Human | — | **[GATE: human]** DoD checklist in spec all green → PR (template) → merge → move spec to `specs/done/` → fill scorecard, and the Lesson section if the numbers show real friction (spec's own template). |
