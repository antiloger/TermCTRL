package weidget

import (
	"fmt"

	"github.com/antiloger/termctlr/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Weidget interface {
	tea.Model
	SetPosition(x, y int)
}

type WeidgetScreen struct {
	weidgets   []Weidget
	focus      int
	screenSize types.Position
	idle       bool
	layout     Layout
	Tick       int
}

func NewWeidgetScreen(layout Layout, weidgets ...Weidget) WeidgetScreen {
	return WeidgetScreen{
		weidgets: weidgets,
		focus:    0,
		idle:     true,
		layout:   layout,
		Tick:     0,
	}
}

func (W WeidgetScreen) Init() tea.Cmd {
	var cmds []tea.Cmd
	cmds = append(cmds, types.Tick()) // ONE tick for all widgets
	for _, widget := range W.weidgets {
		cmds = append(cmds, widget.Init())
	}
	return tea.Batch(cmds...)
}

func (W WeidgetScreen) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Focus navigation
		switch msg.String() {
		case "tab":
			W.focus = (W.focus + 1) % len(W.weidgets)
			return W, nil
		case "shift+tab":
			W.focus--
			if W.focus < 0 {
				W.focus = len(W.weidgets) - 1
			}
			return W, nil
		}

		// ← KeyMsg only goes to focused widget
		if len(W.weidgets) > 0 {
			updated, cmd := W.weidgets[W.focus].Update(msg)
			W.weidgets[W.focus] = updated.(Weidget)
			return W, cmd
		}
	case tea.WindowSizeMsg:
		W.screenSize.X = msg.Width
		W.screenSize.Y = msg.Height

	case types.TickMsg:
		// Broadcast to ALL widgets but screen owns the next tick
		for i, widget := range W.weidgets {
			updated, _ := widget.Update(msg) // ignore widget's cmd
			W.weidgets[i] = updated.(Weidget)
		}
		W.Tick++
		return W, types.Tick() // ONE tick continues ✅

	default:
		// ← EVERYTHING else (tickMsg, WindowSizeMsg, etc.) → ALL widgets
		var cmds []tea.Cmd
		for i, widget := range W.weidgets {
			updated, cmd := widget.Update(msg)
			W.weidgets[i] = updated.(Weidget)
			cmds = append(cmds, cmd)
		}
		return W, tea.Batch(cmds...)
	}

	return W, nil
}

func (W WeidgetScreen) View() string {
	layout := W.applyLayout()

	layoutH := lipgloss.Height(layout)
	layoutW := lipgloss.Width(layout)

	// Check if content fits
	status := fmt.Sprintf(
		"screen=%dx%d | layout=%dx%d | fits_h=%v | fits_w=%v | widgets=%d | focus=%d | tick=%d",
		W.screenSize.X, W.screenSize.Y,
		layoutW, layoutH,
		layoutH <= W.screenSize.Y,
		layoutW <= W.screenSize.X,
		len(W.weidgets),
		W.focus,
		W.Tick,
	)

	centered := lipgloss.Place(
		W.screenSize.X,
		W.screenSize.Y-1,
		lipgloss.Center,
		lipgloss.Center,
		layout,
	)

	return lipgloss.JoinVertical(lipgloss.Left, centered, status)
}
