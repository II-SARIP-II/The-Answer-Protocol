package components

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type PlayerInvite struct {
	Name string
	Btn  *Button
}

func NewPlayerInvite(name string, index int, panelX, panelY, panelW float32) *PlayerInvite {
	paddingX := float32(10.0)
	paddingY := float32(35.0)
	gapY := float32(6.0)

	btnW := panelW - (paddingX * 2)
	btnH := float32(35.0)

	xPos := panelX + paddingX
	yPos := panelY + paddingY + float32(index)*(btnH+gapY)

	btn := NewButton(name, xPos, yPos, btnW, btnH)
	return &PlayerInvite{
		Name: name,
		Btn:  btn,
	}
}

func CreateList(name string, playerList []string, panelX, panelY, panelW float32) []*PlayerInvite {
	listPlayerInvite := make([]*PlayerInvite, 0, len(playerList))
	idx := 0
	for _, player := range playerList {
		if player == name {
			continue
		}
		inv := NewPlayerInvite(player, idx, panelX, panelY, panelW)
		listPlayerInvite = append(listPlayerInvite, inv)
		idx++
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
