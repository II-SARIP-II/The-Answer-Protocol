package panels

import (
	"fmt"
	"image/color"

	"tap/internal/gui/components"
)

func CreateActionsPanel(topH, centerH, colLeftW, colMidW, colRightW float32) *components.Panel {
	actions := map[string]func(){
		"LOOK":   Look,
		"MOVE":   Move,
		"TAKE":   Take,
		"DROP":   Drop,
		"TALK":   Talk,
		"ATTACK": Attack,
		"STATUS": Status,
		"QUEST":  Quest,
		"QUESTS": Quests,
		"WHO":    Who,
		"GROUP":  Group,
		"QUIT":   Quit,
	}
	
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

		actionFunc := actions[label]
		if actionFunc == nil {
			actionFunc = func() { fmt.Printf("Action Not Found : %s\n", label) }
		}

		btn := components.NewButton(label, bx, by, btnW, btnH, actionFunc)

		content = append(content, btn)
	}

	return &components.Panel{
		Name:			"ACTIONS",
		DisplayName:	true,
		X:				crX,
		Y:				crY,
		W:				crW,
		H:				crH,
		Content:		content,
		BgColor:		color.RGBA{R: 160, G: 160, B: 160, A: 255},
		BorderColor:	color.RGBA{R: 50, G: 50, B: 50, A: 255},
		TextSize:		float32(25.0),
	}
}

func Look() {
	fmt.Print("LOOK clicked\n")
}

func Move() {
	fmt.Print("MOVE clicked\n")
}

func Take() {
	fmt.Print("TAKE clicked\n")
}

func Drop() {
	fmt.Print("DROP clicked\n")
}

func Talk() {
	fmt.Print("TALK clicked\n")
}

func Attack() {
	fmt.Print("ATTACK clicked\n")
}

func Status() {
	fmt.Print("STATUS clicked\n")
}

func Quest() {
	fmt.Print("QUEST clicked\n")
}

func Quests() {
	fmt.Print("QUESTS clicked\n")
}

func Who() {
	fmt.Print("WHO clicked\n")
}

func Group() {
	fmt.Print("GROUP clicked\n")
}

func Quit() {
	fmt.Print("QUIT clicked\n")
}
