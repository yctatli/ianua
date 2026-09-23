package main

import (
	"os"
	"path/filepath"
)

// Workspace locates an Ianua installation relative to the current directory, mirroring the
// same CORE_ROOT / EXEC_ROOT split scripts/check, scripts/doctor, and scripts/init already use
// (docs/decisions/0010-nested-install.md): in a nested install, the real content (docs/, specs/,
// workflows/, prompts/, scripts/, adapters/) lives under .ianua/, while the actual project's own
// files sit one level up, at the project root.
type Workspace struct {
	CoreRoot string // where AGENTS.md, docs/, specs/, workflows/, scripts/check.conf live
	ExecRoot string // where the project's own code + adapter files (.claude/, CLAUDE.md) live
	Nested   bool
}

// FindWorkspace walks up from the current directory looking for an Ianua core, checking a
// nested ".ianua" subdirectory at each level before falling back to a classic (root) install.
func FindWorkspace() (*Workspace, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	for {
		nested := filepath.Join(dir, ".ianua")
		if isIanuaCore(nested) {
			return &Workspace{CoreRoot: nested, ExecRoot: dir, Nested: true}, nil
		}
		if isIanuaCore(dir) {
			return &Workspace{CoreRoot: dir, ExecRoot: dir, Nested: false}, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return nil, errNotAWorkspace
		}
		dir = parent
	}
}

func isIanuaCore(dir string) bool {
	_, err1 := os.Stat(filepath.Join(dir, "AGENTS.md"))
	_, err2 := os.Stat(filepath.Join(dir, "workflows", "bootstrap.md"))
	return err1 == nil && err2 == nil
}

var errNotAWorkspace = &workspaceError{}

type workspaceError struct{}

func (e *workspaceError) Error() string {
	return "not an Ianua workspace: no AGENTS.md + workflows/bootstrap.md found here, in .ianua/, or in any parent directory"
}
