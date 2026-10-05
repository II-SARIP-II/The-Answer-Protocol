package components

import (
    "image/color"

    "github.com/hajimehoshi/ebiten/v2"
)

type Panel struct {
    Name        string
    X, Y, W, H  float32
    Buttons     []*Button
    BgColor     color.Color
    BorderColor color.Color
}

func (p *Panel) Update() {
    if p == nil {
        return
    }
    for _, btn := range p.Buttons {
        btn.Update()
    }
}

func (p *Panel) Draw(screen *ebiten.Image) {
    if p == nil {
        return
    }
    DrawFilledRect(screen, p.X, p.Y, p.W, p.H, p.BgColor)
    DrawStrokeRect(screen, p.X, p.Y, p.W, p.H, 1, p.BorderColor)

    for _, btn := range p.Buttons {
        btn.Draw(screen)
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
