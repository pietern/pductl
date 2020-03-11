package main

import (
	"flag"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/pietern/pductl/pdu"
	"github.com/pietern/pductl/watchdog"
)

func main() {
	var path string
	flag.StringVar(&path, "config", "", "Path to configuration file")
	flag.Parse()

	var c Configuration
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		log.Fatal(err)
	}

	pdu, err := pdu.Dial("tcp", c.PDU.Address, c.PDU.Timeout.Duration)
	if err != nil {
		log.Fatal(err)
	}

	err = pdu.Authenticate(c.PDU.Username, c.PDU.Password)
	if err != nil {
		log.Fatal(err)
	}

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

	// Main loop.
	keepalive := time.NewTicker(c.PDU.Timeout.Duration)
	for {
		select {
		case t := <-ch:
			switch t.State {
			case watchdog.Present:
				log.Printf("Turning %v on\n", t.Outlet.Name)
				_, err := pdu.Run("on", t.Outlet.Name)
				if err != nil {
					log.Fatal(err)
				}
			case watchdog.Absent:
				log.Printf("Turning %v off\n", t.Outlet.Name)
				_, err := pdu.Run("off", t.Outlet.Name)
				if err != nil {
					log.Fatal(err)
				}
			}
		case <-keepalive.C:
			_, err := pdu.Whoami()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
