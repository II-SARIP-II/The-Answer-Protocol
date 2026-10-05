package main

import (
	"log"
	"tap/internal/server"
)

func main() {
	error := server.Start(":4242")
	if error != nil {
		log.Fatal(error)
	}
}
