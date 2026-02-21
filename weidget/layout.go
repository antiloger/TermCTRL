package weidget

import (
	"github.com/charmbracelet/lipgloss"
)

type Layout int

const (
	Vertical Layout = iota
)

func (W *WeidgetScreen) applyLayout() string {
	switch W.layout {
	case Vertical:
		var views []string
		for i, widget := range W.weidgets {
			v := widget.View()
			// Debug: print each widget view info
			if i == W.focus {
				v = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true).Render(v)
			}
			views = append(views, v)
		}

		joined := lipgloss.JoinVertical(lipgloss.Center, views...)

		return joined
	}
	return "ERROR"
}
