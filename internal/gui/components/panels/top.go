package panels

import (
	"fmt"
	"image/color"

	"tap/internal/gui/components"
)

func CreateTopPanel(w, topH float32, playerInGame, playerInRoom int) *components.Panel {
	textColor := color.RGBA{R: 0, G: 255, B: 180, A: 255}

	inGameText := components.NewTextWidget(
		20,
		15,
		20,
		fmt.Sprintf("Players in Game: %d", playerInGame),
		textColor,
	)

	inRoomText := components.NewTextWidget(
		250,
		15,
		20,
		fmt.Sprintf("Players in Room: %d", playerInRoom),
		textColor,
	)

	top := &components.Panel{
		Name:        "THE ANSWER PROTOCOL",
		DisplayName: true,
		TextColor:   color.RGBA{R: 0, G: 255, B: 200, A: 255},
		ToDraw:      true,
		X:           0,
		Y:           0,
		W:           w,
		H:           topH,
		Content:     []components.Widget{inGameText, inRoomText},
		BgColor:     color.RGBA{R: 10, G: 16, B: 26, A: 240},
		BorderColor: color.RGBA{R: 0, G: 180, B: 220, A: 255},
		TextSize:    float32(40.0),
	}

	return top
}
