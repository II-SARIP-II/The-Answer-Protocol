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
}

type Quest struct {
	ID             string
	Status         string
	Steps          int
	CompletedSteps int
}

func NewPlayer(username string, conn net.Conn) *Player {
	return &Player{
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
	}
}

func (p *Player) SendOK(data string) {
	fmt.Fprint(p.Conn, protocol.FormatOK(data))
	slog.Info("Response sent", "player", p.Username, "response", "OK", "data", data)
}

func (p *Player) SendErr(code int, msg string) {
	fmt.Fprint(p.Conn, protocol.FormatErr(code, msg))
	slog.Warn("Error response sent", "player", p.Username, "code", code, "message", msg)
}

func (p *Player) SendEvt(eventType string, eventData string) {
	fmt.Fprint(p.Conn, protocol.FormatEvt(eventType, eventData))
	slog.Info("Event pushed", "player", p.Username, "event type", eventType, "event data", eventData)
}
