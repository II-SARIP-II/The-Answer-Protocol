package main

import (
	"log"
	"tap/internal/gui"
	"os"
	"strings"
)

func main() {
	addr := "localhost:4242"
	if len(os.Args) > 1 {
		addr = os.Args[1]
		if !strings.Contains(addr, ":") {
			addr = "localhost:" + addr
		}
	}
	if err := gui.Run(addr); err != nil {
		log.Fatal(err)
	}
}
