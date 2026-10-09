package main

import (
	"os"
	"tap/internal/server"
	"log/slog"
)

func main() {
	serv, err := server.NewServer("data/world.json")
	if err != nil{
		slog.Error("Failed to initialize server", "error", err)
		os.Exit(1)
	}
	if err := serv.Start(":4242"); err != nil {
		os.Exit(1)
	}
}
