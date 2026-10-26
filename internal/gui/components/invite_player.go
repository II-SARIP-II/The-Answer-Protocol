package components

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerInvite struct {
	Name string
	Btn  *Button
}

func NewPlayerInvite(name string, index int, panelX, panelY, panelW float32) *PlayerInvite {
	paddingX := float32(15.0)
	paddingY := float32(100.0)
	gapX := float32(10.0)
	gapY := float32(8.0)

	cols := 2
	col := index % cols
	row := index / cols

	btnW := (panelW - (paddingX * 2) - gapX) / float32(cols)
	btnH := float32(40.0)

	xPos := panelX + paddingX + float32(col)*(btnW+gapX)
	yPos := panelY + paddingY + float32(row)*(btnH+gapY)

	btn := NewButton(name, xPos, yPos, btnW, btnH)
	return &PlayerInvite{
		Name: name,
		Btn:  btn,
	}
}

func CreateList(playerList []string, panelX, panelY, panelW float32) []*PlayerInvite {
	listPlayerInvite := make([]*PlayerInvite, 0, len(playerList))
	for i, player := range playerList {
		inv := NewPlayerInvite(player, i, panelX, panelY, panelW)
		listPlayerInvite = append(listPlayerInvite, inv)
	}
	return listPlayerInvite
}

func DrawList(screen *ebiten.Image, playerList []*PlayerInvite) {
	for _, player := range playerList {
		if player != nil && player.Btn != nil {
			player.Btn.Draw(screen)
		}
	}
}
