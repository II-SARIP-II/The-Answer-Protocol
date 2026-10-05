package main

import (
	"log"
	"tap/internal/server"
)

func main() {
	if error := server.Start(":4242"); error != nil {
		log.Fatal(error)
	}
}
