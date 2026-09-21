#!/usr/bin/env python3
# PreToolUse policy hook: specs/done/ is immutable (AGENTS.md invariant rule 7), and a handful of
# destructive git/shell ops are blocked outright (docs/git.md "Forbidden"). Ported from the
# Claude Code adapter's protect-shipped.sh hook to Codex's PreToolUse hook contract.
#
# This is a courtesy layer, not a security boundary — for a real boundary, run risky sessions
# (e.g. REVIEW) under Codex's read-only sandbox (sandbox_mode = "read-only").
import json
import re
import sys


def deny(reason):
    print(json.dumps({
        "hookSpecificOutput": {
            "hookEventName": "PreToolUse",
            "permissionDecision": "deny",
            "permissionDecisionReason": reason,
        }
    }))
    sys.exit(0)


try:
    event = json.load(sys.stdin)
except ValueError:
    sys.exit(0)

tool_name = event.get("tool_name", "")
command = event.get("tool_input", {}).get("command", "") or ""

SHIPPED_MSG = (
    "specs/done/ is immutable (AGENTS.md rule 7). Behavior changing? Open a NEW spec in "
    "specs/active/ referencing the old one."
)

if tool_name == "apply_patch":
    # tool_input.command is the patch body; any touched path under specs/done/ is a write attempt.
    if "specs/done/" in command:
        deny(SHIPPED_MSG)

if tool_name == "Bash":
    write_markers = (">", ">>", "rm ", "mv ", "cp ", "sed -i", "tee ", "git mv", "git rm")
    if "specs/done/" in command and any(marker in command for marker in write_markers):
        deny(SHIPPED_MSG)

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
