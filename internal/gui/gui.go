package gui

import (
	"github.com/hajimehoshi/ebiten/v2"

	"tap/internal/gui/components"
	"tap/internal/world"
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
}

func NewGame(winWidth int, winHeight int, data world.GameData) *Game {
	t, cl, cm, cr, bl, bcl, bcr, br := InitPanels(winWidth, winHeight, data)
	g := &Game{
		top:               t,
		centerLeft:        cl,
		centerMid:         cm,
		centerRight:       cr,
		bottomLeft:        bl,
		bottomCenterLeft:  bcl,
		bottomCenterRight: bcr,
		bottomRight:       br,
	}

	return g
}

func (g *Game) Update() error {
	components.UpdatePanels(
		g.centerRight,
	)
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
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1480, 1090
}

func Run() error {
	windowWidth := 1480
	windowHeight := 1090
	data := world.ReadJson()
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("The Answer Protocol")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	newGame := NewGame(windowWidth, windowHeight, data)
	return ebiten.RunGame(newGame)
}
