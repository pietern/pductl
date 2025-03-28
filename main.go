package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/pietern/pductl/snmp"
	"github.com/pietern/pductl/watchdog"
)

type wrapper struct {
	config Configuration
	snmp   *snmp.Connection
}

func newWrapper(config Configuration) *wrapper {
	var err error

	w := &wrapper{config: config}
	w.snmp, err = snmp.Dial(snmp.Config{
		Address:  config.SNMP.Address,
		Username: config.SNMP.Username,
		Password: config.SNMP.Password,
		Key:      config.SNMP.Key,
	})

	if err != nil {
		log.Fatal(err)
	}

	return w
}

func (w *wrapper) PowerOn(oid string) error {
	return w.snmp.Set(oid, 1)
}

func (w *wrapper) PowerOff(oid string) error {
	return w.snmp.Set(oid, 2)
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

	// Convenience wrapper to power outlets on and off.
	w := newWrapper(c)

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

	// // Run something to see if the credentials are valid.
	// err = w.Run("whoami")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Main loop.
	for t := range ch {
		switch t.State {
		case watchdog.Present:
			log.Printf("Turning %v on\n", t.Outlet.Name)
			err := w.PowerOn(t.Outlet.OID)
			if err != nil {
				log.Fatal(err)
			}
		case watchdog.Absent:
			log.Printf("Turning %v off\n", t.Outlet.Name)
			err := w.PowerOff(t.Outlet.OID)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
