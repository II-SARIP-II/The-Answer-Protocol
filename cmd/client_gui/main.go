package main

import (
	"log"
	"tap/internal/gui"
)

func main() {
	if err := gui.Run("localhost:4242"); err != nil {
		log.Fatal(err)
	}
}
