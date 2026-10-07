package server

import (
	"fmt"
	"net"
	"tap/internal/protocol"
)

type Player struct {
	Username string
	Conn net.Conn
	State protocol.State
	// TODO RoomID, HP, Inventory...
}

func NewPlayer(username string, conn net.Conn) *Player {
	return &Player{
		Username: 	username,
		Conn: 		conn,
		State:		protocol.StateConnected,
	}
}

func (p *Player) SendOK(data string) {
	fmt.Fprintf(p.Conn, protocol.FormatOK(data))
}

func (p *Player) SendErr(code int, msg string) {
	fmt.Fprintf(p.Conn, protocol.FormatErr(code, msg))
}

func (p *Player) SendEvt(eventType string, eventData string) {
	fmt.Fprintf(p.Conn, protocol.FormatEvt(eventType, eventData))
}
