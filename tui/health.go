package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CheckLevel int

const (
	LevelOK CheckLevel = iota
	LevelWarn
	LevelFail
)

type CheckResult struct {
	Level CheckLevel
	Text  string
}

// requiredCoreFiles mirrors the list scripts/doctor checks for structural completeness.
var requiredCoreFiles = []string{
	"AGENTS.md", "README.md",
	"docs/architecture.md", "docs/domain.md", "docs/conventions.md", "docs/testing.md",
	"docs/security.md", "docs/git.md",
	"docs/decisions/TEMPLATE.md", "docs/roles/README.md", "docs/roles/analyst.md",
	"docs/roles/developer.md", "docs/roles/reviewer.md", "docs/roles/security.md", "docs/roles/qa.md",
	"specs/TEMPLATE.md", "specs/plans/TEMPLATE.md", "specs/epics/README.md", "specs/epics/TEMPLATE.md",
	"workflows/README.md", "workflows/bootstrap.md", "workflows/feature-development.md",
	"workflows/bug-fix.md", "workflows/refactor.md", "workflows/incident.md",
	"prompts/README.md", "prompts/bootstrap.md", "prompts/clarify.md", "prompts/spec.md",
	"prompts/plan.md", "prompts/build.md", "prompts/review.md", "prompts/security-review.md",
	"prompts/verify.md", "prompts/adr.md", "prompts/recovery/README.md",
	"scripts/check", "scripts/status",
}

// RunHealthChecks re-implements scripts/doctor's checks natively, so the dashboard doesn't need
// to shell out and re-parse text output.
func RunHealthChecks(ws *Workspace) []CheckResult {
	var results []CheckResult

	missing := 0
	for _, f := range requiredCoreFiles {
		if _, err := os.Stat(filepath.Join(ws.CoreRoot, f)); err != nil {
			results = append(results, CheckResult{LevelFail, "missing: " + f})
			missing++
		}
	}
	if missing == 0 {
		results = append(results, CheckResult{LevelOK, "core structure complete"})
	}

	ramps, _ := filepath.Glob(filepath.Join(ws.CoreRoot, "prompts", "recovery", "*.md"))
	rampCount := 0
	for _, r := range ramps {
		if filepath.Base(r) != "README.md" {
			rampCount++
		}
	}
	if rampCount >= 13 {
		results = append(results, CheckResult{LevelOK, "recovery catalog present (" + strconv.Itoa(rampCount) + " ramps)"})
	} else {
		results = append(results, CheckResult{LevelWarn, "recovery catalog incomplete (found " + strconv.Itoa(rampCount) + ", expected 13)"})
	}

	agents, err := os.ReadFile(filepath.Join(ws.CoreRoot, "AGENTS.md"))
	agentsText := string(agents)
	if err != nil {
		results = append(results, CheckResult{LevelFail, "AGENTS.md unreadable: " + err.Error()})
	} else {
		if strings.Contains(agentsText, "STATUS: NOT CONFIGURED") {
			results = append(results, CheckResult{LevelWarn, "AGENTS.md not configured yet — run the bootstrap workflow"})
		} else {
			results = append(results, CheckResult{LevelOK, "AGENTS.md configured"})
		}

		if strings.Contains(agentsText, "Mode: unset") {
			results = append(results, CheckResult{LevelWarn, "operating mode not set (lite/strict)"})
		} else {
			results = append(results, CheckResult{LevelOK, "operating mode set"})
		}

		if strings.Contains(agentsText, "Chat language: unset") {
			results = append(results, CheckResult{LevelWarn, "chat language not set"})
		} else {
			results = append(results, CheckResult{LevelOK, "chat language set"})
		}

		if strings.Contains(agentsText, "Invariant rules") {
			results = append(results, CheckResult{LevelOK, "AGENTS.md invariant rules block present"})
		} else {
			results = append(results, CheckResult{LevelFail, "AGENTS.md invariant rules block MISSING"})
		}

		lines := strings.Count(agentsText, "\n")
		if lines > 60 {
			results = append(results, CheckResult{LevelWarn, "AGENTS.md is " + strconv.Itoa(lines) + " lines — should be a signpost, not a handbook"})
		}
	}

	if _, err := os.Stat(filepath.Join(ws.CoreRoot, "scripts", "check.conf")); err == nil {
		results = append(results, CheckResult{LevelOK, "verification contract configured (scripts/check.conf)"})
	} else {
		results = append(results, CheckResult{LevelWarn, "scripts/check.conf missing — bootstrap generates it"})
	}

	adapterMarkers := []string{
		"CLAUDE.md", ".github/copilot-instructions.md", ".cursor/rules", ".agents/skills", ".codex/hooks.json",
	}
	adapterFound := false
	for _, m := range adapterMarkers {
		if _, err := os.Stat(filepath.Join(ws.ExecRoot, m)); err == nil {
			adapterFound = true
			break
		}
	}
	if adapterFound {
		results = append(results, CheckResult{LevelOK, "an AI tool adapter is installed"})
	} else {
		results = append(results, CheckResult{LevelWarn, "no adapter installed — run scripts/init <tool>"})
	}

	hookDir := filepath.Join(ws.ExecRoot, ".claude", "hooks")
	if entries, err := os.ReadDir(hookDir); err == nil {
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".sh") {
				info, err := e.Info()
				if err == nil && info.Mode()&0o111 == 0 {
					results = append(results, CheckResult{LevelWarn, "hook not executable: " + e.Name() + " (chmod +x it)"})
				}
			}
		}
	}

	return results
}
