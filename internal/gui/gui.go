package gui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"tap/internal/gui/components"
)

type Game struct {
    top               *components.Panel
    centerLeft        *components.Panel
    centerMid         *components.Panel
    centerRight       *components.Panel
    bottomLeft        *components.Panel
    bottomCenterLeft  *components.Panel
    bottomCenterRight *components.Panel
    bottomRight       *components.Panel
    buttons           []*components.Button
    lastMessage       string
}

func NewGame(winWidth int, winHeight int) *Game {
    t, cl, cm, cr, bl, bcl, bcr, br := InitPanels(winWidth, winHeight)
    g := &Game{
        top:               t,
        centerLeft:        cl,
        centerMid:         cm,
        centerRight:       cr,
        bottomLeft:        bl,
        bottomCenterLeft:  bcl,
        bottomCenterRight: bcr,
        bottomRight:       br,
        lastMessage:       "Game Started",
    }

    return g
}

func (g *Game) Update() error {
	for _, btn := range g.buttons {
		btn.Update()
	}
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
    components.DrawPanels(
        screen,
        g.top,
        g.centerLeft,
        g.centerMid,
        g.centerRight,
        g.bottomLeft,
        g.bottomCenterLeft,
        g.bottomCenterRight,
        g.bottomRight,
    )

    for _, btn := range g.buttons {
        btn.Draw(screen)
    }

    components.DrawText(screen, g.lastMessage, 10, 10, 24, color.White, "title")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1480, 1090
}

func Run() error {
	windowWidth := 1480
	windowHeight := 1090
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Final Test")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	newGame := NewGame(windowWidth, windowHeight)
	return ebiten.RunGame(newGame)
}
