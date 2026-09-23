package main

import "github.com/charmbracelet/lipgloss"

var (
	colorOK     = lipgloss.Color("42")  // green
	colorWarn   = lipgloss.Color("214") // amber
	colorFail   = lipgloss.Color("203") // red
	colorMuted  = lipgloss.Color("246") // gray
	colorAccent = lipgloss.Color("81")  // cyan

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15")).
			Background(colorAccent).
			Padding(0, 1)

	subtitleStyle = lipgloss.NewStyle().Foreground(colorMuted)

	tabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(colorAccent).
			Padding(0, 2)

	tabInactiveStyle = lipgloss.NewStyle().
				Foreground(colorMuted).
				Padding(0, 2)

	okStyle    = lipgloss.NewStyle().Foreground(colorOK).Bold(true)
	warnStyle  = lipgloss.NewStyle().Foreground(colorWarn).Bold(true)
	failStyle  = lipgloss.NewStyle().Foreground(colorFail).Bold(true)
	mutedStyle = lipgloss.NewStyle().Foreground(colorMuted)

	footerStyle = lipgloss.NewStyle().Foreground(colorMuted).Padding(1, 0, 0, 0)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorMuted).
			Padding(0, 1)
)

func levelBadge(l CheckLevel) string {
	switch l {
	case LevelOK:
		return okStyle.Render("OK  ")
	case LevelWarn:
		return warnStyle.Render("WARN")
	default:
		return failStyle.Render("FAIL")
	}
}
