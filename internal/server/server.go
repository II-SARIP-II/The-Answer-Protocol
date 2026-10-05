package server

import (
	"errors"
	"fmt"
	"log/slog"
	"net"
	"time"
	"sync"
	"bufio"
	"strings"
)

var (
    mu      sync.Mutex
    players = make(map[string]net.Conn)
)

// function pour démarrer le serveur
func Start(port string) error {
	// la function Listen ouvre un socket sur le port indiqué et renvoie un objet ln ainsi qu éventuellement un message d erreur.
	ln, err := net.Listen("tcp", port)
	if err != nil {
		// à la différence de log (texte libre), revoie des logs structurés en json
		slog.Error("Cannot open port", "port", port, "err", err)
		return err
	}
	// on anticipe la fermeture du socket juste avant le retour
	// On renvoie directement err sans fmt.Errorf car slog.Error a déjà enregistré tous les détails
	defer ln.Close()
	// On utilise slog plutôt que fmt.Println car le sujet impose des logs structurés (horodatage et niveaux INFO/WARN/ERROR)
	slog.Info("TAP server started", "port", port)
	// on met en route une boucle infinie
	for {
		// on tente le three-way handshake avec le client. On récupère un objet conn et éventuellement une erreur
		conn, err := ln.Accept()
		if err != nil {
			// si l'erreur est une ErrClosed
			if errors.Is(err, net.ErrClosed) {
				// on arrête tout proprement
				slog.Info("Listener closed, server shutting down cleanly")
				return nil
			}
			// si l'erreur est de type net.Error et qu'elle est teporaire 9ok = True)
			netErr, ok := err.(net.Error)
			if ok && netErr.Timeout() {
				// on attend et on réessaye jusqu'à ce que cela fonctionne
				slog.Warn("Temporary network error on Accept(), pausing 10ms...", "err", err)
				time.Sleep(10 * time.Millisecond)
				continue
			}
			// si c ést une autre erreur on la retourne.
			// idem pourquoi pas un fmt.Println?
			slog.Error("Fatal error on Accept()", "err", err)
			return err
		}
		// on lance avec go une goroutine pour le client
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	slog.Info("Connected client.")
	fmt.Fprintf(conn, "OK hello proto=1\n")

	client_scanner := bufio.NewScanner(conn)
	if !client_scanner.Scan() {
		return
	}

	line := client_scanner.Text()
	parts := strings.SplitN(line, " ", 2)
	if len(parts) != 2 || strings.ToUpper(parts[0]) != "CONNECT" {
		fmt.Fprintf(conn, "ERR 400 BAD_REQUEST\n")
		return
	}
	username := strings.TrimSpace(parts[1])
	if username == "" {
		fmt.Fprintf(conn, "ERR 400 BAD_REQUEST\n")
		return
	}

	mu.Lock()
	if _, exists := players[username]; exists {
		fmt.Fprintf(conn, "ERR 201 NAME_IN_USE\n")
		mu.Unlock()
		return
	}
	players[username] = conn

	mu.Unlock()
	fmt.Fprintf(conn, "OK connected\n")
	slog.Info("Player anthenticated", "username", username)

	defer func() {
		mu.Lock()
		delete(players, username)
		slog.Info("Player disconnected", "username", username)
		mu.Unlock()
	}()
}

