# R-13 — Critical Security Finding (owner: Security)

```
A Critical-severity security finding was raised: <finding> (docs/security.md, "Severity & gating").
STOP other work on this change until this is resolved — don't keep building on top of code with a
live Critical finding in it.
1) Confirm the severity is really Critical against docs/security.md's criteria, not a gut call. If
   genuinely unclear, R-05 first (QA reproduction) — don't downgrade the severity to dodge this rule.
2) Surface it to the human immediately, on its own — not buried at the end of a longer report.
3) The only two valid closes: FIXED (normal fix-round path, re-reviewed narrowly), or a dedicated
   ADR that names the accepted risk, its owner, and a revisit trigger (docs/decisions/README.md).
   "Noise" and "we'll get to it later" without an ADR are not valid closes for Critical.
PROTECTION: no ship, and no further build on top of the affected code, while a Critical finding is
open. The fix-or-ADR requirement exists so "we decided to accept this risk" is a decision on
record, not a thing that quietly happened under deadline pressure.
```
