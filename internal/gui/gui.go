package gui

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"net"

	"github.com/hajimehoshi/ebiten/v2"

	"tap/internal/gui/components"
	"tap/internal/protocol"
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
	direction    int
	conn         net.Conn
	messagesChan chan string
	curr_cmd     string
	status       string
}

func NewGame(winWidth int, winHeight int, conn net.Conn, name string) *Game {
	t, cl, cm, cr, bl, bc, br := InitPanels(winWidth, winHeight)
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
		g.sendGroupButton(conn, chatAction)
	}

	g.sendServerDirection(conn)

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

func getUsername(serverScanner *bufio.Scanner, conn net.Conn) (string, error) {
	stdinScanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter your user name: ")
		if !stdinScanner.Scan() {
			return "", fmt.Errorf("failed to read username from stdin")
		}

		username := strings.TrimSpace(stdinScanner.Text())

		if username == "" {
			fmt.Println("Error: Username cannot be empty. Please try again.")
			continue
		}
		if strings.Contains(username, " ") {
			fmt.Println("Error: Username cannot contain spaces. Please try again.")
			continue
		}

		fmt.Fprint(conn, protocol.FormatCmd(protocol.CmdConnect, username))

		if !serverScanner.Scan() {
			return "", fmt.Errorf("connection lost with server")
		}

		response := serverScanner.Text()
		if response == "OK connected" {
			fmt.Println("Server:", response)
			return username, nil
		}

		if response == "ERR 201 NAME_IN_USE" {
			fmt.Println("Error: This username is already taken by another player. Please choose another one.")
		} else {
			fmt.Printf("Error (%s). Please try again.\n", response)
		}
	}
}

func Run(addr string) error {
	windowWidth := 1480
	windowHeight := 1090

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("Cannot connect to the TAP server on %s: %w", addr, err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to the TAP server.")
	serverScanner := bufio.NewScanner(conn)

	if serverScanner.Scan() {
		fmt.Println("Server:", serverScanner.Text())
	}

	username, err := getUsername(serverScanner, conn)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("The Answer Protocol")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	newGame := NewGame(windowWidth, windowHeight, conn, username)
	if err := ebiten.RunGame(newGame); err != nil {
		return fmt.Errorf("game error: %w", err)
	}

	return nil
}
