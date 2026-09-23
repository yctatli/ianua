package main

import (
	"os"
	"path/filepath"
)

type LinkStatus int

const (
	LinkMissing   LinkStatus = iota
	LinkOKPlain              // present, not a symlink (classic install, or hand-placed)
	LinkOKSymlink            // present, symlink, target resolves
	LinkBroken               // symlink present but target doesn't resolve
)

type AdapterFile struct {
	Path   string // relative to ExecRoot
	Status LinkStatus
	Target string // symlink target, if any
}

// adapterTargets are the discovery files each adapter can place at the project root.
var adapterTargets = []string{
	"AGENTS.md", "CLAUDE.md", ".claude",
	".agents", ".codex",
	".cursor/rules/ianua.mdc",
	".github/copilot-instructions.md",
}

func CheckAdapterFiles(ws *Workspace) []AdapterFile {
	var out []AdapterFile
	for _, rel := range adapterTargets {
		full := filepath.Join(ws.ExecRoot, rel)
		info, err := os.Lstat(full)
		if err != nil {
			continue // not installed — don't clutter the view with every possible target
		}
		af := AdapterFile{Path: rel}
		if info.Mode()&os.ModeSymlink != 0 {
			target, _ := os.Readlink(full)
			af.Target = target
			if _, statErr := os.Stat(full); statErr == nil { // Stat follows the link
				af.Status = LinkOKSymlink
			} else {
				af.Status = LinkBroken
			}
		} else {
			af.Status = LinkOKPlain
		}
		out = append(out, af)
	}
	return out
}

// AgentsMeta pulls the two small self-reported settings lines out of AGENTS.md.
type AgentsMeta struct {
	Mode         string
	ChatLanguage string
}

func LoadAgentsMeta(ws *Workspace) AgentsMeta {
	data, err := os.ReadFile(filepath.Join(ws.CoreRoot, "AGENTS.md"))
	if err != nil {
		return AgentsMeta{}
	}
	text := string(data)
	return AgentsMeta{
		Mode:         extractBold(text, "Mode:"),
		ChatLanguage: extractBold(text, "Chat language:"),
	}
}

// extractBold finds "**<label> <value>**" and returns <value>, trimmed.
func extractBold(text, label string) string {
	idx := indexOf(text, "**"+label)
	if idx == -1 {
		return "unset"
	}
	rest := text[idx+len("**"+label):]
	end := indexOf(rest, "**")
	if end == -1 {
		return "unset"
	}
	val := trim(rest[:end])
	if val == "" {
		return "unset"
	}
	return val
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func trim(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
