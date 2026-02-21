package main

import (
	"fmt"
	"log"
	"os"

	"github.com/antiloger/termctlr/weidget"
	"github.com/antiloger/termctlr/weidget/audio"
	"github.com/antiloger/termctlr/weidget/clock"
	sysinfo "github.com/antiloger/termctlr/weidget/sysInfo"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	clockWidget := clock.NewClockWidget()
	specWidget := sysinfo.NewSysInfoWidget()
	audioWidget, err := audio.NewModel()
	if err != nil {
		log.Fatal("Failed to initialize audio widget:", err)
	}

	weidgetScr := weidget.NewWeidgetScreen(weidget.Vertical, &clockWidget, &specWidget, &audioWidget)

	screens := map[string]tea.Model{
		"weidget": weidgetScr,
	}

	m := NewModel(screens)
	m.SetCurrentScreen("weidget")

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
	// clock.Testascii()
}
