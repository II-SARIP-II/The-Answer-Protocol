package components

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Widget interface {
	Update()
	Draw(screen *ebiten.Image)
}

type Hole struct {
	Wall int
	Pos  float32
}

type Panel struct {
	Name        string
	DisplayName bool
	ToDraw      bool
	X, Y, W, H  float32
	BorderHoles []Hole
	Content     []Widget
	BgColor     color.Color
	BorderColor color.Color
	TextSize    float32
}

func (p *Panel) Update() {
	if p == nil {
		return
	}
	for _, widget := range p.Content {
		if widget != nil {
			widget.Update()
		}
	}
}

func PrintTitle(screen *ebiten.Image, x float32, y float32, w float32, h float32, name string, size float32) {
	startX := x + (w / 2.0) - (float32(len(name)) * size * 0.5 / 2.0)
	startY := y + size/2
	col := color.RGBA{R: 50, G: 50, B: 50, A: 255}
	fontType := "feather"
	DrawText(screen, name, float64(startX), float64(startY), float64(size), col, fontType)
}

func (p *Panel) Draw(screen *ebiten.Image) {
	if p == nil {
		return
	}
	if p.ToDraw {
		DrawFilledRect(screen, p.X, p.Y, p.W, p.H, p.BgColor)
		DrawStrokeRect(screen, p.X, p.Y, p.W, p.H, 1, p.BorderColor, p.BorderHoles)

		if p.DisplayName {
			PrintTitle(screen, p.X, p.Y, p.W, p.H, p.Name, p.TextSize)
		}
	}
	for _, widget := range p.Content {
		if widget != nil {
			widget.Draw(screen)
		}
	}
}

func UpdatePanels(panels ...*Panel) {
	for _, p := range panels {
		p.Update()
	}
}

func DrawPanels(screen *ebiten.Image, panels ...*Panel) {
	for _, p := range panels {
		p.Draw(screen)
	}
}
