package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Button struct {
	X, Y, W, H float32
	Label      string
	IsHovered  bool
	IsPressed  bool

	NormalColor  color.RGBA
	HoverColor   color.RGBA
	PressedColor color.RGBA
	BorderColor  color.RGBA

	OnClick func()
}

func NewButton(label string, x, y, w, h float32, onClick func()) *Button {
	return &Button{
		X:            x,
		Y:            y,
		W:            w,
		H:            h,
		Label:        label,
		NormalColor:  color.RGBA{R: 45, G: 55, B: 75, A: 255},
		HoverColor:   color.RGBA{R: 70, G: 95, B: 135, A: 255},
		PressedColor: color.RGBA{R: 25, G: 35, B: 50, A: 255},
		BorderColor:  color.RGBA{R: 120, G: 160, B: 215, A: 255},
		OnClick:      onClick,
	}
}

func (b *Button) Update() {
	mx, my := ebiten.CursorPosition()
	fx, fy := float32(mx), float32(my)

	b.IsHovered = (fx >= b.X && fx <= b.X+b.W && fy >= b.Y && fy <= b.Y+b.H)

	if b.IsHovered {
		b.IsPressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if b.OnClick != nil {
				b.OnClick()
			}
		}
	} else {
		b.IsPressed = false
	}
}

func (b *Button) Draw(screen *ebiten.Image) {
	currentColor := b.NormalColor
	if b.IsPressed {
		currentColor = b.PressedColor
	} else if b.IsHovered {
		currentColor = b.HoverColor
	}

	DrawFilledRect(screen, b.X, b.Y, b.W, b.H, currentColor)

	borderWidth := float32(1)
	if b.IsHovered {
		borderWidth = 2
	}
	var hole []Hole = []Hole{}
	DrawStrokeRect(screen, b.X, b.Y, b.W, b.H, borderWidth, b.BorderColor, hole)

	textX := int(b.X) + int(b.W)/2 - (len(b.Label)*10)/2
	textY := int(b.Y) + int(b.H)/2 - 6
	DrawText(screen, b.Label, float64(textX), float64(textY), 12, color.White, "default")
}
