package sysmonitor

const (
	filled = "⠿"
	empty  = " "
)

func renderBar(percent float64, maxPercent float64, length int) string {
	filledCount := int(percent / maxPercent * float64(length))
	var bar string
	for i := range length {
		if i < filledCount {
			bar += filled
			continue
		}
		bar += empty
	}
	return bar
}
