package gui

import (
	"bufio"
	"fmt"
	"net"
)

func listenServer(conn net.Conn, messagesChan chan string) {
	serverScanner := bufio.NewScanner(conn)

	for serverScanner.Scan() {
		text := serverScanner.Text()
		messagesChan <- text
	}

	if err := serverScanner.Err(); err != nil {
		messagesChan <- "Connection lost: " + err.Error()
	}
}

func sendMessage(conn net.Conn, text string) {
	if conn != nil {
		fmt.Fprintf(conn, "%s\n", text)
	}
}
