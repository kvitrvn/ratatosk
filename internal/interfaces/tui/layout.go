package tui

import "github.com/charmbracelet/lipgloss"

// frameSize is the number of columns/rows consumed by a NormalBorder on each axis.
const frameSize = 2

var (
	colorSubtle    = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
	colorHighlight = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#7D56F4"}
	colorUnread    = lipgloss.AdaptiveColor{Light: "#EE6FF8", Dark: "#AD58B4"}
)

var (
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorSubtle)

	focusedPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(colorHighlight)

	selectedItemStyle = lipgloss.NewStyle().
				Background(colorSubtle)

	overlayStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(colorHighlight).
			Padding(1, 2)
)

var (
	unreadDot = lipgloss.NewStyle().Foreground(colorUnread).Render("●")
	readDot   = lipgloss.NewStyle().Foreground(colorSubtle).Render("○")
)

// dimensions returns the inner content widths for the feed-list and article/detail panels.
func dimensions(w int) (listW, articleW int) {
	feedOuterW := w * 30 / 100
	listW = feedOuterW - frameSize
	if listW < 0 {
		listW = 0
	}
	articleW = w - feedOuterW - frameSize
	if articleW < 0 {
		articleW = 0
	}
	return
}
