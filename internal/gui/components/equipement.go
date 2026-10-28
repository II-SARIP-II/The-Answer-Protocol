package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	slotHeight float32 = 40.0
	paddingTop float32 = 45.0
	paddingX   float32 = 15.0
	gapY       float32 = 8.0
)

type ItemSlot struct {
	Item   string
	Button *Button
}

type Stuff struct {
	X, Y, W, H float32
	Slots      []ItemSlot
}

func NewStuff(x, y, w, h float32) *Stuff {
	return &Stuff{
		X:     x,
		Y:     y,
		W:     w,
		H:     h,
		Slots: []ItemSlot{},
	}
}

func (s *Stuff) recalculatePositions() {
	for i := range s.Slots {
		itemY := s.Y + paddingTop + float32(i)*(slotHeight+gapY)
		s.Slots[i].Button.X = s.X + paddingX
		s.Slots[i].Button.Y = itemY
		s.Slots[i].Button.W = s.W - (paddingX * 2)
		s.Slots[i].Button.H = slotHeight
	}
}

func (s *Stuff) MoveToTop(itemName string) {
	targetIdx := -1
	for i, slot := range s.Slots {
		if slot.Item == itemName {
			targetIdx = i
			break
		}
	}

	if targetIdx > 0 {
		targetSlot := s.Slots[targetIdx]

		remaining := append(s.Slots[:targetIdx], s.Slots[targetIdx+1:]...)
		s.Slots = append([]ItemSlot{targetSlot}, remaining...)
		s.recalculatePositions()
	}
	for i, item := range s.Slots {
		if i == 0 {
			item.Button.NormalColor = color.RGBA{R: 0, G: 180, B: 216, A: 180}
			item.Button.TextColor = color.RGBA{R: 0, G: 0, B: 0, A: 180}

		} else {
			item.Button.NormalColor = color.RGBA{R: 12, G: 24, B: 36, A: 220}
			item.Button.TextColor = color.RGBA{R: 0, G: 180, B: 216, A: 180}
		}
	}
}

func (s *Stuff) AddItem(item string, onItemClick func(string)) {
	btn := NewButton(item, 0, 0, 0, 0)

	itemName := item
	btn.OnClick = func() {
		if onItemClick != nil {
			onItemClick(itemName)
		}
	}

	itemSlot := ItemSlot{
		Item:   item,
		Button: btn,
	}

	s.Slots = append(s.Slots, itemSlot)
	s.recalculatePositions()
}

func (s *Stuff) RemoveItem(item string) {
	for i, slot := range s.Slots {
		if slot.Item == item {
			s.Slots = append(s.Slots[:i], s.Slots[i+1:]...)
			break
		}
	}

	s.recalculatePositions()
}

func (s *Stuff) Update() {
	for i := range s.Slots {
		s.Slots[i].Button.Update()
	}
}

func (s *Stuff) Draw(screen *ebiten.Image) {
	for i := range s.Slots {
		s.Slots[i].Button.Draw(screen)
	}
}
