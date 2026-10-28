package cli

import (
	"bufio"
	_ "embed"
	"fmt"
	"net"
	"os"
	"strings"
	"tap/internal/protocol"
)

// Run connects to the TAP server, authenticates the player, and starts the interactive CLI session.
func Run(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("cannot connect to the TAP server on %s: %w", addr, err)
	}
	defer conn.Close()
	fmt.Println("Successfully connected to the TAP server.")

	serverScanner := bufio.NewScanner(conn)
	stdinScanner := bufio.NewScanner(os.Stdin)

	// Read initial server greeting (OK hello proto=1).
	if !serverScanner.Scan() {
		return fmt.Errorf("server disconnected before greeting: %w", serverScanner.Err())
	}
	fmt.Println("Server:", serverScanner.Text())

	if err := authenticate(conn, serverScanner, stdinScanner); err != nil {
		return err
	}

	return handleCommands(conn, serverScanner, stdinScanner)
}

// authenticate prompts the user for a valid username and negotiates the CONNECT handshake.
func authenticate(conn net.Conn, serverScanner *bufio.Scanner, stdinScanner *bufio.Scanner) error {
	for {
		fmt.Print("Enter your user name: ")
		if !stdinScanner.Scan() {
			return fmt.Errorf("failed to read username from stdin")
		}

		username := strings.TrimSpace(stdinScanner.Text())

		if username == "" {
			fmt.Println("WARN: Username cannot be empty. Please try again.")
			continue
		}
		if strings.Contains(username, " ") {
			fmt.Println("WARN: Username cannot contain spaces. Please try again.")
			continue
		}

		fmt.Fprint(conn, protocol.FormatCmd(protocol.CmdConnect, username))

		if !serverScanner.Scan() {
			return fmt.Errorf("connection lost with server")
		}

		response := serverScanner.Text()
		if response == "OK connected" {
			fmt.Println("Server:", response)
			return nil
		}

		if response == "ERR 201 NAME_IN_USE" {
			fmt.Println("WARN: This username is already taken by another player. Please choose another one.")
		} else {
			fmt.Printf("Server Warning (%s). Please try again.\n", response)
		}
	}
}

// handleCommands loops over stdin, forwarding player commands and awaiting QUIT handshake.
func handleCommands(conn net.Conn, serverScanner *bufio.Scanner, stdinScanner *bufio.Scanner) error {
	fmt.Println("You can now enter commands (or QUIT to exit):")

	quitChan := make(chan struct{})

	go listenServer(serverScanner, quitChan)

	for stdinScanner.Scan() {
		input := strings.TrimSpace(stdinScanner.Text())
		if input == "" {
			continue
		}

		if strings.ToUpper(input) == protocol.CmdQuit {
			fmt.Fprint(conn, protocol.FormatCmd(protocol.CmdQuit, ""))
			<-quitChan
			fmt.Println("Disconnected.")
			return nil
		}

		if strings.EqualFold(input, "HELP") || input == "?" {
			fmt.Print(helpMessage)
			continue
		}

		fmt.Fprintf(conn, "%s\n", input)
	}

	if err := stdinScanner.Err(); err != nil {
		return fmt.Errorf("error reading standard input: %w", err)
	}
	return nil
}

// listenServer continuously prints messages received from the server until disconnected.
func listenServer(serverScanner *bufio.Scanner, quitChan chan struct{}) {
	for serverScanner.Scan() {
		response := serverScanner.Text()
		if response == "" {
			continue
		}
		fmt.Println("Server:", response)

		if response == "OK bye" {
			close(quitChan)
			return
		}
	}

	fmt.Println("Connection to server lost. Exiting...")
	os.Exit(0)
}
