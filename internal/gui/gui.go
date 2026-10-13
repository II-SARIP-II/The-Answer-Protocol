package gui

import (
	"net"

	"github.com/hajimehoshi/ebiten/v2"

	"tap/internal/gui/components"
	"tap/internal/world"
)

type Game struct {
	top					*components.Panel
	centerLeft			*components.Panel
	centerMid			*components.Panel
	centerRight			*components.Panel
	bottomLeft			*components.Panel
	bottomCenterLeft	*components.Panel
	bottomCenterRight	*components.Panel
	bottomRight			*components.Panel
	name 				string
	conn				net.Conn
	messagesChan		chan string
}

func NewGame(winWidth int, winHeight int, data world.GameData) *Game {
	name := "Pâris"
	t, cl, cm, cr, bl, bcl, bcr, br := InitPanels(winWidth, winHeight, data)
	conn, err := Connect("localhost:4242", "Edouard")
	messagesChan := make(chan string, 100)
	go ListenServer(conn, messagesChan)
	if err == nil {
		chat := br.GetChat()
	    if chat != nil {
	        chat.OnSubmit = func(text string) {
	            SendMessage(conn, text)
	        }
	    }
		if btn := cr.GetButton(); btn != nil {
	        btn.OnClick = func(text string) {
				chat.AddLine(name + " > " + text)
	            SendMessage(conn, text)
	        }
	    }
	}

	g := &Game{
        top:               t,
        centerLeft:        cl,
        centerMid:         cm,
        centerRight:       cr,
        bottomLeft:        bl,
        bottomCenterLeft:  bcl,
        bottomCenterRight: bcr,
        bottomRight:       br,
		name: 				name,
        conn:              conn,
        messagesChan:      messagesChan,
    }

	return g
}

func (g *Game) Update() error {
    components.UpdatePanels(
        g.top,
        g.centerLeft,
        g.centerMid,
        g.centerRight,
        g.bottomLeft,
        g.bottomCenterLeft,
        g.bottomCenterRight,
        g.bottomRight,
    )

    for {
        select {
        case msg := <-g.messagesChan:
            if chat := g.bottomRight.GetChat(); chat != nil {
                chat.AddLine("Server > " + msg)
            }
        default:
            goto Done
        }
    }
Done:

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
