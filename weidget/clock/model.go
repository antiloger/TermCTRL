package clock

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/antiloger/termctlr/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ClockModel struct {
	w    types.Position
	size types.Position
	ct   time.Time
}

func NewClockWidget() ClockModel {
	return ClockModel{
		ct: time.Now(),
	}
}

func (C ClockModel) Init() tea.Cmd {
	return nil
}

func (C ClockModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		C.w.X = msg.Width
		C.w.Y = msg.Height
	case types.TickMsg:
		C.ct = time.Time(msg)
		return C, nil // Keep ticking
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return C, tea.Quit
		}
	}
	return C, nil
}

func (C ClockModel) SetPosition(x, y int) {
	C.w.X = x
	C.w.Y = y
}

func (C ClockModel) View() string {
	// Get current time
	hour := C.ct.Hour()
	minute := C.ct.Minute()
	second := C.ct.Second()
	// day := C.ct.Day()
	// year := C.ct.Year()
	// month := C.ct.Month().String()[:3] // Short month name
	// dayOfWeek := C.ct.Weekday().String()[:3]
	timezone, _ := C.ct.Zone()

	// Convert to digits
	h1 := hour / 10
	h2 := hour % 10
	m1 := minute / 10
	m2 := minute % 10
	s1 := second / 10
	s2 := second % 10

	// Build the display: HH:MM:SS
	// display := joinHorizontal(
	// 	asciiDigits[h1],
	// 	asciiDigits[h2],
	// 	asciiColon,
	// 	asciiDigits[m1],
	// 	asciiDigits[m2],
	// 	asciiColon,
	// 	asciiDigits[s1],
	// 	asciiDigits[s2],
	// )

	display := lipgloss.JoinHorizontal(lipgloss.Center,
		asciiDigits[h1],
		asciiDigits[h2],
		asciiColon,
		asciiDigits[m1],
		asciiDigits[m2],
		asciiColon,
		asciiDigits[s1],
		asciiDigits[s2],
	)

	display = lipgloss.JoinHorizontal(lipgloss.Center, display, "      ", TestCal2)

	style := lipgloss.NewStyle().
		Foreground(lipgloss.Color("ff")).
		Bold(true)

	timeStyle := lipgloss.NewStyle()
	// datestlyle := lipgloss.NewStyle().
	// 	Bold(true).
	// 	Italic(true)
	// dateStr := fmt.Sprintf(" %s %02d, %d || %s ", month, day, year, dayOfWeek)
	bottomInfo := fmt.Sprintf("TZ: %s", timezone)
	return lipgloss.Place(
		C.w.X,
		C.w.Y,
		lipgloss.Center, // horizontal center
		lipgloss.Center, // vertical center
		style.Render(timeStyle.Render(display)+"\n"+bottomInfo),
	)
}

func joinHorizontal(arts ...string) string {
	// Split each art into lines
	var allLines [][]string
	maxHeight := 0
	widths := make([]int, len(arts))
	for i, art := range arts {
		lines := strings.Split(strings.TrimPrefix(art, "\n"), "\n")
		allLines = append(allLines, lines)
		if len(lines) > maxHeight {
			maxHeight = len(lines)
		}
		// Calculate max width for this block using rune count
		for _, line := range lines {
			runeCount := utf8.RuneCountInString(line)
			if runeCount > widths[i] {
				widths[i] = runeCount
			}
		}
	}
	// Join line by line
	var result strings.Builder
	for i := 0; i < maxHeight; i++ {
		for j, lines := range allLines {
			if i < len(lines) {
				result.WriteString(lines[i])
				// Pad remaining space if not the last block
				if j < len(allLines)-1 {
					padding := widths[j] - utf8.RuneCountInString(lines[i])
					result.WriteString(strings.Repeat(" ", padding))
				}
			} else {
				// No line at this index, write spaces for width
				if j < len(allLines)-1 {
					result.WriteString(strings.Repeat(" ", widths[j]))
				}
			}
		}
		if i < maxHeight-1 {
			result.WriteString("\n")
		}
	}
	return result.String()
}
