package server

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"tap/internal/protocol"
	"tap/internal/world"
	"time"
)

type connNumber struct {
	timeElapsed time.Time
	count       int
}

// Server manages game state, connected players, and TCP listening.
type Server struct {
	mu                sync.Mutex
	listener          net.Listener
	players           map[string]*Player
	port              string
	world             *world.GameData
	spawnsQuestsItems map[string]string
	connTracker       map[string]connNumber
}

// NewServer initializes and returns a new Server instance with world data.
func NewServer(worldPath string) (*Server, error) {
	w, err := world.ReadJson(worldPath)
	if err != nil {
		return nil, err
	}

	for id, npc := range w.NPCs {
		npc.InitialRole = npc.Role
		w.NPCs[id] = npc
	}

	spawnsQuestsItems := make(map[string]string)
	for roomID, room := range w.Rooms {
		for _, itemID := range room.Items {
			spawnsQuestsItems[itemID] = roomID
		}
	}

	return &Server{
		players:           make(map[string]*Player),
		world:             w,
		spawnsQuestsItems: spawnsQuestsItems,
		connTracker:       make(map[string]connNumber),
	}, nil
}

// Start starts listening for client connections on the given port.
func (s *Server) Start(port string) error {
	s.port = port
	ln, err := net.Listen("tcp", port)
	if err != nil {
		slog.Error("Cannot open port", "port", port, "err", err)
		return err
	}
	defer ln.Close()
	s.listener = ln
	slog.Info("TAP server started", "port", port)

	// Intercept shutdown signals for graceful shutdown.
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stopChan)
	go func() {
		<-stopChan
		s.Stop()
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				slog.Info("Server stopped cleanly")
				return nil
			}

			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				slog.Warn("Temporary network error on Accept(), pausing 10ms...", "err", err)
				time.Sleep(10 * time.Millisecond)
				continue
			}

			slog.Error("Fatal error on Accept()", "err", err)
			return err
		}

		s.trackConnection(conn.RemoteAddr().String())
		go s.handleClientSession(conn)
	}
}

// trackConnection logs a warning if the same IP opens more than 5 connections within 2 seconds.
func (s *Server) trackConnection(remoteAddr string) {
	host, _, _ := net.SplitHostPort(remoteAddr)

	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	connNumber := s.connTracker[host]

	if now.Sub(connNumber.timeElapsed) < 2*time.Second {
		connNumber.count++
		if connNumber.count > 5 {
			slog.Warn("Potential abuse detected: rapid connections",
				"ip", host,
				"connections_in_window", connNumber.count,
			)
		}
	} else {
		connNumber.timeElapsed = now
		connNumber.count = 1
	}
	s.connTracker[host] = connNumber
}

func (s *Server) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, player := range s.players {
		player.Conn.Close()
	}

	if s.listener != nil {
		s.listener.Close()
	}
}

func (s *Server) handleClientSession(conn net.Conn) {
	defer conn.Close()

	state := protocol.StateConnected
	slog.Info("Client connected", "addr", conn.RemoteAddr().String(), "state", state)

	fmt.Fprint(conn, protocol.MsgHello)

	clientScanner := bufio.NewScanner(conn)
	clientScanner.Buffer(make([]byte, 1024), 1024)

	player, err := s.newPlayer(conn, clientScanner)
	if err != nil {
		return
	}

	defer func() {
		s.handleDisconnect(player)
	}()

	for clientScanner.Scan() {
		command := strings.TrimSpace(clientScanner.Text())
		if command == "" {
			continue
		}

		player.handleCommand(command)

		if player.State == protocol.StateTerminated {
			break
		}
	}
	if err := clientScanner.Err(); err != nil && !errors.Is(err, net.ErrClosed) {
		slog.Warn("Client connection error", "username", player.Username, "error", err)
	}
}

func (s *Server) newPlayer(conn net.Conn, clientScanner *bufio.Scanner) (*Player, error) {
	var player *Player
	failedAttempts := 0
	for clientScanner.Scan() {
		p, err := s.authenticate(conn, clientScanner.Text())
		if err == nil {
			player = p
			break
		}
		failedAttempts++
		if failedAttempts >= 3 {
			slog.Warn("Potential abuse detected: repeated authentication failures",
				"addr", conn.RemoteAddr().String(),
				"failed_attempts", failedAttempts,
			)
			fmt.Fprint(conn, protocol.FormatErr(protocol.CodeConnectionFailed, protocol.MsgConnectionFailed))
			return nil, errors.New("too many failed authentication attempts")
		}
	}
	if player == nil {
		return nil, errors.New("client disconnected before authenticating")
	}

	player.State = protocol.StateAuthenticated
	fmt.Fprint(conn, protocol.MsgConnected)
	slog.Info("Client authenticated", "username", player.Username, "addr", conn.RemoteAddr().String(), "state", player.State)

	s.mu.Lock()
	for _, otherPlayer := range s.players {
		if otherPlayer.Username != player.Username && otherPlayer.CurrentRoom == player.CurrentRoom {
			otherPlayer.SendEvt(protocol.EvtRoomEnter, player.Username)
		}
	}

	data := fmt.Sprintf("players=%d", len(s.players))
	for _, otherPlayer := range s.players {
		otherPlayer.SendEvt(protocol.EvtPlayers, data)
	}

	s.mu.Unlock()
	return player, nil
}

func (s *Server) authenticate(conn net.Conn, command string) (*Player, error) {
	parts := strings.SplitN(strings.TrimSpace(command), " ", 2)
	if len(parts) == 0 || strings.ToUpper(parts[0]) != protocol.CmdConnect {
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeBadRequest, protocol.MsgUnknownCommand))
		slog.Warn("Authentication failed", "addr", conn.RemoteAddr().String(), "code", protocol.CodeBadRequest, "protocol_msg", protocol.MsgUnknownCommand)
		return nil, errors.New("Authentication failed: invalid command")
	}

	if len(parts) == 1 || strings.TrimSpace(parts[1]) == "" {
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeBadRequest, protocol.MsgArgsError))
		slog.Warn("Authentication failed: empty username", "addr", conn.RemoteAddr().String(), "code", protocol.CodeBadRequest)
		return nil, errors.New("username is empty")
	}

	username := strings.TrimSpace(parts[1])

	if strings.Contains(username, " ") {
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeBadRequest, protocol.MsgArgsError))
		slog.Warn("Authentication failed: space in username", "addr", conn.RemoteAddr().String(), "code", protocol.CodeBadRequest, "protocol_msg", protocol.MsgArgsError)
		return nil, errors.New("space in username.")
	}

	s.mu.Lock()
	if _, exists := s.players[username]; exists {
		s.mu.Unlock()
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeNameInUse, protocol.MsgNameInUse))
		slog.Warn("Authentication failed", "addr", conn.RemoteAddr().String(), "code", protocol.CodeNameInUse, "protocol_msg", protocol.MsgNameInUse)

		return nil, errors.New("name in use")
	}

	player := NewPlayer(username, conn)
	player.Server = s
	s.players[username] = player

	s.mu.Unlock()
	return player, nil

}

func (s *Server) handleDisconnect(player *Player) {
	s.mu.Lock()
	player.State = protocol.StateTerminated
	player.leaveGroupGetModif()

	for _, itemID := range player.Inventory {
		if origRoomID, exists := s.spawnsQuestsItems[itemID]; exists {
			origRoom := s.world.Rooms[origRoomID]
			origRoom.Items = append(origRoom.Items, itemID)
			s.world.Rooms[origRoomID] = origRoom
			s.resetPassiveEnemy(origRoomID, itemID)
			slog.Info("World item restored to origin on disconnect", "item", itemID, "room", origRoomID)
		}
	}

	player.Inventory = nil

	delete(s.players, player.Username)

	for _, otherPlayer := range s.players {
		if otherPlayer.CurrentRoom == player.CurrentRoom {
			otherPlayer.SendEvt(protocol.EvtRoomLeave, player.Username)
		}
	}

	data := fmt.Sprintf("players=%d", len(s.players))
	for _, otherPlayer := range s.players {
		otherPlayer.SendEvt(protocol.EvtPlayers, data)
	}

	slog.Info("Client disconnected", "username", player.Username, "addr", player.Conn.RemoteAddr().String(), "state", player.State)
	s.mu.Unlock()
	close(player.sendChan)
}
