package panels

import (
	"image/color"

	"tap/internal/gui/components"
	"tap/internal/world"
)

func CreateMapPanel(topH, centerH, colLeftW, colMidW float32, data world.GameData) *components.Panel {
	var content []components.Widget

	const scale float32 = 30.0
	const compassW float32 = 60
	const compassH float32 = 60
	
	compassRoseNorth := components.NewButton("North", colLeftW + colMidW/2 - compassW/2, topH + centerH - 3 * compassH, compassW, compassH)
	compassRoseEast := components.NewButton("West", colLeftW + colMidW/2 - compassW/2 - compassW, topH + centerH - 2 * compassH, compassW, compassH)
	compassRoseSouth := components.NewButton("South", colLeftW + colMidW/2 - compassW/2, topH + centerH - 1 * compassH, compassW, compassH)
	compassRoseWest := components.NewButton("East", colLeftW + colMidW/2 - compassW/2 + compassW, topH + centerH - 2 * compassH, compassW, compassH)
	content = append(content, compassRoseNorth)
	content = append(content, compassRoseEast)
	content = append(content, compassRoseSouth)
	content = append(content, compassRoseWest)
	return &components.Panel{
		Name:        "Hotel MAP",
		DisplayName: true,
		TextColor:   color.RGBA{R: 0, G: 255, B: 200, A: 255},
		ToDraw:      true,
		X:           colLeftW,
		Y:           topH,
		W:           colMidW,
		H:           centerH,
		Content:     content,
		BgColor:     color.RGBA{R: 10, G: 16, B: 26, A: 240},
		BorderColor: color.RGBA{R: 0, G: 180, B: 220, A: 255},
		TextSize:    float32(20.0),
	}
}
