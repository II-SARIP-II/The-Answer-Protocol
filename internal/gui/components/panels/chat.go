package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreateChatPanel(topH, centerH, bottomH, colLeftW, colMidW float32) *components.Panel {
	panelX := colLeftW
	panelY := topH + centerH
	margin := float32(10.0)

	chat := components.NewChatWidget(
		panelX+margin,
		panelY+40.0,
		colMidW-(margin*2),
		bottomH-50.0,
		20.0,
		color.RGBA{R: 200, G: 200, B: 200, A: 255},
		"Team",
		true,
	)

	return &components.Panel{
		Name:        "Team Chat",
		DisplayName: true,
		TextColor:   color.RGBA{R: 200, G: 200, B: 200, A: 255},
		ToDraw:      true,
		X:           colLeftW,
		Y:           topH + centerH,
		W:           colMidW,
		H:           bottomH,
		Content:     []components.Widget{chat},
		BgColor:     color.RGBA{R: 50, G: 50, B: 50, A: 255},
		BorderColor: color.RGBA{R: 160, G: 160, B: 160, A: 255},
		TextSize:    float32(25.0),
	}
}
