package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreateStatsPanel(topH, centerH, bottomH, colLeftW float32) *components.Panel {
	panelX := float32(0)
	panelY := topH + centerH

	stuff := components.NewStuff(
		panelX,
		panelY,
		colLeftW,
		bottomH,
	)

	return &components.Panel{
		Name:        "EQUIPMENT / INVENTORY",
		DisplayName: true,
		TextColor:   color.RGBA{R: 0, G: 255, B: 200, A: 255},
		ToDraw:      true,
		X:           panelX,
		Y:           panelY,
		W:           colLeftW,
		H:           bottomH,
		Content:     []components.Widget{stuff},
		BgColor:     color.RGBA{R: 10, G: 16, B: 26, A: 240},
		BorderColor: color.RGBA{R: 0, G: 180, B: 220, A: 255},
		TextSize:    float32(20.0),
	}
}

