package server

import (
	"fmt"
	"log/slog"
	"net"
	"tap/internal/protocol"
	"time"
)

type Player struct {
	Server           *Server
	Username         string
	Conn             net.Conn
	State            protocol.State
	CurrentRoom      string
	Inventory        []string
	MaxHP            int
	HP               int
	EquippedWeapon   string
	DefendCount      int
	Status           string
	GroupID          string
	GroupInvitations []string
	Quests           map[string]Quest
	lastCmdWindow    time.Time
	cmdCountInSec    int
	sendChan         chan string
}

type Quest struct {
	ID             string
	Status         string
	Steps          int
	CompletedSteps int
}

func NewPlayer(username string, conn net.Conn) *Player {
	player := &Player{
		Username:         username,
		Conn:             conn,
		State:            protocol.StateAuthenticated,
		CurrentRoom:      "hotel_entrance",
		Inventory:        make([]string, 0),
		MaxHP:            100,
		HP:               100,
		EquippedWeapon:   "",
		DefendCount:      0,
		GroupID:          "",
		GroupInvitations: make([]string, 0),
		Status:           "healthy",
		Quests:           make(map[string]Quest),
		lastCmdWindow:    time.Now(),
		cmdCountInSec:    0,
		sendChan:         make(chan string, 128),
	}

	go player.writeMessages()
	return player
}

func (p *Player) writeMessages() {
	for msg := range p.sendChan {
		if _, err := fmt.Fprint(p.Conn, msg); err != nil {
			slog.Warn("TCP write error", "player", p.Username, "error", err)
			return
		}
	}
}

func (p *Player) SendOK(data string) {
	p.sendChan <- protocol.FormatOK(data)
	slog.Info("Response sent", "player", p.Username, "response", "OK", "data", data)
}

func (p *Player) SendErr(code int, msg string) {
	p.sendChan <- protocol.FormatErr(code, msg)
	slog.Warn("Error response sent", "player", p.Username, "code", code, "message", msg)
}

func (p *Player) SendEvt(eventType string, eventData string) {
	p.sendChan <- protocol.FormatEvt(eventType, eventData)
	slog.Info("Event pushed", "player", p.Username, "event type", eventType, "event data", eventData)
}
