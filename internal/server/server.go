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
	"time"
)

var (
	mu      sync.Mutex
	players = make(map[string]*Player)
)

func Start(port string) error {
	// la fonction Listen ouvre un socket sur le port indiqué et renvoie un listener ln ainsi qu'éventuellement une erreur
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

		go handleClientSession(conn)
	}
}

func handleClientSession(conn net.Conn) {
	defer conn.Close()

	state := protocol.StateConnected
	slog.Info("Client connected", "addr", conn.RemoteAddr().String(), "state", state)

	fmt.Fprint(conn, protocol.MsgHello)

	clientScanner := bufio.NewScanner(conn)
	if !clientScanner.Scan() {
		return
	}

	player, err := authenticate(conn, clientScanner.Text())
	if err != nil {
		return
	}

	player.State = protocol.StateAuthenticated
	fmt.Fprint(conn, protocol.MsgConnected)
	slog.Info("Client authenticated", "username", player.Username, "addr", conn.RemoteAddr().String(), "state", player.State)

	defer func() {
		mu.Lock()
		player.State = protocol.StateTerminated
		delete(players, player.Username)
		slog.Info("Client disconnected", "username", player.Username, "addr", conn.RemoteAddr().String(), "state", player.State)
		mu.Unlock()
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

func authenticate(conn net.Conn, command string) (*Player, error) {
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

	mu.Lock()
	if _, exists := players[username]; exists {
		mu.Unlock()
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeNameInUse, protocol.MsgNameInUse))
		slog.Warn("Authentication failed", "addr", conn.RemoteAddr().String(), "code", protocol.CodeNameInUse, "msg", protocol.MsgNameInUse)

		return nil, errors.New("name in use")
	}
	player := &Player{
		Username: username,
		Conn:     conn,
		State:    protocol.StateAuthenticated,
	}
	players[username] = player
	mu.Unlock()
	return player, nil

}
