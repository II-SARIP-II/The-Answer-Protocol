package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreateMapPanel(topH, centerH, colLeftW, colMidW float32) *components.Panel {
	return &components.Panel{
		Name:        "MAP",
		X:           colLeftW,
		Y:           topH,
		W:           colMidW,
		H:           centerH,
		BgColor:     color.RGBA{R: 160, G: 160, B: 160, A: 255},
		BorderColor: color.RGBA{R: 50, G: 50, B: 50, A: 255},
	}
}
