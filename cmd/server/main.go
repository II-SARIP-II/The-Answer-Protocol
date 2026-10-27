package main

import (
	"log/slog"
	"os"
	"strings"
	"tap/internal/server"
)

func main() {
	server.SetLogger()
	
	serv, err := server.NewServer("data/world.json")
	if err != nil {
		slog.Error("Failed to initialize server", "error", err)
		os.Exit(1)
	}

	addr := ":4242"
	if len(os.Args) > 1 {
		addr = os.Args[1]
		if !strings.Contains(addr, ":") {
			addr = ":" + addr
		}
	}

	if err := serv.Start(addr); err != nil {
		os.Exit(1)
	}
}

