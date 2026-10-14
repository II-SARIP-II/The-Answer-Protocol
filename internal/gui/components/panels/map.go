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
		ToDraw := false
		var holes []components.Hole = []components.Hole{}
		rx := offsetX + float32(roomData.X)*scale
		ry := offsetY + float32(roomData.Y)*scale
		rw := float32(roomData.W) * scale
		rh := float32(roomData.H) * scale
		for _, exitData := range roomData.Exits {
			holes = append(holes, components.Hole{
				Wall: exitData.Wall,
				Pos:  float32(exitData.DoorPos) * scale,
			})
		}

		if roomData.Name == "Hotel Entrance" {
			ToDraw = true
		}

		roomWidget := &components.Panel{
			Name:        roomData.Name,
			DisplayName: true,
			TextColor:   color.RGBA{R: 0, G: 255, B: 180, A: 255},
			ToDraw:      ToDraw,
			X:           rx,
			Y:           ry,
			W:           rw,
			H:           rh,
			Content:     nil,
			BorderHoles: holes,
			BgColor:     color.RGBA{R: 15, G: 35, B: 45, A: 200},
			BorderColor: color.RGBA{R: 0, G: 255, B: 180, A: 255},
			TextSize:    float32(10.0),
		}

		content = append(content, roomWidget)
	}

	return &components.Panel{
		Name:        "TACTICAL MAP // RADAR",
		DisplayName: true,
		TextColor:   color.RGBA{R: 0, G: 255, B: 200, A: 255},
		ToDraw:      true,
		X:           colLeftW,
		Y:           topH,
		W:           colMidW,
		H:           centerH,
		Content:     content,
		BgColor:     color.RGBA{R: 8, G: 12, B: 20, A: 240},
		BorderColor: color.RGBA{R: 0, G: 180, B: 220, A: 255},
		TextSize:    float32(20.0),
	}
}
