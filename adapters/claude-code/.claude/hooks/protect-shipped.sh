#!/bin/sh
# PreToolUse hook — two protected tiers (AGENTS.md invariant rule 7):
#  1) specs/done/ — always immutable. No override, ever.
#  2) This workspace's own process files (AGENTS.md, workflows/, prompts/, scripts/, adapters/,
#     docs/roles/, docs/decisions/, spec/plan templates, and the installed .claude/.codex/.agents/
#     adapter copies) — immutable by default, editable only when the human consciously sets
#     ANEW_ALLOW_CORE_EDIT=1 (bootstrap, a dedicated ADR, or deliberate framework maintenance).
#     Never as a silent side effect of ordinary feature/bugfix work.
#
# Runs on Edit/Write (file_path) AND Bash (command text, heuristically — a courtesy layer, not a
# security boundary: a sufficiently indirect shell command can still slip past a substring check).
input=$(cat)
file_path=$(printf '%s' "$input" | grep -o '"file_path"[[:space:]]*:[[:space:]]*"[^"]*"' | head -n 1 | sed 's/.*"file_path"[[:space:]]*:[[:space:]]*"//; s/"$//')
command=$(printf '%s' "$input" | grep -o '"command"[[:space:]]*:[[:space:]]*"[^"]*"' | head -n 1 | sed 's/.*"command"[[:space:]]*:[[:space:]]*"//; s/"$//')

SHIPPED_MSG="specs/done/ is immutable (shipped specs are the historical record — AGENTS.md rule 7). Behavior changing? Open a NEW spec in specs/active/ referencing the old one."
CORE_MSG_PREFIX="is a workspace process file (AGENTS.md rule 7). These change only via bootstrap or a dedicated ADR — set ANEW_ALLOW_CORE_EDIT=1 to do this consciously, don't let it happen as a side effect of feature/bugfix work."

block() { echo "Blocked: $1" >&2; exit 2; }

has_write_marker() {
  case "$1" in
    *'>'*|*'rm '*|*'mv '*|*'cp '*|*'sed -i'*|*'tee '*|*'git mv'*|*'git rm'*) return 0 ;;
    *) return 1 ;;
  esac
}

is_core_path() {
  case "$1" in
    AGENTS.md|CLAUDE.md|workflows/*|prompts/*|scripts/check|scripts/check.conf|scripts/doctor|scripts/init|scripts/status|adapters/*|docs/roles/*|docs/decisions/*|specs/TEMPLATE.md|specs/plans/TEMPLATE.md|specs/epics/TEMPLATE.md|.claude/*|.codex/*|.agents/*)
      return 0 ;;
    *) return 1 ;;
  esac
}

# --- tier 1: specs/done/, absolute, no override ---
case "$file_path" in *specs/done/*) block "$SHIPPED_MSG" ;; esac
if [ -n "$command" ] && has_write_marker "$command"; then
  case "$command" in *specs/done/*) block "$SHIPPED_MSG" ;; esac
fi

# --- tier 2: workspace process files, override with ANEW_ALLOW_CORE_EDIT=1|true|yes ---
case "${ANEW_ALLOW_CORE_EDIT:-}" in 1|true|yes) exit 0 ;; esac

if [ -n "$file_path" ] && is_core_path "$file_path"; then
  block "$file_path $CORE_MSG_PREFIX"
fi

if [ -n "$command" ] && has_write_marker "$command"; then
  for pat in AGENTS.md CLAUDE.md workflows/ prompts/ scripts/check scripts/doctor scripts/init scripts/status adapters/ docs/roles/ docs/decisions/ specs/TEMPLATE.md specs/plans/TEMPLATE.md specs/epics/TEMPLATE.md .claude/ .codex/ .agents/; do
    case "$command" in
      *"$pat"*) block "this command touches $pat, which $CORE_MSG_PREFIX" ;;
    esac
  done
fi

exit 0
