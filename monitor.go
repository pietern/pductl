package main

import (
	"log"
	"net"

	"github.com/pietern/pductl/watchdog"
)

// Monitor wraps a watchdog and kicks it when it receives a packet on
// a UDP socket that it listens on.
type Monitor struct {
	*watchdog.Watchdog
}

func NewMonitor(outlet Outlet) (*Monitor, error) {
	addr, err := net.ResolveUDPAddr("udp", outlet.UDP)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	log.Printf(
		"New watchdog for %s with delay %v",
		outlet.Name,
		outlet.Delay.Duration)
	m := Monitor{
		Watchdog: watchdog.NewWatchdog(outlet.Delay.Duration),
	}

	go func() {
		// Assume packets are never bigger...
		// We don't care what they contain anyway.
		buf := make([]byte, 1024)
		for {
			_, _, err := conn.ReadFrom(buf)
			if err != nil {
				log.Fatal(err)
			}

			// Received packet, kick watchdog.
			m.Watchdog.Kick()
		}
	}()

	return &m, nil
}
