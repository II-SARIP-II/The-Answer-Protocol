package main

import (
	"tap/internal/server"
	"os"

)

func main() {
	if error := server.Start(":4242"); error != nil {
		os.Exit(1)
	}
}
