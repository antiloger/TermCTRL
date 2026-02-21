package audio

const (
	FildBox  = "█"
	EmptyBox = "░"
	MutedBox = "░░░░░░░ mute ░░░░░░░"
)

func (m *Model) UIVolumeOut() string {
	if m.audio.OutMuted {
		return MutedBox
	}
	maxblock := 20

	filled := int(float64(m.audio.OutVolume) / float64(m.audio.MaxInVolume) * float64(maxblock))
	var bar string
	for i := range maxblock {
		if i < filled {
			bar += FildBox
			continue
		}
		bar += EmptyBox
	}

	return bar
}

func (m *Model) UIVolumeIn() string {
	if m.audio.InMuted {
		return MutedBox
	}
	maxblock := 20

	filled := int(float64(m.audio.InVolume) / float64(m.audio.MaxInVolume) * float64(maxblock))
	var bar string
	for i := range maxblock {
		if i < filled {
			bar += FildBox
			continue
		}
		bar += EmptyBox
	}

	return bar
}
