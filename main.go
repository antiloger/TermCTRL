package main

import (
	"fmt"
	"os"

	"github.com/antiloger/termctlr/weidget/clock"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if _, err := tea.NewProgram(clock.NewClockWidget(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	// clock.Testascii()
}
