package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreateStuffPanel(topH, centerH, bottomH, colLeftW float32) *components.Panel {
	return &components.Panel{
		Name:        "STUFF",
		X:           colLeftW * 0.5,
		Y:           topH + centerH,
		W:           colLeftW * 0.5,
		H:           bottomH,
		BgColor:     color.RGBA{R: 160, G: 160, B: 160, A: 255},
		BorderColor: color.RGBA{R: 50, G: 50, B: 50, A: 255},
	}
}
