package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreateTopPanel(w, topH float32) *components.Panel {
	return &components.Panel{
		Name:			"Top Panel",
		DisplayName:	true,
		X:				0,
		Y:				0,
		W:				w,
		H:				topH,
		Content:		[]components.Widget{},
		BgColor:		color.RGBA{R: 180, G: 180, B: 180, A: 255},
		BorderColor:	color.RGBA{R: 50, G: 50, B: 50, A: 255},
		TextSize:		float32(25.0),
	}
}
