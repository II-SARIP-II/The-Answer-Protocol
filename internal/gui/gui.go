package gui

import (
	"errors"
	"fmt"
	"strings"

	//	"fmt"
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
	status		string
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

		if buttons := cm.GetButtons(); buttons != nil {
			for _, button := range buttons {
				btn := button
				switch btn.Label {
				case "North":
					btn.OnClick = func() {
						if g.status == "AWAITING_MOVE_DIR" {
							msg := "MOVE North"
							SendMessage(conn, msg)
							g.status = ""
						}
					}
				case "East":
					btn.OnClick = func() {
						if g.status == "AWAITING_MOVE_DIR" {
							msg := "MOVE East"
							SendMessage(conn, msg)
							g.status = ""
						}
					}
				case "South":
					btn.OnClick = func() {
						if g.status == "AWAITING_MOVE_DIR" {
							msg := "MOVE South"
							SendMessage(conn, msg)
							g.status = ""
						}
					}
				case "West":
					btn.OnClick = func() {
						if g.status == "AWAITING_MOVE_DIR" {
							msg := "MOVE West"
							SendMessage(conn, msg)
							g.status = ""
						}
					}
				}
			}
		}

		if buttons := cr.GetButtons(); buttons != nil {
			for _, button := range buttons {
				btn := button

				switch btn.Label {
				case "MOVE":
					btn.OnClick = func() {
						g.status = "AWAITING_MOVE_DIR"
						g.curr_cmd = "MOVE"
						if chat != nil {
							chat.AddLine("Wait a direction")
						}
					}

				case "TAKE":
					btn.OnClick = func() {
						g.status = "AWAITING_TAKE_ITEM"
						g.curr_cmd = "TAKE"
					}

				case "DROP":
					btn.OnClick = func() {
						g.status = "AWAITING_DROP_ITEM"
						g.curr_cmd = "DROP"
					}

				case "TALK":
					btn.OnClick = func() {
						g.status = "IDLE"
						g.curr_cmd = "STATUS"
						if chat != nil {
							chat.AddLine(name + " > STATUS")
						}
						SendMessage(conn, "STATUS")
					}

				case "STATUS":
					btn.OnClick = func() {
						g.status = "IDLE"
						g.curr_cmd = "STATUS"
						if chat != nil {
							chat.AddLine(name + " > STATUS")
						}
						SendMessage(conn, "STATUS")
					}

				case "QUEST":
					btn.OnClick = func() {
						g.status = "IDLE"
						g.curr_cmd = "STATUS"
						if chat != nil {
							chat.AddLine(name + " > STATUS")
						}
						SendMessage(conn, "STATUS")
					}

				case "QUESTS":
					btn.OnClick = func() {
						g.status = "IDLE"
						g.curr_cmd = "STATUS"
						if chat != nil {
							chat.AddLine(name + " > STATUS")
						}
						SendMessage(conn, "STATUS")
					}

				case "WHO":
					btn.OnClick = func() {
						g.status = "IDLE"
						g.curr_cmd = "STATUS"
						if chat != nil {
							chat.AddLine(name + " > STATUS")
						}
						SendMessage(conn, "STATUS")
					}

				case "GROUP":
					btn.OnClick = func() {
						g.status = "IDLE"
						g.curr_cmd = "STATUS"
						if chat != nil {
							chat.AddLine(name + " > STATUS")
						}
						SendMessage(conn, "STATUS")
					}

				default:
					btn.OnClick = func() {
						g.status = "IDLE"
						g.curr_cmd = btn.Label
						if chat != nil {
							chat.AddLine(name + " > " + btn.Label)
						}
						SendMessage(conn, btn.Label)
					}
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
				fmt.Print(msg)
				jsonStart := strings.Index(msg, "{")
        		response := msg
        		if jsonStart != -1 {
        		    response = msg[jsonStart:]
        		}

				stuff := g.bottomLeft.GetStuff()
				switch g.curr_cmd {
				case "LOOK":
    				lookResp, err := ParseLookResponse(msg)
    				if err == nil {
    				    room := lookResp.ToWorldRoom()

    				    g.centerMid.AddRoom(
    				        room.Name,
    				        room.Exits,
    				        room.Items,
    				        func(itemName string) {
    				            if g.status == "AWAITING_TAKE_ITEM" {
    				                cmd := "TAKE " + itemName
    				                g.curr_cmd = "TAKE"
    				                if chat != nil {
    				                    chat.AddLine(g.name + " > " + cmd)
    				                }
    				                SendMessage(g.conn, cmd)
    				                g.status = ""
    				            }
    				        },
    				    )
    				    chat.AddLine(g.curr_cmd + " > " + room.Name)
    				} else {
    				    chat.AddLine("Server > " + msg)
    				}
                case "MOVE":
                    chat.AddLine(g.curr_cmd + " > " + response)
                case "TAKE":
    				if strings.Contains(msg, "taken=") {
    				    parts := strings.Split(msg, "taken=")
    				    if len(parts) > 1 {
    				        itemName := strings.TrimSpace(parts[1])

    				        if stuff != nil {
    				            stuff.AddItem(world.Item{Name: itemName})
    				        }
    				        chat.AddLine(g.curr_cmd + " > " + itemName)

    				        g.curr_cmd = "LOOK"
    				        SendMessage(g.conn, "LOOK")
    				    }
    				} else {
    				    chat.AddLine("Server > " + msg)
    				}
                case "DROP":
                    item := world.ParseItem(response)
                    if stuff != nil {
                        stuff.RemoveItem(item)
                    }
                    chat.AddLine(g.curr_cmd + " > " + fmt.Sprintf("%+v", item))
//				case "TALK":
//    			    world.ParseRoom(msg)
//				case "ATTACK":
//    			    // attack function
//				case "STATUS":
//    			    // status function
//				case "QUEST":
//    			    // quest function
//				case "QUESTS":
//    			    // quests function
				case "WHO":
    			    player := ParsePlayer(response)
					chat.AddLine(g.curr_cmd + fmt.Sprintf("%+v", player))
//				case "GROUP":
//    			    world.ParseRoom(msg)
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
