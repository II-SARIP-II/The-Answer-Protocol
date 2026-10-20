package server

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"
	"tap/internal/protocol"
	"tap/internal/world"
	"time"
)

const NPCRespawnDelay = 30 * time.Second

type Server struct {
	mu            sync.Mutex
	players       (map[string]*Player)
	port          string
	world         *world.GameData
	npcSpawnRooms map[string]string
	initialNPCs   map[string]world.NPC
}

func NewServer(worldPath string) (*Server, error) {
	w, err := world.ReadJson(worldPath)
	if err != nil {
		return nil, err
	}

	npcSpawnRooms := make(map[string]string)
	initialNPCs := make(map[string]world.NPC)

	for roomID, room := range w.Rooms {
		for _, npcID := range room.NPCs {
			npcSpawnRooms[npcID] = roomID
		}
	}
	for npcID, npc := range w.NPCs {
		initialNPCs[npcID] = npc
	}

	return &Server{
		players:       make(map[string]*Player),
		world:         w,
		npcSpawnRooms: npcSpawnRooms,
		initialNPCs:   initialNPCs,
	}, nil
}

func (s *Server) scheduleNPCRespawn(enemyID string) {
	spawnRoom, okRoom := s.npcSpawnRooms[enemyID]
	initialData, okData := s.initialNPCs[enemyID]
	if !okRoom || !okData {
		return
	}

	time.AfterFunc(NPCRespawnDelay, func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		s.world.NPCs[enemyID] = initialData

		room := s.world.Rooms[spawnRoom]
		alreadyPresent := false
		for _, id := range room.NPCs {
			if id == enemyID {
				alreadyPresent = true
				break
			}
		}
		if !alreadyPresent {
			room.NPCs = append(room.NPCs, enemyID)
			s.world.Rooms[spawnRoom] = room
		}

		slog.Info("NPC respawned", "npc", enemyID, "room", spawnRoom)

		for _, player := range s.players {
			if strings.EqualFold(player.CurrentRoom, spawnRoom) {
				player.SendEvt("ROOM CHAT", fmt.Sprintf("Server %s has reappeared!", initialData.Name))
			}
		}
	})
}

func (s *Server) Start(port string) error {
	// la fonction Listen ouvre un socket sur le port indiqué et renvoie un listener ln ainsi qu'éventuellement une erreur
	s.port = port
	ln, err := net.Listen("tcp", port)
	if err != nil {
		slog.Error("Cannot open port", "port", port, "err", err)
		return err
	}
	defer ln.Close()

	slog.Info("TAP server started", "port", port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			// si le listener a été fermé proprement
			if errors.Is(err, net.ErrClosed) {
				slog.Info("Listener closed, server shutting down cleanly")
				return nil
			}

			// gestion d'erreur réseau temporaire
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

func (s *Server) handleClientSession(conn net.Conn) {
	defer conn.Close()

	state := protocol.StateConnected
	slog.Info("Client connected", "addr", conn.RemoteAddr().String(), "state", state)

	fmt.Fprint(conn, protocol.MsgHello)

	clientScanner := bufio.NewScanner(conn)
	if !clientScanner.Scan() {
		return
	}

	player, err := s.authenticate(conn, clientScanner.Text())
	if err != nil {
		return
	}

	player.State = protocol.StateAuthenticated
	fmt.Fprint(conn, protocol.MsgConnected)
	slog.Info("Client authenticated", "username", player.Username, "addr", conn.RemoteAddr().String(), "state", player.State)

	s.mu.Lock()
	for _, otherPlayer := range s.players {
		if otherPlayer.Username != player.Username {
			otherPlayer.SendEvt(protocol.EvtGameEnter, player.Username)
		}
	}
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		player.State = protocol.StateTerminated
		player.leaveGroupGetModif()
		delete(s.players, player.Username)
		for _, otherPlayer := range s.players {
			otherPlayer.SendEvt(protocol.EvtGameLeave, player.Username)
		}
		slog.Info("Client disconnected", "username", player.Username, "addr", conn.RemoteAddr().String(), "state", player.State)
		s.mu.Unlock()
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
	if error := clientScanner.Err(); error != nil {
		slog.Warn("Client connection error", "username", player.Username, "error", error)
	}
}

func (s *Server) authenticate(conn net.Conn, command string) (*Player, error) {
	parts := strings.SplitN(command, " ", 2)
	if len(parts) != 2 || strings.ToUpper(parts[0]) != protocol.CmdConnect {
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeBadRequest, protocol.MsgUnknownCommand))
		slog.Warn("Authentication failed", "addr", conn.RemoteAddr().String(), "code", protocol.CodeBadRequest, "msg", protocol.MsgUnknownCommand)
		return nil, errors.New("Authentication failed: invalid command")
	}

	username := strings.TrimSpace(parts[1])
	if username == "" {
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeBadRequest, protocol.MsgArgsError))
		slog.Warn("Authentication failed: empty username", "addr", conn.RemoteAddr().String(), "code", protocol.CodeBadRequest, "msg", protocol.MsgArgsError)
		return nil, errors.New("username is empty.")
	}

	if strings.Contains(username, " ") {
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeBadRequest, protocol.MsgArgsError))
		slog.Warn("Authentication failed: space in username", "addr", conn.RemoteAddr().String(), "code", protocol.CodeBadRequest, "msg", protocol.MsgArgsError)
		return nil, errors.New("space in username.")
	}

	s.mu.Lock()
	if _, exists := s.players[username]; exists {
		s.mu.Unlock()
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeNameInUse, protocol.MsgNameInUse))
		slog.Warn("Authentication failed", "addr", conn.RemoteAddr().String(), "code", protocol.CodeNameInUse, "msg", protocol.MsgNameInUse)

		return nil, errors.New("name in use")
	}

	player := NewPlayer(username, conn)
	player.Server = s
	s.players[username] = player

	s.mu.Unlock()
	return player, nil

}
