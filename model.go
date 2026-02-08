package main

import (
	"fmt"

	"github.com/antiloger/termctlr/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	window      types.Position
	screens     map[string]tea.Model
	currScrreen string
	shared      map[string]interface{}
	quit        bool
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.setWindowSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	// Define the box style with padding
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Align(lipgloss.Center)

	// Your content inside the box
	content := lipgloss.JoinVertical(
		lipgloss.Center,
		"Welcome to My App",
		"",
		"This box is centered",
		"in the terminal",
		"",
		fmt.Sprintf("Terminal Size: %d x %d", m.window.X, m.window.Y),
		"",
		"Press 'q' to quit",
	)

	// Render the box with content
	box := boxStyle.Render(content)

	// Center the entire box in the terminal
	return lipgloss.Place(
		m.window.X,
		m.window.Y,
		lipgloss.Center, // horizontal center
		lipgloss.Center, // vertical center
		box,
	)
}

func (m *Model) setWindowSize(width, height int) {
	m.window.X = width
	m.window.Y = height
}
