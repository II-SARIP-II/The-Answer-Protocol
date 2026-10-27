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

func (p *Panel) GetStuff() *Stuff {
	if p == nil {
		return nil
	}
	for _, widget := range p.Content {
		if stuff, ok := widget.(*Stuff); ok {
			return stuff
		}
	}
	return nil
}

func (p *Panel) GetChat() *ChatWidget {
	if p == nil {
		return nil
	}
	for _, widget := range p.Content {
		if chat, ok := widget.(*ChatWidget); ok {
			return chat
		}
	}
	return nil
}

func (p *Panel) GetButtons() []*Button {
	if p == nil {
		return nil
	}
	var buttons []*Button
	for _, widget := range p.Content {
		if btn, ok := widget.(*Button); ok {
			buttons = append(buttons, btn)
		}
	}
	return buttons
}

func (p *Panel) GetSubPanel(name string) *Panel {
	if p == nil {
		return nil
	}
	for _, widget := range p.Content {
		if subPanel, ok := widget.(*Panel); ok && subPanel.Name == name {
			return subPanel
		}
	}
	return nil
}

type Panel struct {
	Name          string
	DisplayName   bool
	TextColor     color.Color
	ToDraw        bool
	X, Y, W, H    float32
	BorderHoles   []Hole
	Content       []Widget
	BgColor       color.Color
	BorderColor   color.Color
	TextSize      float32
	PlayerInvites []*PlayerInvite
	AcceptInvites []*PlayerInvite
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
	for _, invite := range p.PlayerInvites {
		if invite != nil && invite.Btn != nil {
			invite.Btn.Update()
		}
	}
	for _, invite := range p.AcceptInvites {
		if invite != nil && invite.Btn != nil {
			invite.Btn.Update()
		}
	}
}

func printTitle(screen *ebiten.Image, x float32, y float32, w float32, h float32, name string, size float32, textColor color.Color) {
	startX := x + (w / 2.0) - (float32(len(name)) * size * 0.5 / 2.0)
	startY := y + size/2
	fontType := "Title"
	drawText(screen, name, float64(startX), float64(startY), float64(size), textColor, fontType)
}

func (p *Panel) Draw(screen *ebiten.Image) {
	if p == nil || !p.ToDraw {
		return
	}

	drawFilledRect(screen, p.X, p.Y, p.W, p.H, p.BgColor)
	drawStrokeRect(screen, p.X, p.Y, p.W, p.H, 1, p.BorderColor, p.BorderHoles)

	if p.DisplayName {
		printTitle(screen, p.X, p.Y, p.W, p.H, p.Name, p.TextSize, p.TextColor)
	}

	for _, widget := range p.Content {
		if widget != nil {
			widget.Draw(screen)
		}
	}

	if len(p.PlayerInvites) > 0 {
		drawList(screen, p.PlayerInvites)
	}
	if len(p.AcceptInvites) > 0 {
		drawList(screen, p.AcceptInvites)
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
