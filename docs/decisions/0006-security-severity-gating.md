# ADR 0006 — Security finding severity taxonomy and ship-gating

- Status: Accepted
- Date: 2026-09-22

## Context
The human asked directly: when Security actually finds something, what's *supposed* to happen, and
what does this workspace's tooling *actually* do right now? Tracing the real path (ADR 0004's
Security role → `workflows/feature-development.md` TRIAGE step) exposed two gaps:

1. **No severity taxonomy.** `prompts/security-review.md` told the model to "order findings by
   severity" without ever defining what severity means. Every finder invented its own scale.
2. **No gating tied to severity.** Every finding — Reviewer's "this variable name is unclear" and
   Security's "this endpoint has no authorization check" — went through the exact same generic
   triage (real / noise / investigate). Nothing stopped a human from triaging a critical, actually-
   exploitable vulnerability as "noise" with a one-line rationale and shipping anyway. For a
   correctness nitpick that discretion is the whole point of human triage; for a Critical security
   finding it's a silent, undocumented risk acceptance — exactly what `docs/security.md`'s
   Automated-checks section (ADR 0003) already argued against for enforcement in general.

A second finding, mid-analysis: nothing distinguished a vulnerability found *pre-ship*, during a
normal REVIEW, from one found in code that's *already shipped* (e.g. an ad hoc `/security` pass on
a live branch). The former is this workspace's normal spec flow; the latter is
`workflows/incident.md` territory and was going completely unrouted.

## Decision
Add a four-level severity taxonomy to `docs/security.md` ("Severity & gating"): Critical / High /
Medium / Low, each with concrete examples and an explicit ship rule:
- **Critical** — cannot be triaged as "noise." Closes only by being fixed, or by a **dedicated ADR**
  naming the accepted risk, its owner, and a revisit trigger — reusing the existing ADR mechanism
  rather than inventing a new "risk acceptance" artifact type. Triggers recovery ramp **R-13**
  (new: `prompts/recovery/critical-security-finding.md`) — stop other work on the change until
  it's resolved.
- **High** — blocks ship in strict mode; lite mode requires explicit written human sign-off, not a
  silent noise triage.
- **Medium** — may ship, but only with a tracked follow-up spec, not a triage note that evaporates.
- **Low** — normal triage, human's call, freely deferrable.

`docs/roles/security.md`, `prompts/security-review.md`,
`workflows/feature-development.md`'s TRIAGE row, and `specs/TEMPLATE.md`'s Definition of Done all
updated to point at this table instead of leaving severity undefined. A vulnerability found in
already-shipped code is explicitly routed to `workflows/incident.md`, with the same Critical/High
ship rules applying there.

## Consequences
**Buys us:** a security finding's fate now depends on what it actually is, not on how much
attention one triage pass happened to give it. "We decided to ship with this open" becomes a
recorded decision (an ADR) instead of a thing that can happen silently in a triage note. Reuses
existing machinery (ADRs, recovery ramps) rather than adding a new process artifact.

**Costs us:** triage now requires an extra judgment call (assign severity, correctly) before the
existing real/noise/investigate call — more friction on every security finding, deliberately, since
the friction is precisely what stops a Critical finding from being waved through by habit. Depends
on the finder (Security, or Reviewer in lite-mode collapse) actually using the table honestly;
nothing stops a bad-faith or careless severity downgrade to dodge the gate — same trust-boundary
limitation every prose rule in this workspace has, mitigated only by review of the finding itself.

## Alternatives considered
- **CVSS-style numeric scoring**: rejected for v1 — more rigorous but heavier than this workspace's
  Markdown-and-judgment style elsewhere; the four-level qualitative scale with concrete examples is
  more consistent with how the rest of ANEW specifies rules. Revisit if a project bootstrapped from
  this template is in a regulated context that expects CVSS.
- **A new "risk acceptance" document type, separate from ADRs**: rejected — an accepted security
  risk is structurally identical to any other hard-to-reverse decision with consequences and a
  revisit trigger; ADRs already have exactly those fields. Inventing a parallel artifact would
  duplicate a mechanism that already fits.
- **Block ship on ANY open finding regardless of severity, no gating table**: rejected — that's
  strict mode's existing ceremony pushed onto every severity indiscriminately, which is what
  produces the "just triage everything as noise to move on" pressure this ADR is trying to relieve
  for Low/Medium findings specifically.

## Revisit triggers
- A Critical/High finding gets shipped via a rushed, low-substance risk-acceptance ADR that's
  really just "noise" wearing an ADR's clothes — tighten what a valid risk-acceptance ADR must
  contain (e.g. require a second human's sign-off, not just the one who triaged it).
- The four severity levels prove too coarse or too fine in practice for a real project bootstrapped
  from this template — adjust the table, but keep it a fixed, written table, never back-slide to
  "order by severity" with no definition.
- A bootstrapped project needs CVSS or another formal scoring system for compliance reasons —
  swap the table, keep the same ship-gating structure (Critical/High block, Medium needs follow-up,
  Low is free).
