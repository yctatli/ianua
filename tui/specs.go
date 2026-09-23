package main

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Spec struct {
	File     string
	Name     string
	Status   string
	Mode     string
	Epic     string
	DoDDone  int
	DoDTotal int
	PlanFile string
	PlanOK   bool
	Findings FindingsCount
}

// FindingsCount mirrors scripts/status's findings_summary — a rollup of a plan's
// "## Findings log" table (specs/plans/TEMPLATE.md), never hand-maintained.
type FindingsCount struct {
	Open, Fixed, Noise, Deferred, Other int
}

func (f FindingsCount) Total() int {
	return f.Open + f.Fixed + f.Noise + f.Deferred + f.Other
}

type Epic struct {
	File   string
	Name   string
	Status string
}

var fieldRe = regexp.MustCompile(`^- ([A-Za-z ]+): *(.*)$`)
var dodRe = regexp.MustCompile(`^- \[([ xX])\]`)

// LoadSpecs mirrors scripts/status: reads every spec in specs/active/, deriving everything from
// the files themselves — nothing here is hand-maintained.
func LoadSpecs(ws *Workspace) []Spec {
	dir := filepath.Join(ws.CoreRoot, "specs", "active")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var specs []Spec
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || name == "README.md" || name == "TEMPLATE.md" {
			continue
		}
		s := parseSpec(filepath.Join(dir, name))
		s.File = name
		if len(name) >= 4 {
			num := name[:4]
			planPath := filepath.Join(ws.CoreRoot, "specs", "plans", num+"-plan.md")
			if _, err := os.Stat(planPath); err == nil {
				s.PlanFile = num + "-plan.md"
				s.PlanOK = true
				s.Findings = parseFindingsLog(planPath)
			}
		}
		specs = append(specs, s)
	}
	return specs
}

func LoadEpics(ws *Workspace) []Epic {
	dir := filepath.Join(ws.CoreRoot, "specs", "epics")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var epics []Epic
	for _, e := range entries {
		name := e.Name()
		if e.IsDir() || !strings.HasSuffix(name, ".md") || name == "README.md" || name == "TEMPLATE.md" {
			continue
		}
		f, err := os.Open(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		ep := Epic{File: name}
		sc := bufio.NewScanner(f)
		first := true
		for sc.Scan() {
			line := sc.Text()
			if first {
				ep.Name = strings.TrimPrefix(strings.TrimSpace(line), "# ")
				first = false
			}
			if m := fieldRe.FindStringSubmatch(line); m != nil && strings.TrimSpace(m[1]) == "Status" {
				ep.Status = m[2]
			}
		}
		f.Close()
		epics = append(epics, ep)
	}
	return epics
}

// parseFindingsLog reads a plan file's "## Findings log" table and tallies the Status column —
// the same logic as scripts/status's findings_summary awk function, kept in sync by hand
// (docs/decisions/0012-tui-dashboard.md flags this as the real cost of a second implementation).
func parseFindingsLog(planPath string) FindingsCount {
	var fc FindingsCount
	f, err := os.Open(planPath)
	if err != nil {
		return fc
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	inSection := false
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "## Findings log") {
			inSection = true
			continue
		}
		if strings.HasPrefix(line, "## ") {
			if inSection {
				break
			}
			continue
		}
		if !inSection || !strings.HasPrefix(line, "|") {
			continue
		}
		cols := strings.Split(line, "|")
		if len(cols) < 8 {
			continue
		}
		id := strings.TrimSpace(cols[1])
		if id == "" || id == "ID" || isDashes(id) {
			continue
		}
		status := strings.ToLower(strings.TrimSpace(cols[6]))
		switch {
		case strings.HasPrefix(status, "open"):
			fc.Open++
		case strings.HasPrefix(status, "fixed"):
			fc.Fixed++
		case strings.HasPrefix(status, "noise"):
			fc.Noise++
		case strings.HasPrefix(status, "deferred"):
			fc.Deferred++
		default:
			fc.Other++
		}
	}
	return fc
}

func isDashes(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if c != '-' {
			return false
		}
	}
	return true
}

func parseSpec(path string) Spec {
	s := Spec{}
	f, err := os.Open(path)
	if err != nil {
		return s
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	first := true
	for sc.Scan() {
		line := sc.Text()
		if first {
			s.Name = strings.TrimPrefix(strings.TrimSpace(line), "# ")
			first = false
			continue
		}
		if m := fieldRe.FindStringSubmatch(line); m != nil {
			switch strings.TrimSpace(m[1]) {
			case "Status":
				if s.Status == "" {
					s.Status = m[2]
				}
			case "Mode":
				if s.Mode == "" {
					s.Mode = m[2]
				}
			case "Epic":
				if s.Epic == "" {
					s.Epic = m[2]
				}
			}
		}
		if m := dodRe.FindStringSubmatch(line); m != nil {
			s.DoDTotal++
			if m[1] == "x" || m[1] == "X" {
				s.DoDDone++
			}
		}
	}
	return s
}
