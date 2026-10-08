package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreatePlayerPanel(topH, centerH, colLeftW float32) *components.Panel {
	return &components.Panel{
		Name:			"Player",
		DisplayName:	true,
		ToDraw:			true,
		X:				0,
		Y:				topH,
		W:				colLeftW,
		H:				centerH,
		Content:		[]components.Widget{},
		BgColor:		color.RGBA{R: 160, G: 160, B: 160, A: 255},
		BorderColor:	color.RGBA{R: 50, G: 50, B: 50, A: 255},
		TextSize:		float32(25.0),
	}
}
