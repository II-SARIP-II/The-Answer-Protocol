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

	defer conn.Close()
	fmt.Println("Successfully connected to the TAP server.")

	serverScanner := bufio.NewScanner(conn)

	if serverScanner.Scan() {
		fmt.Println("Server:", serverScanner.Text())
	}

	stdinScanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Enter your user name: ")
		if !stdinScanner.Scan() {
			return fmt.Errorf("failed to read username from stdin")
		}

		username := strings.TrimSpace(stdinScanner.Text())

		if username == "" {
			fmt.Println("Error: Username cannot be empty. Please try again.")
		}
		if strings.Contains(username, " ") {
			fmt.Println("Error: Username cannot contain spaces. Please try again.")
		}

		fmt.Fprint(conn, protocol.FormatCmd(protocol.CmdConnect, username))

		if !serverScanner.Scan() {
			return fmt.Errorf("connection lost with server")
		}

		response := serverScanner.Text()
		if response == "OK connected" {
			fmt.Println("Server:", response)
			break
		}

		if response == "ERR 201 NAME_IN_USE" {
			fmt.Println("Error: This username is already taken by another player. Please choose another one.")
		} 
	}

	fmt.Println("You can now enter commands (or QUIT to exit):")

	// channel de signalisation (structure vide utilisée en go car 0 octet !)
	// permet de fermer la socket
	quitChan := make(chan struct{})
	go func() {
		for serverScanner.Scan() {
			response := serverScanner.Text()

			if response == "" {
				continue
			}
			fmt.Println("Server:", response)

			if response == "OK bye" {
				close(quitChan)
			}
		}
	}()

	for stdinScanner.Scan() {
		input := strings.TrimSpace(stdinScanner.Text())

		if input == "" {
			continue
		}
		if strings.ToUpper(input) == protocol.CmdQuit {
			fmt.Fprint(conn, protocol.FormatCmd(protocol.CmdQuit, ""))
			// avec <- on extrait une valeur du channel
			<-quitChan
				fmt.Println("Disconnected.")
				return nil
		}
		fmt.Fprintf(conn, "%s\n", input)
	}
	if err := stdinScanner.Err(); err != nil {
		return fmt.Errorf("Error reading standard input: %w", err)
	}
	if err := serverScanner.Err(); err != nil {
		fmt.Println("Connection to server lost: %w", err)
	}
	return nil
}
