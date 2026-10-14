package gui

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"tap/internal/protocol"
)

func RunGui(addr string, username string) error {
	// Dial = ouverture d'une socket à une adresse et sur un port donnés
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("Cannot connect to the TAP server on %s: %w", addr, err)
	}
	defer conn.Close()
	// pour l instant affichage terminal, à remplacer si tu veux par des fenêtres popup
	fmt.Println("Successfully connected to the TAP server.")

	// on crée un scanner pour lire les messages que le serveur écrit dans la socket
	serverScanner := bufio.NewScanner(conn)
	if serverScanner.Scan() {
		fmt.Println("Server:", serverScanner.Text())
	}

	fmt.Fprint(conn, protocol.FormatCmd(protocol.CmdConnect, username))

	// Validation de l'authentification : OK connected ou ERR
	if serverScanner.Scan() {
		response := serverScanner.Text()
		// pour l instant affichage terminal, à remplacer si tu veux par des fenêtres popup
		fmt.Println("Server:", response)
		if strings.HasPrefix(response, protocol.PrefixERR) {
			return fmt.Errorf("authentication failed: %s", response)
		}
	}

	// Goroutine en arrière-plan qui écoute les messages du serveur
	go func() {
		for serverScanner.Scan() {
			response := serverScanner.Text()
			// pour l instant affichage terminal, à remplacer si tu veux par des fenêtres popup
			fmt.Println("Server:", response)
			// TODO: mise à jour fenêtre gui
		}
		os.Exit(0)
	}()

	// TODO: La boucle graphique de ta fenêtre prend le relais ici.
	// Pour envoyer une commande lors d'un clic de bouton :
	// fmt.Fprintf(conn, "%s\n", "MOVE north") ou fmt.Fprintf(conn, "%s\n", "LOOK")
	// Pour éviter les string en dur, tu peux utiliser les constantes de protocol.go
	// (par exemple protocol.CmdMove au lieu de "MOVE")

	return nil
}
