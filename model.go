package main

import (
	"github.com/antiloger/termctlr/message"
	"github.com/antiloger/termctlr/types"
	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	window      types.Position
	screens     map[string]tea.Model
	currScrreen string
	shared      map[string]interface{}
	quit        bool
}

func NewModel(screens map[string]tea.Model) Model {
	return Model{
		screens: screens,
		quit:    false,
	}
}

func (m *Model) AddScreen(name string, screen tea.Model) {
	m.screens[name] = screen
}

func (m *Model) SetShared(key string, value interface{}) {
	m.shared[key] = value
}

func (m *Model) SetCurrentScreen(name string) {
	m.currScrreen = name
}

func (m Model) Init() tea.Cmd {
	// Initialize ALL screens (important for clocks, spinners etc.)
	var cmds []tea.Cmd
	for _, screen := range m.screens {
		cmds = append(cmds, screen.Init())
	}
	return tea.Batch(cmds...)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case message.SwitchScreenMsg:
		m.currScrreen = string(msg)
		return m, nil
	case message.QuitMsg:
		return m, tea.Quit
	case tea.WindowSizeMsg:
		m.window.X = msg.Width
		m.window.Y = msg.Height
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" || msg.String() == "q" {
			if len(m.screens) == 0 {
				return m, tea.Quit
			}
		}
	}
	currM, ok := m.screens[m.currScrreen]
	if ok {
		updated, cmd := currM.Update(msg)
		m.screens[m.currScrreen] = updated
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	if currM, ok := m.screens[m.currScrreen]; ok {
		return currM.View()
	}
	return "no screen found: " + m.currScrreen
}
