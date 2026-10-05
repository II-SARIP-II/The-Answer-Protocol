package main

import (
	"log"

	"tap/internal/gui"
)

func main() {
	if err := gui.Run(); err != nil {
		log.Fatal(err)
	}
}
