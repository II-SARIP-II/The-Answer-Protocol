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

// Server manages game state, connected players, and TCP listening.
type Server struct {
	mu                sync.Mutex
	listener          net.Listener
	players           map[string]*Player
	port              string
	world             *world.GameData
	spawnsQuestsItems map[string]string
}

// NewServer initializes and returns a new Server instance with world data.
func NewServer(worldPath string) (*Server, error) {
	w, err := world.ReadJson(worldPath)
	if err != nil {
		return nil, err
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

		go s.handleClientSession(conn)
	}
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

		player.HandleCommand(command)

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
	for clientScanner.Scan() {
		p, err := s.authenticate(conn, clientScanner.Text())
		if err == nil {
			player = p
			break
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

	room, exists := s.world.Rooms[player.CurrentRoom]
	if exists && len(player.Inventory) > 0 {
		room.Items = append(room.Items, player.Inventory...)
		s.world.Rooms[player.CurrentRoom] = room
		player.Inventory = nil
		slog.Info("Player items dropped on disconnect", "player", player.Username, "room", player.CurrentRoom)
	}

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
}
