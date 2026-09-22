#!/usr/bin/env python3
# PreToolUse policy hook — two protected tiers (AGENTS.md invariant rule 7), plus a handful of
# destructive git ops blocked outright (docs/git.md "Forbidden"). Ported from the Claude Code
# adapter's protect-shipped.sh hook to Codex's PreToolUse hook contract.
#
#  1) specs/done/ — always immutable. No override, ever.
#  2) This workspace's own process files (AGENTS.md, workflows/, prompts/, scripts/, adapters/,
#     docs/roles/, docs/decisions/, spec/plan templates, and the installed .claude/.codex/.agents/
#     adapter copies) — immutable by default, editable only when the human consciously sets
#     ANEW_ALLOW_CORE_EDIT=1 (bootstrap, a dedicated ADR, or deliberate framework maintenance).
#     Never as a silent side effect of ordinary feature/bugfix work.
#
# This is a courtesy layer, not a security boundary — for a real boundary, run risky sessions
# (e.g. REVIEW) under Codex's read-only sandbox (sandbox_mode = "read-only").
import json
import os
import re
import sys

CORE_PATTERNS = (
    "AGENTS.md", "CLAUDE.md", "workflows/", "prompts/", "scripts/check", "scripts/check.conf",
    "scripts/doctor", "scripts/init", "scripts/status", "adapters/", "docs/roles/",
    "docs/decisions/", "specs/TEMPLATE.md", "specs/plans/TEMPLATE.md", "specs/epics/TEMPLATE.md",
    ".claude/", ".codex/", ".agents/",
)
WRITE_MARKERS = (">", ">>", "rm ", "mv ", "cp ", "sed -i", "tee ", "git mv", "git rm")

SHIPPED_MSG = (
    "specs/done/ is immutable (AGENTS.md rule 7). Behavior changing? Open a NEW spec in "
    "specs/active/ referencing the old one."
)
CORE_MSG = (
    "this touches a workspace process file (AGENTS.md rule 7). These change only via bootstrap or "
    "a dedicated ADR — set ANEW_ALLOW_CORE_EDIT=1 to do this consciously, don't let it happen as a "
    "side effect of feature/bugfix work."
)


def deny(reason):
    print(json.dumps({
        "hookSpecificOutput": {
            "hookEventName": "PreToolUse",
            "permissionDecision": "deny",
            "permissionDecisionReason": reason,
        }
    }))
    sys.exit(0)


def core_edit_allowed():
    return os.environ.get("ANEW_ALLOW_CORE_EDIT", "").lower() in ("1", "true", "yes")


def touches_core(text):
    return any(p in text for p in CORE_PATTERNS)


try:
    event = json.load(sys.stdin)
except ValueError:
    sys.exit(0)

tool_name = event.get("tool_name", "")
command = event.get("tool_input", {}).get("command", "") or ""

if tool_name == "apply_patch":
    # tool_input.command is the patch body; any touched path under specs/done/ is a write attempt.
    if "specs/done/" in command:
        deny(SHIPPED_MSG)
    if not core_edit_allowed() and touches_core(command):
        deny(CORE_MSG)

if tool_name == "Bash":
    if "specs/done/" in command and any(marker in command for marker in WRITE_MARKERS):
        deny(SHIPPED_MSG)

    if any(marker in command for marker in WRITE_MARKERS):
        if not core_edit_allowed() and touches_core(command):
            deny(CORE_MSG)

    forbidden = [
        (r"git\s+push\s+(?:[^\n]*\s)?(--force(?:-with-lease)?|-f)\b", "force push"),
        (r"git\s+reset\s+(?:[^\n]*\s)?--hard\b", "git reset --hard"),
        (r"git\s+rebase\b", "git rebase"),
        (r"rm\s+(?:-\w*\s+)*-\w*[rR]\w*[fF]\w*|rm\s+(?:-\w*\s+)*-\w*[fF]\w*[rR]\w*", "rm -rf"),
    ]
    for pattern, label in forbidden:
        if re.search(pattern, command):
            deny(f"{label} is forbidden by docs/git.md. Ask the human if this is genuinely needed.")

sys.exit(0)
