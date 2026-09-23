// Command ianua-tui is a read-only terminal dashboard for an Ianua workspace (classic or
// .ianua/-nested install, docs/decisions/0010-nested-install.md): the same information
// scripts/doctor, scripts/status, and scripts/init already expose, in one interactive view,
// in the spirit of lazyskills.sh. It never writes any file — every check re-implements the
// shell scripts' own logic, read-only, so nothing here can drift into being a second source
// of truth. See tui/README.md.
package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	ws, err := FindWorkspace()
	if err != nil {
		fmt.Fprintln(os.Stderr, "ianua-tui:", err)
		fmt.Fprintln(os.Stderr, "Run this from an Ianua project root (classic install) or the")
		fmt.Fprintln(os.Stderr, "directory containing a .ianua/ folder (nested install) — never")
		fmt.Fprintln(os.Stderr, "from inside .ianua/ itself.")
		os.Exit(1)
	}

	// --check: print the health view as plain text and exit, no TUI. Useful in CI/scripts,
	// and for anyone who wants scripts/doctor's checks without the interactive dashboard.
	if len(os.Args) > 1 && (os.Args[1] == "--check" || os.Args[1] == "-check") {
		fails := 0
		for _, r := range RunHealthChecks(ws) {
			fmt.Println(plainBadge(r.Level), r.Text)
			if r.Level == LevelFail {
				fails++
			}
		}
		if fails > 0 {
			os.Exit(1)
		}
		return
	}

	p := tea.NewProgram(newModel(ws))
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "ianua-tui:", err)
		os.Exit(1)
	}
}

func plainBadge(l CheckLevel) string {
	switch l {
	case LevelOK:
		return "ok  "
	case LevelWarn:
		return "WARN"
	default:
		return "FAIL"
	}
}
