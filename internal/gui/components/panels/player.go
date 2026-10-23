package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreatePlayerPanel(topH, centerH, colLeftW float32) *components.Panel {
	panelX := float32(0.0)
	panelY := topH

	marginX := float32(25.0)
	marginY := float32(20.0)

	btn := components.NewButton("CREATE GROUP", panelX + marginX, panelY + float32(marginY) + 20, 150, 50)

	return &components.Panel{
		Name:        "AGENT PROFILE",
		DisplayName: true,
		TextColor:   color.RGBA{R: 0, G: 255, B: 200, A: 255},
		ToDraw:      true,
		X:           panelX,
		Y:           panelY,
		W:           colLeftW,
		H:           centerH,
		Content:     []components.Widget{btn},
		BgColor:     color.RGBA{R: 10, G: 16, B: 26, A: 240},
		BorderColor: color.RGBA{R: 0, G: 180, B: 220, A: 255},
		TextSize:    float32(20.0),
	}
}
