package main

import (
	"log"
	"time"

	"github.com/pietern/pductl/pdu"
)

func main() {
	addr := "192.168.1.24:23"
	p, err := pdu.Dial("tcp", addr, time.Second)
	if err != nil {
		log.Fatalf("Error dialing %s: %v", addr, err)
	}

	err = p.Authenticate("apc", "apc")
	if err != nil {
		log.Fatalf("Error authenticating: %v", err)
	}

	// Run whoami (can be used for keepalive)
	whoami, err := p.Whoami()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	log.Printf("Who am I? I am \"%s\"!", whoami)

	// Fetch status of all outlets
	lines, err := p.Run("status", "all")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for _, line := range lines {
		log.Println(line)
	}
}
