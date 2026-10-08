package main

import (
	"log"
	"tap/internal/cli"
)

func main() {
	if error := cli.Run("localhost:4242"); error != nil {
		log.Fatal(error)
	}
}
