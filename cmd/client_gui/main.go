package main

import (
	"log"
	"os"
	"tap/internal/gui"
)

func main() {
	username := "player"
	if len(os.Args) > 1 {
		username = os.Args[1]
	}

	if err := gui.Run("localhost:4242", username); err != nil {
		log.Fatal(err)
	}
}

