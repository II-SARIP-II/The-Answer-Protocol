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

func PrintTitle(screen *ebiten.Image, x float32, y float32, w float32, h float32, name string) {
    size := float32(20.0)
    startX := x + (w / 2.0) - (float32(len(name)) * size * 0.5 / 2.0)
    startY := y + size
    color := color.RGBA{R: 50, G: 50, B: 50, A: 255}
    fontType := "Default"
    DrawText(screen, name, float64(startX), float64(startY), float64(size), color, fontType)
}

func (p *Panel) Draw(screen *ebiten.Image) {
    if p == nil {
        return
    }
    DrawFilledRect(screen, p.X, p.Y, p.W, p.H, p.BgColor)
    DrawStrokeRect(screen, p.X, p.Y, p.W, p.H, 1, p.BorderColor)
    PrintTitle(screen, p.X, p.Y, p.W, p.H, p.Name)

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
