package main

import (
	"os"
	"tap/internal/server"
)

func main() {
	if error := server.Start(":4242"); error != nil {
		os.Exit(1)
	}
}
