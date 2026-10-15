package panels

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"tap/internal/gui/components"
	"tap/internal/world"
)

type Stuff struct {
	X, Y, W, H float32
	items      []world.Item
}

func (s *Stuff) NewItem(item world.Item) {
	s.items = append(s.items, item)
}

func (s *Stuff) Update() {
}

func (s *Stuff) Draw(screen *ebiten.Image) {
	const lineHeight float32 = 22.0
	const paddingTop float32 = 40.0

	for i, item := range s.items {
		msg := fmt.Sprintf(" - %s (DMG: %d)", item.Name, item.Damage)

		itemY := s.Y + paddingTop + float32(i)*lineHeight

		components.DrawText(
			screen,
			msg,
			float64(s.X+15),
			float64(itemY),
			14.0,
			color.RGBA{R: 0, G: 255, B: 180, A: 255},
			"default",
		)
	}
}

func CreateStatsPanel(topH, centerH, bottomH, colLeftW float32) *components.Panel {
	panelX := float32(0)
	panelY := topH + centerH

	stuff := &Stuff{
		X:     panelX,
		Y:     panelY,
		W:     colLeftW,
		H:     bottomH,
		items: []world.Item{},
	}

	item1 := world.Item{
		ID:     1,
		Name:   "couteau",
		Damage: 5,
	}
	stuff.NewItem(item1)

	return &components.Panel{
		Name:        "EQUIPMENT / INVENTORY",
		DisplayName: true,
		TextColor:   color.RGBA{R: 0, G: 255, B: 200, A: 255},
		ToDraw:      true,
		X:           panelX,
		Y:           panelY,
		W:           colLeftW,
		H:           bottomH,
		Content:     []components.Widget{stuff},
		BgColor:     color.RGBA{R: 10, G: 16, B: 26, A: 240},
		BorderColor: color.RGBA{R: 0, G: 180, B: 220, A: 255},
		TextSize:    float32(20.0),
	}
}
