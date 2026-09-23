package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type tab int

const (
	tabHealth tab = iota
	tabStatus
	tabAdapters
	tabCount
)

var tabNames = [tabCount]string{"Health", "Status", "Adapters"}

type model struct {
	ws     *Workspace
	active tab
	width  int
	height int

	health       []CheckResult
	specs        []Spec
	epics        []Epic
	adapterFiles []AdapterFile
	meta         AgentsMeta

	message string
}

func newModel(ws *Workspace) model {
	m := model{ws: ws}
	m.refresh()
	return m
}

func (m *model) refresh() {
	m.health = RunHealthChecks(m.ws)
	m.specs = LoadSpecs(m.ws)
	m.epics = LoadEpics(m.ws)
	m.adapterFiles = CheckAdapterFiles(m.ws)
	m.meta = LoadAgentsMeta(m.ws)
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "tab", "right", "l":
			m.active = (m.active + 1) % tabCount
			m.message = ""
		case "shift+tab", "left", "h":
			m.active = (m.active - 1 + tabCount) % tabCount
			m.message = ""
		case "1":
			m.active = tabHealth
		case "2":
			m.active = tabStatus
		case "3":
			m.active = tabAdapters
		case "r":
			m.refresh()
			m.message = "refreshed"
		}
	}
	return m, nil
}

func (m model) View() string {
	var b strings.Builder

	header := titleStyle.Render("ianua-tui") + "  " +
		subtitleStyle.Render(headerSubtitle(m.ws))
	b.WriteString(header + "\n\n")

	for i := tab(0); i < tabCount; i++ {
		if i == m.active {
			b.WriteString(tabActiveStyle.Render(tabNames[i]))
		} else {
			b.WriteString(tabInactiveStyle.Render(tabNames[i]))
		}
	}
	b.WriteString("\n\n")

	switch m.active {
	case tabHealth:
		b.WriteString(m.viewHealth())
	case tabStatus:
		b.WriteString(m.viewStatus())
	case tabAdapters:
		b.WriteString(m.viewAdapters())
	}

	footer := "tab/←→ switch · r refresh · q quit"
	if m.message != "" {
		footer = m.message + "  ·  " + footer
	}
	b.WriteString("\n" + footerStyle.Render(footer))

	return b.String()
}

func headerSubtitle(ws *Workspace) string {
	kind := "classic install"
	if ws.Nested {
		kind = "nested install (.ianua/)"
	}
	return fmt.Sprintf("%s — core: %s", kind, ws.CoreRoot)
}

func (m model) viewHealth() string {
	var b strings.Builder
	fails, warns := 0, 0
	for _, r := range m.health {
		b.WriteString(levelBadge(r.Level) + "  " + r.Text + "\n")
		switch r.Level {
		case LevelFail:
			fails++
		case LevelWarn:
			warns++
		}
	}
	summary := okStyle.Render(fmt.Sprintf("%d ok", len(m.health)-fails-warns))
	if warns > 0 {
		summary += "  " + warnStyle.Render(fmt.Sprintf("%d warn", warns))
	}
	if fails > 0 {
		summary += "  " + failStyle.Render(fmt.Sprintf("%d fail", fails))
	}
	b.WriteString("\n" + summary + "\n")
	return b.String()
}

func (m model) viewStatus() string {
	var b strings.Builder
	if len(m.specs) == 0 {
		b.WriteString(mutedStyle.Render("(no specs in specs/active/ — nothing in flight)") + "\n")
	}
	for _, s := range m.specs {
		name := s.Name
		if name == "" {
			name = s.File
		}
		b.WriteString(fmt.Sprintf("%s — %s\n", s.File, name))
		plan := "none yet"
		if s.PlanOK {
			plan = "written"
		}
		b.WriteString(mutedStyle.Render(fmt.Sprintf(
			"  status: %s · mode: %s · plan: %s\n  DoD: %d/%d checked\n",
			orDash(s.Status), orDash(s.Mode), plan, s.DoDDone, s.DoDTotal,
		)))
		if s.Findings.Total() > 0 {
			line := fmt.Sprintf("  findings: %d open, %d fixed, %d noise, %d deferred",
				s.Findings.Open, s.Findings.Fixed, s.Findings.Noise, s.Findings.Deferred)
			if s.Findings.Open > 0 {
				b.WriteString(warnStyle.Render(line) + "\n")
			} else {
				b.WriteString(mutedStyle.Render(line) + "\n")
			}
		}
		if s.Epic != "" {
			b.WriteString(mutedStyle.Render("  epic: "+s.Epic) + "\n")
		}
		b.WriteString("\n")
	}
	if len(m.epics) > 0 {
		b.WriteString("epics:\n")
		for _, e := range m.epics {
			b.WriteString(fmt.Sprintf("  %s — %s (%s)\n", e.File, orDash(e.Name), orDash(e.Status)))
		}
	}
	return b.String()
}

func (m model) viewAdapters() string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Mode: %s   Chat language: %s\n\n", m.meta.Mode, m.meta.ChatLanguage))

	if len(m.adapterFiles) == 0 {
		b.WriteString(mutedStyle.Render("no adapter discovery files found at the project root — run scripts/init <tool>") + "\n")
		return b.String()
	}
	for _, af := range m.adapterFiles {
		switch af.Status {
		case LinkOKSymlink:
			b.WriteString(okStyle.Render("OK  ") + "  " + af.Path + mutedStyle.Render(" -> "+af.Target) + "\n")
		case LinkOKPlain:
			b.WriteString(okStyle.Render("OK  ") + "  " + af.Path + mutedStyle.Render(" (plain file/dir, classic install)") + "\n")
		case LinkBroken:
			b.WriteString(failStyle.Render("FAIL") + "  " + af.Path + mutedStyle.Render(" -> "+af.Target+" (target missing!)") + "\n")
		}
	}
	return b.String()
}

func orDash(s string) string {
	if s == "" {
		return "?"
	}
	return s
}
