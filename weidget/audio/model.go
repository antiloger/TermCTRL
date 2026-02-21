package audio

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/antiloger/termctlr/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gen2brain/malgo"
)

// AudioWidget manages the primary system input/output audio devices.
// It tracks OS volume levels and provides real-time RMS signal levels.
//
// Struct fixes from original:
//   - ctx: *malgo.Context → *malgo.AllocatedContext (correct return type of InitContext)
//   - Muted bool → InMuted + OutMuted (need separate mute state per device)
type AudioWidget struct {
	InVolume     int  // current mic volume  (0–MaxInVolume)
	MaxInVolume  int  // ceiling for mic volume (e.g. 100)
	InMuted      bool // mic mute state
	OutVolume    int  // current speaker volume (0–MaxOutVolume)
	MaxOutVolume int  // ceiling for speaker volume (e.g. 100)
	OutMuted     bool // speaker mute state
	Hop          int  // step size for Inc/Dec (e.g. 5 = 5%)

	// internal — not exported
	ctx       *malgo.AllocatedContext
	inDevice  *malgo.Device
	outDevice *malgo.Device
	inLevel   atomic.Uint64 // float64 bits of RMS
	outLevel  atomic.Uint64 // float64 bits of RMS
	mu        sync.Mutex
}

type Model struct {
	outVolume int
	audio     *AudioWidget
	pos       types.Position
	err       error
	pactlCh   chan struct{} // held for lifetime of model
}

func NewModel() (Model, error) {
	w, err := New(5, 100, 100)
	if err := w.startInputMonitor(); err != nil {
		return Model{}, err
	}
	if err := w.startOutputMonitor(); err != nil {
		return Model{}, err
	}
	if err != nil {
		return Model{}, err
	}
	return Model{
		audio:   w,
		pactlCh: NewPactlSubscriber(),
	}, nil
}

func (m Model) Init() tea.Cmd {
	return WaitForVolumeChange(m.pactlCh)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "=":
			m.err = m.audio.IncOut() // mutates through pointer — safe
		case "-":
			m.err = m.audio.DecOut()
		case "m":
			m.err = m.audio.ToggleMuteOut()
		case "q":
			m.audio.Close()
			return m, tea.Quit
		}
	case VolumeChangedMsg:
		_ = m.audio.Sync()
		return m, WaitForVolumeChange(m.pactlCh) // re-issue to keep waiting for next event
	}
	return m, nil
}

func (m Model) View() string {
	// rms, db := m.audio.OutLevel() // atomic.Load inside — safe
	// return fmt.Sprintf("Vol: %d%%  Muted: %v  RMS: %.3f  dB: %.1f | scr x:%d y:%d ",
	// 	m.audio.OutVolume, m.audio.OutMuted, rms, db, m.pos.X, m.pos.Y)
	outVol := fmt.Sprintf("%d%%", m.audio.OutVolume)
	inVol := fmt.Sprintf("%d%%", m.audio.InVolume)
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, "Vol: ", m.UIVolumeOut(), "  ", outVol),
		" ",
		lipgloss.JoinHorizontal(lipgloss.Center, "Mic: ", m.UIVolumeIn(), "  ", inVol),
	)
}

func (m Model) SetPosition(x, y int) {
	m.pos.X = x
	m.pos.Y = y
}
