package sysmonitor

import (
	"context"
	"fmt"

	"github.com/antiloger/termctlr/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	info *SystemStats
	pos  types.Position
	err  error
}

func NewModel() Model {
	var info SystemStats
	return Model{
		info: &info,
	}
}

func (m Model) Init() tea.Cmd {
	ctx := context.Background()
	m.info.Start(ctx)
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg.(type) {
	// case types.TickMsg:
	// 	m.info.Read()
	// 	return m, nil
	// }
	return m, nil
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left, "CPU:  ", renderBar(m.info.CPUPercent, 100.0, 20), fmt.Sprintf("  %.1f%%", m.info.CPUPercent)),
		lipgloss.JoinHorizontal(lipgloss.Left, "RAM:  ", renderBar(m.info.RAMPercent, 100.0, 20), fmt.Sprintf("  %.1f%%", m.info.RAMPercent)),
		lipgloss.JoinHorizontal(lipgloss.Left, "Disk: ", renderBar(m.info.DiskPercent, 100.0, 20), fmt.Sprintf("  %.1f%%", m.info.DiskPercent)),
	)
}

func (m Model) SetPosition(x, y int) {
	m.pos.X = x
	m.pos.Y = y
}
