package audio

import (
	"encoding/binary"
	"fmt"
	"math"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/gen2brain/malgo"
)

// ── Constructor ───────────────────────────────────────────────────────────────

// New initialises the AudioWidget and starts the input/output monitor streams.
// hop is the volume step size (e.g. 5 for 5%), max values cap Inc operations.
func New(hop, maxInVolume, maxOutVolume int) (*AudioWidget, error) {
	if hop <= 0 {
		hop = 5
	}
	if maxInVolume <= 0 {
		maxInVolume = 100
	}
	if maxOutVolume <= 0 {
		maxOutVolume = 100
	}

	w := &AudioWidget{
		Hop:          hop,
		MaxInVolume:  maxInVolume,
		MaxOutVolume: maxOutVolume,
	}

	// init malgo context (no system libs needed — C is bundled)
	ctx, err := malgo.InitContext(nil, malgo.ContextConfig{}, func(string) {})
	if err != nil {
		return nil, fmt.Errorf("malgo context: %w", err)
	}
	w.ctx = ctx

	if err := w.startInputMonitor(); err != nil {
		ctx.Uninit()
		return nil, fmt.Errorf("input monitor: %w", err)
	}
	if err := w.startOutputMonitor(); err != nil {
		w.inDevice.Uninit()
		ctx.Uninit()
		return nil, fmt.Errorf("output monitor: %w", err)
	}

	// pull current OS volumes into struct fields
	_ = w.Sync()
	return w, nil
}

// Close stops all streams and frees resources.
func (w *AudioWidget) Close() {
	if w.outDevice != nil {
		w.outDevice.Stop()
		w.outDevice.Uninit()
	}
	if w.inDevice != nil {
		w.inDevice.Stop()
		w.inDevice.Uninit()
	}
	if w.ctx != nil {
		w.ctx.Uninit()
	}
}

// ── OS Volume control (pactl) ─────────────────────────────────────────────────

// IncOut raises speaker volume by Hop, capped at MaxOutVolume.
func (w *AudioWidget) IncOut() error {
	return w.setOutVol(clamp(w.OutVolume+w.Hop, 0, w.MaxOutVolume))
}

// DecOut lowers speaker volume by Hop.
func (w *AudioWidget) DecOut() error {
	return w.setOutVol(clamp(w.OutVolume-w.Hop, 0, w.MaxOutVolume))
}

// IncIn raises mic volume by Hop, capped at MaxInVolume.
func (w *AudioWidget) IncIn() error {
	return w.setInVol(clamp(w.InVolume+w.Hop, 0, w.MaxInVolume))
}

// DecIn lowers mic volume by Hop.
func (w *AudioWidget) DecIn() error {
	return w.setInVol(clamp(w.InVolume-w.Hop, 0, w.MaxInVolume))
}

// MuteOut mutes the speaker.
func (w *AudioWidget) MuteOut() error { return w.setOutMute(true) }

// UnmuteOut unmutes the speaker.
func (w *AudioWidget) UnmuteOut() error { return w.setOutMute(false) }

// ToggleMuteOut flips the speaker mute state.
func (w *AudioWidget) ToggleMuteOut() error { return w.setOutMute(!w.OutMuted) }

// MuteIn mutes the microphone.
func (w *AudioWidget) MuteIn() error { return w.setInMute(true) }

// UnmuteIn unmutes the microphone.
func (w *AudioWidget) UnmuteIn() error { return w.setInMute(false) }

// ToggleMuteIn flips the mic mute state.
func (w *AudioWidget) ToggleMuteIn() error { return w.setInMute(!w.InMuted) }

// Sync reads current OS volume and mute state into the struct fields.
func (w *AudioWidget) Sync() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	vol, muted, err := osGetVolume("sink", "@DEFAULT_SINK@")
	if err != nil {
		return err
	}
	w.OutVolume = vol
	w.OutMuted = muted

	vol, muted, err = osGetVolume("source", "@DEFAULT_SOURCE@")
	if err != nil {
		return err
	}
	w.InVolume = vol
	w.InMuted = muted
	return nil
}

// ── Real-time signal levels (from malgo) ──────────────────────────────────────

// InLevel returns the current microphone RMS and dB values.
func (w *AudioWidget) InLevel() (rms, db float64) {
	rms = math.Float64frombits(w.inLevel.Load())
	db = rmsToDb(rms)
	return
}

// OutLevel returns the current speaker output RMS and dB values.
func (w *AudioWidget) OutLevel() (rms, db float64) {
	rms = math.Float64frombits(w.outLevel.Load())
	db = rmsToDb(rms)
	return
}

// ── Internal helpers ──────────────────────────────────────────────────────────

func (w *AudioWidget) setOutVol(v int) error {
	if err := osSetVolume("sink", "@DEFAULT_SINK@", v); err != nil {
		return err
	}
	w.mu.Lock()
	w.OutVolume = v
	w.mu.Unlock()
	return nil
}

func (w *AudioWidget) setInVol(v int) error {
	if err := osSetVolume("source", "@DEFAULT_SOURCE@", v); err != nil {
		return err
	}
	w.mu.Lock()
	w.InVolume = v
	w.mu.Unlock()
	return nil
}

func (w *AudioWidget) setOutMute(mute bool) error {
	if err := osSetMute("sink", "@DEFAULT_SINK@", mute); err != nil {
		return err
	}
	w.mu.Lock()
	w.OutMuted = mute
	w.mu.Unlock()
	return nil
}

func (w *AudioWidget) setInMute(mute bool) error {
	if err := osSetMute("source", "@DEFAULT_SOURCE@", mute); err != nil {
		return err
	}
	w.mu.Lock()
	w.InMuted = mute
	w.mu.Unlock()
	return nil
}

func (w *AudioWidget) startInputMonitor() error {
	cfg := malgo.DefaultDeviceConfig(malgo.Capture)
	cfg.Capture.Format = malgo.FormatS16
	cfg.Capture.Channels = 1
	cfg.SampleRate = 44100
	cfg.Alsa.NoMMap = 1

	dev, err := malgo.InitDevice(w.ctx.Context, cfg, malgo.DeviceCallbacks{
		Data: func(_, input []byte, _ uint32) {
			w.inLevel.Store(math.Float64bits(calcRMS(input)))
		},
	})
	if err != nil {
		return err
	}
	w.inDevice = dev
	return dev.Start()
}

func (w *AudioWidget) startOutputMonitor() error {
	cfg := malgo.DefaultDeviceConfig(malgo.Playback)
	cfg.Playback.Format = malgo.FormatS16
	cfg.Playback.Channels = 2
	cfg.SampleRate = 44100
	cfg.Alsa.NoMMap = 1

	dev, err := malgo.InitDevice(w.ctx.Context, cfg, malgo.DeviceCallbacks{
		Data: func(output, _ []byte, _ uint32) {
			// silence — we only monitor existing system output
			for i := range output {
				output[i] = 0
			}
			w.outLevel.Store(math.Float64bits(calcRMS(output)))
		},
	})
	if err != nil {
		return err
	}
	w.outDevice = dev
	return dev.Start()
}

// ── DSP helpers ───────────────────────────────────────────────────────────────

func calcRMS(data []byte) float64 {
	n := len(data) / 2
	if n == 0 {
		return 0
	}
	var sum float64
	for i := 0; i < n*2; i += 2 {
		s := int16(binary.LittleEndian.Uint16(data[i:]))
		v := float64(s) / 32768.0
		sum += v * v
	}
	return math.Sqrt(sum / float64(n))
}

func rmsToDb(rms float64) float64 {
	if rms <= 0 {
		return -90
	}
	return 20 * math.Log10(rms)
}

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

// ── pactl OS calls ────────────────────────────────────────────────────────────

var volRe = regexp.MustCompile(`(\d+)%`)

func osGetVolume(kind, target string) (vol int, muted bool, err error) {
	out, err := exec.Command("pactl", "get-"+kind+"-volume", target).Output()
	if err != nil {
		return 0, false, fmt.Errorf("pactl get-%s-volume: %w", kind, err)
	}
	m := volRe.FindStringSubmatch(string(out))
	if m == nil {
		return 0, false, fmt.Errorf("pactl: could not parse volume: %q", out)
	}
	vol, _ = strconv.Atoi(m[1])

	muteOut, _ := exec.Command("pactl", "get-"+kind+"-mute", target).Output()
	muted = strings.Contains(string(muteOut), "yes")
	return vol, muted, nil
}

func osSetVolume(kind, target string, percent int) error {
	_, err := exec.Command("pactl",
		"set-"+kind+"-volume", target,
		strconv.Itoa(percent)+"%",
	).Output()
	return err
}

func osSetMute(kind, target string, mute bool) error {
	val := "0"
	if mute {
		val = "1"
	}
	_, err := exec.Command("pactl", "set-"+kind+"-mute", target, val).Output()
	return err
}
