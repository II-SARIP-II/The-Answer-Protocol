package components

import (
	"image/color"
	"tap/internal/world"
)

const (
	roomW int = 3
	roomH int = 3
)

type Door struct {
	Dir int
	Pos int
}

func GetHolesForRoom(Exits map[string]world.Exit, scale float32) []Hole {
	var door []Hole
	for exit := range Exits {
		exitPos := Hole{}
		switch exit {
		case "north":
			exitPos.Wall = 0
			exitPos.Pos = 2 * scale
		case "east":
			exitPos.Wall = 1
			exitPos.Pos = 1 * scale
		case "south":
			exitPos.Wall = 2
			exitPos.Pos = 2 * scale
		case "west":
			exitPos.Wall = 3
			exitPos.Pos = 1 * scale
		}
		door = append(door, exitPos)
	}
	return door
}

func unpack2(arr [2]int) (int, int) {
	return arr[0], arr[1]
}

func (p *Panel) RemoveRooms() {
	if p == nil {
		return
	}
	var cleanContent []Widget
	for _, widget := range p.Content {
		if btn, isButton := widget.(*Button); isButton {
			switch btn.Label {
			case "North", "East", "South", "West":
				cleanContent = append(cleanContent, widget)
			default:
			}
		} else if _, isPanel := widget.(*Panel); isPanel {
			continue
		} else {
			cleanContent = append(cleanContent, widget)
		}
	}

	p.Content = cleanContent
}

func (p *Panel) AddRoom(name string, location [2]int, exits map[string]world.Exit, items []string, onItemClick func(string), npcs []string, onNPCClick func(string)) {
	if p == nil {
		return
	}

	const scale float32 = 30.0
	const padding float32 = 15.0
	gridX, gridY := unpack2(location)

	rx := p.X + float32(gridX)*scale*float32(roomW)
	ry := p.Y + float32(gridY)*scale*float32(roomH)
	rw := float32(roomW) * scale
	rh := float32(roomH) * scale

	holes := GetHolesForRoom(exits, scale)

	var roomPanel *Panel
	for _, widget := range p.Content {
		if panel, ok := widget.(*Panel); ok && panel.Name == name {
			roomPanel = panel
			break
		}
	}

	if roomPanel == nil {
		roomPanel = &Panel{
			Name:        name,
			DisplayName: true,
			TextColor:   color.RGBA{R: 0, G: 255, B: 180, A: 255},
			ToDraw:      true,
			X:           rx,
			Y:           ry,
			W:           rw,
			H:           rh,
			BorderHoles: holes,
			BgColor:     color.RGBA{R: 15, G: 35, B: 45, A: 200},
			BorderColor: color.RGBA{R: 0, G: 255, B: 180, A: 255},
			TextSize:    float32(10.0),
		}
		p.Content = append(p.Content, roomPanel)
	} else {
		roomPanel.X = rx
		roomPanel.Y = ry
		roomPanel.W = rw
		roomPanel.H = rh
		roomPanel.BorderHoles = holes
		roomPanel.ToDraw = true
	}

	var cleanContent []Widget
	for _, widget := range p.Content {
		if btn, isButton := widget.(*Button); isButton {
			switch btn.Label {
			case "North", "East", "South", "West":
				cleanContent = append(cleanContent, widget)
			default:
			}
		} else {
			cleanContent = append(cleanContent, widget)
		}
	}

	p.Content = cleanContent

	// add item in room
	for i, item := range items {
		btnW := float32(175)
		btnY := p.Y + p.H - float32((i+1)*50) - padding*float32(i)
		btn := NewButton(item, p.X+p.W-btnW-padding, btnY-padding, btnW, 50)

		itemName := item
		btn.OnClick = func() {
			if onItemClick != nil {
				onItemClick(itemName)
			}
		}

		p.Content = append(p.Content, btn)
	}

	// add npc in room
	for i, npc := range npcs {
		btnW := float32(175)
		btnY := p.Y + p.H - float32((i+1)*50) - padding*float32(i)
		btn := NewButton(npc, p.X+p.W/2-btnW/2, btnY-padding, btnW, 50)

		npcName := npc
		btn.OnClick = func() {
			if onNPCClick != nil {
				onNPCClick(npcName)
			}
		}

		p.Content = append(p.Content, btn)
	}

}
