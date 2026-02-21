package message

import tea "github.com/charmbracelet/bubbletea"

// Screen switching
type SwitchScreenMsg string

func SwitchScreen(index string) tea.Cmd {
	return func() tea.Msg {
		return SwitchScreenMsg(index)
	}
}

// Other custom messages
type QuitMsg struct{}

func Quit() tea.Cmd {
	return func() tea.Msg {
		return QuitMsg{}
	}
}
