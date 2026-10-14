package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreateTopPanel(w, topH float32) *components.Panel {
	return &components.Panel{
		Name:        "Top Panel",
		DisplayName: true,
		TextColor:   color.RGBA{R: 200, G: 200, B: 200, A: 255},
		ToDraw:      true,
		X:           0,
		Y:           0,
		W:           w,
		H:           topH,
		Content:     []components.Widget{},
		BgColor:     color.RGBA{R: 50, G: 50, B: 50, A: 255},
		BorderColor: color.RGBA{R: 160, G: 160, B: 160, A: 255},
		TextSize:    float32(40.0),
	}
}
