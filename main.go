package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/pietern/pductl/pdu"
	"github.com/pietern/pductl/watchdog"
)

type wrapper struct {
	c Configuration
}

func (w *wrapper) Run(args ...string) error {
	pdu, err := pdu.Dial("tcp", w.c.PDU.Address, w.c.PDU.Timeout.Duration)
	if err != nil {
		return err
	}

	err = pdu.Authenticate(w.c.PDU.Username, w.c.PDU.Password)
	if err != nil {
		return err
	}

	_, err = pdu.Run(args...)
	if err != nil {
		return err
	}

	return pdu.Logout()
}

func main() {
	var path string
	flag.StringVar(&path, "config", "", "Path to configuration file")
	flag.Parse()

	var c Configuration
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		log.Fatal(err)
	}

	// Convenience wrapper to connect, authenticate, run a
	// command, and logout again.
	//
	// It proves difficult to keep a single connection alive for a
	// longer period of time, so instead we just recreate a
	// connection every time we need to perform some action.
	//
	w := &wrapper{c}

	// Create aggregate channel for tuples containing both the
	// presence or absence signal and the outlet information.
	type Tuple struct {
		State  watchdog.State
		Outlet Outlet
	}

	ch := make(chan Tuple)
	for _, outlet := range c.Outlet {
		monitor, err := NewMonitor(outlet)
		if err != nil {
			log.Fatal(err)
		}

		go func(monitor *Monitor, outlet Outlet) {
			for state := range monitor.C {
				ch <- Tuple{State: state, Outlet: outlet}
			}
		}(monitor, outlet)
	}

	// Run something to see if the credentials are valid.
	err = w.Run("whoami")
	if err != nil {
		log.Fatal(err)
	}

	// Main loop.
	for t := range ch {
		switch t.State {
		case watchdog.Present:
			log.Printf("Turning %v on\n", t.Outlet.Name)
			err := w.Run("on", t.Outlet.Name)
			if err != nil {
				log.Fatal(err)
			}
		case watchdog.Absent:
			log.Printf("Turning %v off\n", t.Outlet.Name)
			err := w.Run("off", t.Outlet.Name)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
