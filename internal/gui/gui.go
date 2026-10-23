package gui

import (
	"errors"
	"strings"

	"net"

	"github.com/hajimehoshi/ebiten/v2"

	"tap/internal/gui/components"
)

var ErrQuit = errors.New("game terminated by user")

type Game struct {
	top          *components.Panel
	centerLeft   *components.Panel
	centerMid    *components.Panel
	centerRight  *components.Panel
	bottomLeft   *components.Panel
	bottomCenter *components.Panel
	bottomRight  *components.Panel
	name         string
	hp           int
	location     [2]int
	conn         net.Conn
	messagesChan chan string
	curr_cmd     string
	status       string
}

func NewGame(winWidth int, winHeight int) *Game {
	name := "Pâris"
	t, cl, cm, cr, bl, bc, br := InitPanels(winWidth, winHeight)
	conn, err := Connect("localhost:4242", name)
	messagesChan := make(chan string, 100)

	g := &Game{
		top:          t,
		centerLeft:   cl,
		centerMid:    cm,
		centerRight:  cr,
		bottomLeft:   bl,
		bottomCenter: bc,
		bottomRight:  br,
		name:         name,
		hp:           100,
		location:     [2]int{1, 1},
		conn:         conn,
		messagesChan: messagesChan,
		curr_cmd:     "",
	}

	if err == nil {
		go ListenServer(conn, messagesChan)

		if chatTeam := g.bottomCenter.GetChat(); chatTeam != nil {
			chatTeam.OnSubmit = func(text string) {
				SendMessage(conn, text)
			}
		}

		if chatAction := g.bottomRight.GetChat(); chatAction != nil {
			chatAction.OnSubmit = func(text string) {
				parts := strings.Split(text, " ")
				g.curr_cmd = strings.ToUpper(parts[0])

				SendMessage(conn, text)
			}
			g.sendActionButton(conn, chatAction)
		}

		g.sendServerDirection(conn)

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
		g.bottomCenter,
		g.bottomRight,
	)

	err := g.serverResponse()
	if err != nil {
		return ErrQuit
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
		g.bottomCenter,
		g.bottomRight,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1480, 1090
}

func Run() error {
	windowWidth := 1480
	windowHeight := 1090
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("The Answer Protocol")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	newGame := NewGame(windowWidth, windowHeight)
	return ebiten.RunGame(newGame)
}
