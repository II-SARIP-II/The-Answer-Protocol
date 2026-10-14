package main

import (
	"log/slog"
	"os"
	"tap/internal/server"
)

func main() {
	serv, err := server.NewServer("data/world.json")
	if err != nil {
		slog.Error("Failed to initialize server", "error", err)
		os.Exit(1)
	}
	if err := serv.Start(":4242"); err != nil {
		os.Exit(1)
	}
}
