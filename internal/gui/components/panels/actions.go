package panels

import (
	"fmt"
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

	var actionButtons []*components.Button
	for i, label := range actionLabels {
		col := i % cols
		row := i / cols

		bx := crX + paddingX + float32(col)*(btnW+gapX)
		by := crY + paddingY + float32(row)*(btnH+gapY)

		btnLabel := label
		btn := components.NewButton(label, bx, by, btnW, btnH, func() {
			fmt.Printf("Action cliquée : %s\n", btnLabel)
		})

		actionButtons = append(actionButtons, btn)
	}

	return &components.Panel{
		Name:        "ACTIONS",
		X:           crX,
		Y:           crY,
		W:           crW,
		H:           crH,
		Buttons:     actionButtons,
		BgColor:     color.RGBA{R: 160, G: 160, B: 160, A: 255},
		BorderColor: color.RGBA{R: 50, G: 50, B: 50, A: 255},
	}
}
