package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreateInteractionsPanel(topH, centerH, bottomH, colLeftW, colMidW, colRightW float32) *components.Panel {
	return &components.Panel{
		Name:			"Game Interactions",
		DisplayName:	true,
		ToDraw:			true,
		X:				colLeftW + colMidW,
		Y:				topH + centerH,
		W:				colRightW,
		H:				bottomH,
		Content:		[]components.Widget{},
		BgColor:		color.RGBA{R: 160, G: 160, B: 160, A: 255},
		BorderColor:	color.RGBA{R: 50, G: 50, B: 50, A: 255},
		TextSize:		float32(25.0),
	}
}
