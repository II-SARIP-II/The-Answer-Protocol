package components

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"

	"tap/internal/world"
)

const (
	slotHeight float32 = 40.0
	paddingTop float32 = 45.0
	paddingX   float32 = 15.0
	gapY       float32 = 8.0
)

type ItemSlot struct {
	Item   world.Item
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

func (s *Stuff) RecalculatePositions() {
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
		if slot.Item.Name == itemName {
			targetIdx = i
			break
		}
	}

	if targetIdx > 0 {
		targetSlot := s.Slots[targetIdx]

		remaining := append(s.Slots[:targetIdx], s.Slots[targetIdx+1:]...)
		s.Slots = append([]ItemSlot{targetSlot}, remaining...)
		s.RecalculatePositions()
	}
}

func (s *Stuff) AddItem(item world.Item, onItemClick func(string)) {
	label := fmt.Sprintf("%s (DMG: %d)", item.Name, item.Damage)
	btn := NewButton(label, 0, 0, 0, 0)

	itemName := item.Name
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
	s.RecalculatePositions()
}

func (s *Stuff) RemoveItem(item world.Item) {
	for i, slot := range s.Slots {
		if slot.Item.Name == item.Name {
			s.Slots = append(s.Slots[:i], s.Slots[i+1:]...)
			break
		}
	}

	s.RecalculatePositions()
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
