package gui

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"tap/internal/protocol"
)
// conn.Close() tjrs faire ca a la fin !!!


func Connect(addr string, username string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("Cannot connect to the TAP server on %s: %w", addr, err)
	}
	serverScanner := bufio.NewScanner(conn)
	if serverScanner.Scan() {
		fmt.Println("Server:", serverScanner.Text())
	}

	fmt.Fprint(conn, protocol.FormatCmd(protocol.CmdConnect, username))

	if serverScanner.Scan() {
		response := serverScanner.Text()
		fmt.Println("Server:", response)
		if strings.HasPrefix(response, protocol.PrefixERR) {
			conn.Close()
			return nil, fmt.Errorf("erreur serveur : %s", response)
		}
	}
	fmt.Println("Successfully connected to the TAP server.")
	return conn, nil
}

func ListenServer(conn net.Conn, messagesChan chan string) {	
	serverScanner := bufio.NewScanner(conn)

    for serverScanner.Scan() {
        text := serverScanner.Text()
        messagesChan <- text
    }

    if err := serverScanner.Err(); err != nil {
        messagesChan <- "Connection lost: " + err.Error()
    }
}

func SendMessage(conn net.Conn, text string) {
    if conn != nil {
        fmt.Fprintf(conn, "%s\n", text)
    }
}
