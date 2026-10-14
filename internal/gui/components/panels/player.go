package panels

import (
	"image/color"
	"log"

	"tap/internal/gui/components"
)

func CreatePlayerPanel(topH, centerH, colLeftW float32) *components.Panel {
	panelX := float32(0.0)
	panelY := topH

	margin := float32(20.0)
	paddingTop := float32(60.0)

	image, err := components.NewImage(
		"internal/gui/components/assets/agent.png",
		panelX+margin,
		panelY+paddingTop,
		colLeftW-(margin*2),
		centerH-paddingTop-margin,
	)

	if err != nil {
		log.Printf("[Warning] Player panel image load error: %v", err)
	}

	var widgetList []components.Widget
	if image != nil {
		widgetList = []components.Widget{image}
	}

	return &components.Panel{
		Name:        "AGENT PROFILE",
		DisplayName: true,
		TextColor:   color.RGBA{R: 0, G: 255, B: 200, A: 255},
		ToDraw:      true,
		X:           panelX,
		Y:           panelY,
		W:           colLeftW,
		H:           centerH,
		Content:     widgetList,
		BgColor:     color.RGBA{R: 10, G: 16, B: 26, A: 240},
		BorderColor: color.RGBA{R: 0, G: 180, B: 220, A: 255},
		TextSize:    float32(20.0),
	}
}
