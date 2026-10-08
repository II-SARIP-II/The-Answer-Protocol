package cli

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"tap/internal/protocol"
)

// en go on renvoie un objet de type erreur quand une erreur est possible
func Run(addr string) error {
	// Dial est une fonction qui ouvre une socket à une adresse et sur un port donnés
	conn, err := net.Dial("tcp", addr)
	// si il y a une erreur, on la recupère proprement avec la fonction Errorf qui formate un message d'erreur
	// grâce à %w l'erreur est emballée, conservée et associée au message
	if err != nil {
		return fmt.Errorf("Cannot connect to the TAP server on %s: %w", addr, err)
	}
	// on prévoit la fermeture de la socket juste avant le retour
	defer conn.Close()
	fmt.Println("Successfully connected to the TAP server.")

	// on crée un scanner pour lire les messages que le serveur écrit dans la socket
	serverScanner := bufio.NewScanner(conn)
	if serverScanner.Scan() {
		fmt.Println("Server:", serverScanner.Text())
	}
	fmt.Println("Enter your user name: ")

	// on crée un scanner qui va lire les entrées clavier
	stdinScanner := bufio.NewScanner(os.Stdin)
	// si la lecture échoue, message d'erreur
	if !stdinScanner.Scan() {
		return fmt.Errorf("Failed to read username from stdin.")
	}
	// sinon on enregistre le nom d'utilisateur dans une variable et on supprime les espaces éventuels
	username := strings.TrimSpace(stdinScanner.Text())

	// on envoie au serveur la commande CONNECT avec le nom de l'utilisateur
	fmt.Fprint(conn, protocol.FormatCmd(protocol.CmdConnect, username))

	// on attend et affiche la réponse du serveur (OK connected ou ERR...)
	if serverScanner.Scan() {
		response := serverScanner.Text()
		fmt.Println("Server:", response)
		if strings.HasPrefix(response, protocol.PrefixERR) {
			return nil
		}
	}

	fmt.Println("You can now enter commands (or QUIT to exit):")

	// nécessaire d'ouvrir une go routine en arrière plan qui ecoute le serveru en continu
	go func() {
		for serverScanner.Scan() {
			response := serverScanner.Text()
			fmt.Println("Server:", response)
		}
		if err := stdinScanner.Err(); err != nil {
			fmt.Println("Error reading standard input: %w", err)
		}
		os.Exit(0)
	}()

	for stdinScanner.Scan() {
		input := strings.TrimSpace(stdinScanner.Text())

		if input == "" {
			continue
		}
		if strings.ToUpper(input) == protocol.CmdQuit {
			fmt.Fprint(conn, protocol.FormatCmd(protocol.CmdQuit, ""))
			if serverScanner.Scan() {
				fmt.Println("Server:", serverScanner.Text())
			}
			break
		}
		fmt.Fprintf(conn, "%s\n", input)

		if serverScanner.Scan() {
			fmt.Println("Server:", serverScanner.Text())
		}
	}
	if err := stdinScanner.Err(); err != nil {
		return fmt.Errorf("Error reading standard input: %w", err)
	}
	if err := serverScanner.Err(); err != nil {
		fmt.Println("Connection to server lost: %w", err)
	}
	return nil
}
