package panels

import (
	"image/color"

	"tap/internal/gui/components"
	"tap/internal/world"
)

func CreateMapPanel(topH, centerH, colLeftW, colMidW float32, data world.GameData) *components.Panel {
	var content []components.Widget

	const scale float32 = 30.0

	offsetX := colLeftW + 20
	offsetY := topH + 40

	for _, roomData := range data.Rooms {
		rx := offsetX + float32(roomData.X)*scale
		ry := offsetY + float32(roomData.Y)*scale
		rw := float32(roomData.W) * scale
		rh := float32(roomData.H) * scale

		roomWidget := &components.Panel{
			Name:			roomData.Name,
			DisplayName:	true,
			X:				rx,
			Y:				ry,
			W:				rw,
			H:				rh,
			Content:		nil,
			BgColor:		color.RGBA{R: 200, G: 200, B: 200, A: 255},
			BorderColor:	color.RGBA{R: 30, G: 30, B: 30, A: 255},
			TextSize:		float32(10.0),
		}

		content = append(content, roomWidget)
	}

	return &components.Panel{
		Name:			"MAP",
		DisplayName:	true,
		X:				colLeftW,
		Y:				topH,
		W:				colMidW,
		H:				centerH,
		Content:		content,
		BgColor:		color.RGBA{R: 120, G: 120, B: 120, A: 255},
		BorderColor:	color.RGBA{R: 50, G: 50, B: 50, A: 255},
		TextSize:		float32(25.0),
	}
}
