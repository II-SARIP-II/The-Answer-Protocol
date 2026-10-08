package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreateChatPanel(topH, centerH, bottomH, colLeftW, colMidW float32) *components.Panel {
	return &components.Panel{
		Name:			"Team Chat",
		DisplayName:	true,
		X:				colLeftW,
		Y:				topH + centerH,
		W:				colMidW,
		H:				bottomH,
		Content:		[]components.Widget{},
		BgColor:		color.RGBA{R: 160, G: 160, B: 160, A: 255},
		BorderColor:	color.RGBA{R: 50, G: 50, B: 50, A: 255},
		TextSize:		float32(25.0),
	}
}
