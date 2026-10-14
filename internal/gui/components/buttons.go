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
	TextColor    color.RGBA

	OnClick func()
}

func NewButton(label string, x, y, w, h float32) *Button {
	return &Button{
		X:            x,
		Y:            y,
		W:            w,
		H:            h,
		Label:        label,
		NormalColor:  color.RGBA{R: 12, G: 24, B: 36, A: 220},
		HoverColor:   color.RGBA{R: 0, G: 180, B: 216, A: 180},
		PressedColor: color.RGBA{R: 0, G: 255, B: 170, A: 230},
		BorderColor:  color.RGBA{R: 0, G: 220, B: 255, A: 255},
		TextColor:    color.RGBA{R: 0, G: 255, B: 200, A: 255},
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
	textColor := b.TextColor
	borderWidth := float32(1)

	if b.IsPressed {
		currentColor = b.PressedColor
		textColor = color.RGBA{R: 10, G: 15, B: 20, A: 255} // Texte sombre sur fond brillant
		borderWidth = 2
	} else if b.IsHovered {
		currentColor = b.HoverColor
		textColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		borderWidth = 2
	}

	// Fond translucide
	DrawFilledRect(screen, b.X, b.Y, b.W, b.H, currentColor)

	// Contour style HUD
	var hole []Hole = []Hole{}
	DrawStrokeRect(screen, b.X, b.Y, b.W, b.H, borderWidth, b.BorderColor, hole)

	// Petit détail coin HUD si survolé
	if b.IsHovered {
		DrawFilledRect(screen, b.X, b.Y, 4, 4, color.RGBA{R: 255, G: 255, B: 255, A: 255})
		DrawFilledRect(screen, b.X+b.W-4, b.Y+b.H-4, 4, 4, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	}

	textX := int(b.X) + int(b.W)/2 - (len(b.Label)*10)/2
	textY := int(b.Y) + int(b.H)/2 - 6
	DrawText(screen, b.Label, float64(textX), float64(textY), 20, textColor, "Title")
}
