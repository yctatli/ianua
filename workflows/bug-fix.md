# Workflow: Bug Fix

**Goal:** fix the cause, prove it stays fixed. **Iron rule: no fix before a failing reproduction.**

```
REPORT → REPRODUCE (red test) → DIAGNOSE → FIX → PROVE → REVIEW → SHIP
```

1. **REPORT.** Capture: expected vs actual, environment, frequency. If it came from review triage,
   the finding is the report.
2. **REPRODUCE** — Role: QA. Write the *minimal failing test* that captures the bug
   (`prompts/recovery/unverified-finding.md`, R-05, is the pattern). No reproduction after honest
   effort → propose a reasoned close. **[GATE: human]**
3. **DIAGNOSE** — Role: Developer. Root cause, not symptom (R-01 discipline): which file, which
   assumption, why now. Diagnosis before any change.
4. **FIX.** Smallest change that makes the red test green. Touch only files related to the root
   cause. Symptom-silencing (try/catch burial, test edits) is forbidden.
5. **PROVE.** Red test now green **and** whole suite green (`scripts/check`) — no new regressions
   (else R-04: revert first). The reproduction test stays forever as a regression guard.
6. **REVIEW** — fresh session, narrow scope: the fix diff. Strict mode, if the fix touches
   auth/input/data handling: also a separate Security pass (`docs/roles/security.md`).
   **[GATE: human — strict]** triage. If this fix escalated to a mini-spec with a plan (step 7),
   record findings in that plan's **Findings log**, same convention as feature-development — a
   one-off fix with no plan file can skip this, the fix diff itself is the record.
7. **SHIP.** Small bugs: PR referencing the report. Behavior-changing fixes: they're features —
   write a mini-spec first. If this fix has a spec (mini or full), fill its Scorecard and — if the
   numbers show real friction — its Lesson section too.
