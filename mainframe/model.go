package mainframe

import (
	"github.com/antiloger/termctlr/types"
	tea "github.com/charmbracelet/bubbletea"
)

type ModelI interface {
	switchScreen(screenName string) tea.Cmd
	getCurrentScreen() string
}

type MainFrame struct {
	m       ModelI
	weidget []tea.Model
	focus   int
	idle    bool
	layout  Layout
	window  types.Position
}

func NewMainFrame(m ModelI, weidget ...tea.Model) MainFrame {
	return MainFrame{
		m:       m,
		weidget: weidget,
		layout:  LayoutHorizontal,
		idle:    true,
		focus:   0,
	}
}

func (M MainFrame) Init() tea.Cmd {
	return nil
}

func (M MainFrame) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			if M.idle {
				return M, tea.Quit
			}
			M.idle = true
		}
	}
	return M, nil
}

func (M MainFrame) View() string {
	return ""
}
