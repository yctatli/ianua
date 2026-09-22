# Epics

An epic groups several specs that only make sense together toward one shared outcome
(`workflows/README.md`, "Epics"). It is a pointer, not a shortcut — each spec inside it still runs
the full spine independently (Spec → Plan → Build → Review → Verify → Ship). Use
`TEMPLATE.md`, numbered sequentially (`0001-....md`), same convention as `specs/active/`.

Most work doesn't need one. Open an epic only when several genuinely separate specs share one
outcome — not as a container for a single large spec, and not as a substitute for writing the
spec itself.
