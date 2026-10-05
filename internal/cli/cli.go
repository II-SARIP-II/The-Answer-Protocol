package cli

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// en go on renvoie un objet de type erreur quand une erreur est possible
func Run(addr string) error {
	// Dial est une fonctiopn qui ouvre un socket à une adresse et sur un port donnés
	conn, err := net.Dial("tcp", addr)
	//síl y a une erreur, on la recupere proprement avec la fonction Errorf qui formate un message d'erreur
	//grâce à %w l 'erreur est emballée, conservée et associée au message
	if err != nil {
		return fmt.Errorf("Cannot connect to the TAP server on %s: %w", addr, err)
	}
	//on prévoit la fermeture du socket juste avant le retour
	defer conn.Close()
	fmt.Println("Successfully connected to the TAP server.")
	//on cree un scanner pour lire les messages que le serveur écrit dans le socket
	serverScanner := bufio.NewScanner(conn)
	if serverScanner.Scan() {
		fmt.Println("Server:", serverScanner.Text())
	}
	fmt.Println("Enter your user name: ")
	// on crée un scanner qui va lire les entrées clavier
	stdinScanner := bufio.NewScanner(os.Stdin)
	// si la lecture echoue, message d'erreur
	if !stdinScanner.Scan() {
		return fmt.Errorf("Failed to read username form stdin.")
	}
	//  sinon on enregistre le nom d útilisateur dans une variable et on supprime les espaces éventuels
	username := strings.TrimSpace(stdinScanner.Text())
	//on envoie au serveur le nom de l'utilisateur
	fmt.Fprintf(conn, "CONNECT %s\n", username)
	// si le serveur ecrit dans le socket, on l affiche le message sur la sortie standard
	if serverScanner.Scan() {
		response := serverScanner.Text()
		fmt.Println(response)
	}

	return nil
}
