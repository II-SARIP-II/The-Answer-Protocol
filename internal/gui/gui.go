package gui

import (
	"errors"
	"fmt"
	"net"

	"github.com/hajimehoshi/ebiten/v2"

	"tap/internal/gui/components"
	"tap/internal/world"
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
	conn         net.Conn
	messagesChan chan string
	curr_cmd     string
}

func NewGame(winWidth int, winHeight int, data world.GameData) *Game {
	name := "Pâris"
	t, cl, cm, cr, bl, bc, br := InitPanels(winWidth, winHeight, data)
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
		conn:         conn,
		messagesChan: messagesChan,
		curr_cmd:     "",
	}

	if err == nil {
		go ListenServer(conn, messagesChan)

		chat := br.GetChat()
		if chat != nil {
			chat.OnSubmit = func(text string) {
				g.curr_cmd = text
				SendMessage(conn, text)
			}
		}

		if buttons := cr.GetButtons(); buttons != nil {
			for _, button := range buttons {
				btnLabel := button.Label

				button.OnClick = func() {
					g.curr_cmd = btnLabel
					if chat != nil {
						chat.AddLine(name + " > " + btnLabel)
					}
					SendMessage(conn, btnLabel)
				}
			}
		}
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

	for {
		select {
		case msg := <-g.messagesChan:
			if chat := g.bottomRight.GetChat(); chat != nil {
				chat.AddLine("Server > " + msg)
				switch g.curr_cmd {
    			case "LOOK":
    			    room := world.ParseRoom(msg)
					chat.AddLine("Room > " + fmt.Sprintf("%+v", room))
    			case "MOVE":
    			    world.ParseRoom(msg)
    			case "TAKE":
    			    world.ParseRoom(msg)
				case "DROP":
    			    world.ParseRoom(msg)
				case "TALK":
    			    world.ParseRoom(msg)
				case "ATTACK":
    			    world.ParseRoom(msg)
				case "STATUS":
    			    world.ParseRoom(msg)
				case "QUEST":
    			    world.ParseRoom(msg)
				case "QUESTS":
    			    world.ParseRoom(msg)
				case "WHO":
    			    world.ParseRoom(msg)
				case "GROUP":
    			    world.ParseRoom(msg)
				case "QUIT":
					fmt.Printf("QUIT")
					return ErrQuit
				}
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
	data, err := world.ReadJson("data/world.json")
	if err != nil {
		return err
	}
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("The Answer Protocol")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	newGame := NewGame(windowWidth, windowHeight, *data)
	return ebiten.RunGame(newGame)
}
