package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreatePlayerPanel(topH, centerH, colLeftW float32) *components.Panel {
	panelX := float32(0.0)
	panelY := topH

	marginX := float32(15.0)
	marginY := float32(15.0)
	btnW := float32(120.0)
	btnH := float32(40.0)

	btnCreate := components.NewButton("CREATE", panelX+marginX, panelY+marginY+15, btnW, btnH)
	btnQuit := components.NewButton("LEAVE", panelX+colLeftW-marginX-btnW, panelY+marginY+15, btnW, btnH)

	colsY := panelY + marginY + btnH + 30.0
	colsH := centerH - (colsY - panelY) - marginY
	gap := float32(10.0)
	colW := (colLeftW - (marginX * 2) - gap) / 2.0

	colMembers := &components.Panel{
		Name:        "INVITE PLAYER",
		DisplayName: true,
		TextColor:   color.RGBA{R: 0, G: 255, B: 180, A: 255},
		ToDraw:      true,
		X:           panelX + marginX,
		Y:           colsY,
		W:           colW,
		H:           colsH,
		Content:     []components.Widget{},
		BgColor:     color.RGBA{R: 15, G: 25, B: 38, A: 220},
		BorderColor: color.RGBA{R: 0, G: 180, B: 220, A: 180},
		TextSize:    float32(14.0),
	}

	colInvites := &components.Panel{
		Name:        "INVITATION",
		DisplayName: true,
		TextColor:   color.RGBA{R: 0, G: 255, B: 180, A: 255},
		ToDraw:      true,
		X:           panelX + marginX + colW + gap,
		Y:           colsY,
		W:           colW,
		H:           colsH,
		Content:     []components.Widget{},
		BgColor:     color.RGBA{R: 15, G: 25, B: 38, A: 220},
		BorderColor: color.RGBA{R: 0, G: 180, B: 220, A: 180},
		TextSize:    float32(14.0),
	}

	widgets := []components.Widget{
		btnCreate,
		btnQuit,
		colMembers,
		colInvites,
	}

	return &components.Panel{
		Name:          "GROUP",
		DisplayName:   true,
		TextColor:     color.RGBA{R: 0, G: 255, B: 200, A: 255},
		ToDraw:        true,
		X:             panelX,
		Y:             panelY,
		W:             colLeftW,
		H:             centerH,
		Content:       widgets,
		BgColor:       color.RGBA{R: 10, G: 16, B: 26, A: 240},
		BorderColor:   color.RGBA{R: 0, G: 180, B: 220, A: 255},
		TextSize:      float32(20.0),
		PlayerInvites: []*components.PlayerInvite{},
	}
}
