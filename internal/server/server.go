package server

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"sync"
	"time"
)

// mot clé var pourdéclarer un bloc de declaration de variables globales
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
		// On renvoie directement err sans fmt.Errorf car slog.Error a déjà enregistré tous les détails
		return err
	}
	// on anticipe la fermeture du socket juste avant le retour

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
	// on anticipe la fermeture de la socket juste après le return
	defer conn.Close()
	// message de log formaté dans le terminal
	slog.Info("Connected client.")
	// ecrit dans la socket réseau du client
	fmt.Fprintf(conn, "OK hello proto=1\n")

	// initialise un nouveau scanner
	client_scanner := bufio.NewScanner(conn)
	// si le scanner n'a rien à scanner
	if !client_scanner.Scan() {
		return
	}

	// recuperation de ce qui se trouve dans le buffer
	line := client_scanner.Text()
	// decoupage de la ligne en 2 parties (permet de gérer les espaces éventuels dans le username)
	parts := strings.SplitN(line, " ", 2)
	// si pb dans la synatxe avec gestion casse insensible
	if len(parts) != 2 || strings.ToUpper(parts[0]) != "CONNECT" {
		fmt.Fprintf(conn, "ERR 400 BAD_REQUEST\n")
		return
	}
	// on enlève les espaces
	username := strings.TrimSpace(parts[1])
	// gesion chîne vide
	if username == "" {
		fmt.Fprintf(conn, "ERR 400 BAD_REQUEST\n")
		return
	}

	// mutex pour protéger lécriture dans la map globale des players
	mu.Lock()
	// vérification doublon username
	if _, exists := players[username]; exists {
		fmt.Fprintf(conn, "ERR 201 NAME_IN_USE\n")
		// déverouillage du mutex si erreur avant de return
		mu.Unlock()
		return
	}
	// ajout du user dans la map
	players[username] = conn

	// déverouillage du mutex
	mu.Unlock()
	// affichage client et log serveur
	fmt.Fprintf(conn, "OK connected\n")
	slog.Info("Player anthenticated", "username", username)

	// on anticipe la deconnexion du player. Apès defer, fonction anonyme dont on lance l'exécution avec  ()
	defer func() {
		mu.Lock()
		delete(players, username)
		slog.Info("Player disconnected", "username", username)
		mu.Unlock()
	}()

	// boucle infinie qui s'arrêtera lors de la déconnexion(en veille donc ne consomme pas de CPU), lorsque la fonction Scan retournera False
	// exécution des defer et fermeture propre. La go routine se ferme automatiquement. Rien à gérer !
	for client_scanner.Scan() {
		command := strings.TrimSpace(client_scanner.Text())

		if command == "" {
			continue
		}

		slog.Info("Command received","username", username,  "cmd", command)

		if strings.ToUpper(command) == "QUIT" {
			fmt.Fprintf(conn, "OK bye\n")
			break
		}
		fmt.Fprintf(conn, "OK received (TODO) %s\n", command)
	}
	if err := client_scanner.Err(); err != nil {
		slog.Warn("Client connexion error", "username", username, "error", err)
	} else {
		slog.Info("Client cleanly disconnected", "username", username)
	}
}

