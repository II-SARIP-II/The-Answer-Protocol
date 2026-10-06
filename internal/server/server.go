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

// mot clé var pour déclarer un bloc de déclaration de variables globales
var (
	mu      sync.Mutex
	players = make(map[string]net.Conn)
)

// fonction pour démarrer le serveur
func Start(port string) error {
	// la fonction Listen ouvre un socket sur le port indiqué et renvoie un listener ln ainsi qu'éventuellement une erreur
	ln, err := net.Listen("tcp", port)
	if err != nil {
		// logs structurés avec slog
		slog.Error("Cannot open port", "port", port, "err", err)
		return err
	}
	defer ln.Close()

	slog.Info("TAP server started", "port", port)

	for {
		// on attend une connexion client
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

		// on lance une goroutine pour chaque client
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	slog.Info("Client connected", "addr", conn.RemoteAddr().String())

	// envoi du message d'accueil RFC
	fmt.Fprint(conn, protocol.MsgHello)

	clientScanner := bufio.NewScanner(conn)
	if !clientScanner.Scan() {
		return
	}

	line := clientScanner.Text()
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 || strings.ToUpper(parts[0]) != protocol.CmdConnect {
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeBadRequest, protocol.MsgBadRequest))
		return
	}

	username := strings.TrimSpace(parts[1])
	if username == "" {
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeBadRequest, protocol.MsgBadRequest))
		return
	}

	// vérification doublon sous verrou
	mu.Lock()
	if _, exists := players[username]; exists {
		mu.Unlock()
		fmt.Fprint(conn, protocol.FormatErr(protocol.CodeNameInUse, protocol.MsgNameInUse))
		return
	}
	players[username] = conn
	mu.Unlock()

	// confirmation de connexion
	fmt.Fprint(conn, protocol.MsgConnected)
	slog.Info("Player authenticated", "username", username)

	// nettoyage automatique du joueur lors de la déconnexion
	defer func() {
		mu.Lock()
		delete(players, username)
		slog.Info("Player disconnected", "username", username)
		mu.Unlock()
	}()

	// boucle de réception des commandes de jeu
	for clientScanner.Scan() {
		command := strings.TrimSpace(clientScanner.Text())
		if command == "" {
			continue
		}

		slog.Info("Command received", "username", username, "cmd", command)

		if strings.ToUpper(command) == protocol.CmdQuit {
			fmt.Fprint(conn, protocol.MsgBye)
			break
		}

		fmt.Fprint(conn, protocol.FormatOK(fmt.Sprintf("received (TODO): %s", command)))
	}

	if err := clientScanner.Err(); err != nil {
		slog.Warn("Client connection error", "username", username, "error", err)
	} else {
		slog.Info("Client cleanly disconnected", "username", username)
	}
}
