package components

import (
	"fmt"
	"image/color"
	"tap/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
)

type Stuff struct {
	X, Y, W, H float32
	Items      []world.Item
}

func NewStuff(x, y, w, h float32) *Stuff {
	return &Stuff{
		X:     x,
		Y:     y,
		W:     w,
		H:     h,
		Items: []world.Item{},
	}
}

func (s *Stuff) AddItem(item world.Item) {
	s.Items = append(s.Items, item)
}

func (s *Stuff) RemoveItem(item world.Item) {
	s.Items = append(s.Items, item)
}

func (s *Stuff) Update() {
}

func (s *Stuff) Draw(screen *ebiten.Image) {
	const lineHeight float32 = 22.0
	const paddingTop float32 = 40.0

	for i, item := range s.Items {
		msg := fmt.Sprintf(" - %s (DMG: %d)", item.Name, item.Damage)

		itemY := s.Y + paddingTop + float32(i)*lineHeight

		DrawText(
			screen,
			msg,
			float64(s.X+15),
			float64(itemY),
			20.0,
			color.RGBA{R: 0, G: 255, B: 180, A: 255},
			"default",
		)
	}
}
