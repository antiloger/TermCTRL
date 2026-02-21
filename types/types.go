package types

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Position struct {
	X int
	Y int
}

type TickMsg time.Time

// One shared tick function
func Tick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
