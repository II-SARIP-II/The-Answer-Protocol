package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreateInteractionsPanel(topH, centerH, bottomH, colLeftW, colMidW, colRightW float32) *components.Panel {
	panelX := colLeftW + colMidW
	panelY := topH + centerH
	margin := float32(10.0)

	chat := components.NewChatWidget(
		panelX+margin,
		panelY+40.0,
		colRightW-(margin*2),
		bottomH-50.0,
		20.0,
		color.RGBA{R: 200, G: 200, B: 200, A: 255},
		"Action > ",
	)

	return &components.Panel{
		Name:        "Game Interactions",
		DisplayName: true,
		TextColor:   color.RGBA{R: 200, G: 200, B: 200, A: 255},
		ToDraw:      true,
		X:           panelX,
		Y:           panelY,
		W:           colRightW,
		H:           bottomH,
		Content:     []components.Widget{chat},
		BgColor:     color.RGBA{R: 50, G: 50, B: 50, A: 255},
		BorderColor: color.RGBA{R: 160, G: 160, B: 160, A: 255},
		TextSize:    float32(20.0),
	}
}
