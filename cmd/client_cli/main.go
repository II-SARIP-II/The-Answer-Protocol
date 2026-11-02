package main

import (
	"log"
	"os"
	"strings"
	"tap/internal/cli"
)

func main() {
	addr := "localhost:4242"
	if len(os.Args) > 1 {
		addr = os.Args[1]
		if !strings.Contains(addr, ":") {
			addr = "localhost:" + addr
		}
	}
	if err := cli.Run(addr); err != nil {
		log.Fatal(err)
	}
}
