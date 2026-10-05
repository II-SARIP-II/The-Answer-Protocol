package main

import (
	"log"
	"tap/internal/cli"
)

func main() {
	error := cli.Run("localhost:4242")
	if error != nil {
		log.Fatal(error)
	}
}
