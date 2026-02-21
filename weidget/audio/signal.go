package audio

import (
	"bufio"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type VolumeChangedMsg struct{}

// NewPactlSubscriber starts a single long-lived goroutine that watches
// pactl subscribe and forwards sink/source events to the returned channel.
// Call this once (e.g. in NewModel) and hold onto the channel.
func NewPactlSubscriber() chan struct{} {
	ch := make(chan struct{})
	go func() {
		cmd := exec.Command("pactl", "subscribe")
		out, err := cmd.StdoutPipe()
		if err != nil {
			return
		}
		if err := cmd.Start(); err != nil {
			return
		}
		defer cmd.Process.Kill()

		scanner := bufio.NewScanner(out)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "sink") || strings.Contains(line, "source") {
				ch <- struct{}{}
			}
		}
	}()
	return ch
}

// WaitForVolumeChange returns a tea.Cmd that blocks until the next event
// arrives on the channel, then emits VolumeChangedMsg.
// This is safe to re-issue after every message — the goroutine is NOT restarted.
func WaitForVolumeChange(ch chan struct{}) tea.Cmd {
	return func() tea.Msg {
		<-ch // blocks until the goroutine sends
		return VolumeChangedMsg{}
	}
}
