package panels

import (
	"image/color"

	"tap/internal/gui/components"
)

func CreateActionsPanel(topH, centerH, colLeftW, colMidW, colRightW float32) *components.Panel {
	actionLabels := []string{
		"LOOK", "MOVE",
		"TAKE", "DROP",
		"TALK", "ATTACK",
		"STATUS", "QUEST",
		"QUESTS", "WHO",
		"GROUP", "QUIT",
	}

	crX := colLeftW + colMidW
	crY := topH
	crW := colRightW
	crH := centerH

	paddingX := float32(15)
	paddingY := float32(45)
	gapX := float32(10)
	gapY := float32(8)

	cols := 2
	rows := 6

	btnW := (crW - (paddingX * 2) - gapX) / float32(cols)
	btnH := (crH - paddingY - float32(20) - (gapY * float32(rows-1))) / float32(rows)

	var content []components.Widget

	for i, label := range actionLabels {
		col := i % cols
		row := i / cols

		bx := crX + paddingX + float32(col)*(btnW+gapX)
		by := crY + paddingY + float32(row)*(btnH+gapY)

		btn := components.NewButton(label, bx, by, btnW, btnH)

		content = append(content, btn)
	}

	return &components.Panel{
		Name:        "DIRECTIVE ACTIONS",
		DisplayName: true,
		TextColor:   color.RGBA{R: 0, G: 255, B: 200, A: 255},
		ToDraw:      true,
		X:           crX,
		Y:           crY,
		W:           crW,
		H:           crH,
		Content:     content,
		BgColor:     color.RGBA{R: 10, G: 16, B: 26, A: 240}, // Fond sombre Cyber
		BorderColor: color.RGBA{R: 0, G: 180, B: 220, A: 255},
		TextSize:    float32(20.0),
	}
}
