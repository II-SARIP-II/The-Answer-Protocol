package components

import (
	"image/color"
	"tap/internal/world"
)

const (
	roomW int = 4
	roomH int = 3
)

type Door struct {
	Dir int
	Pos int
}

var mapPos = map[string][2]int{
	"Hotel Entrance":  {1, 1},
	"Broom Closet":    {1, 2},
	"Hotel Reception": {2, 1},
	"Hotel Kitchen":   {3, 1},
	"Corridor 1":      {2, 2},
	"Dining Room":     {3, 2},
	"Corridor 2":      {2, 3},
	"VIP Room":        {3, 3},
	"Secret Corridor": {3, 4},
	"Secret Room":     {4, 4},
}

func GetHolesForRoom(Exits map[string]world.Exit, scale float32) []Hole {
	var door []Hole
	for exit := range Exits {
		exitPos := Hole{}
		switch exit {
		case "north":
			exitPos.Wall = 0
			exitPos.Pos = 2
		case "east":
			exitPos.Wall = 1
			exitPos.Pos = 1
		case "south":
			exitPos.Wall = 2
			exitPos.Pos = 2
		case "west":
			exitPos.Wall = 3
			exitPos.Pos = 1
		}
		door = append(door, exitPos)
	}
	return door
}

func GetRoomPos(name string) (int, int) {
	if pos, ok := mapPos[name]; ok {
		return pos[0], pos[1]
	}
	return 1, 1
}

func (p *Panel) AddRoom(name string, exits map[string]world.Exit, items []string, onItemClick func(string), npcs []string, onNPCClick func(string)) {
	if p == nil {
		return
	}

	const scale float32 = 30.0
	const padding float32 = 15.0
	gridX, gridY := GetRoomPos(name)

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
		btnY := p.Y + p.H - float32((i+1)*50)
		btn := NewButton(item, p.X+p.W-150, btnY, 150, 50)

		btn.OnClick = func() {
			if onItemClick != nil {
				onItemClick(item)
			}
		}

		p.Content = append(p.Content, btn)
	}

	// add npc in room
	for i, npc := range npcs {
		btnY := p.Y + p.H - float32((i+1)*50)
		btn := NewButton(npc, p.X + padding, btnY, 175, 50)

		btn.OnClick = func() {
			if onItemClick != nil {
				onItemClick(npc)
			}
		}

		p.Content = append(p.Content, btn)
	}

}
