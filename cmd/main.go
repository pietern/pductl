package main

import (
	"fmt"
	"log"

	"github.com/pietern/pductl/snmp"
)

func main() {
	conn, err := snmp.Dial(snmp.Config{
		Address:  "192.168.1.240",
		Username: "username",
		Password: "password",
		Key:      "password",
	})

	if err != nil {
		panic(err)
	}

	result, err := conn.WalkAll(".1.3.6.1.4.1.318.1.1.12.3.3.1.1.4")
	if err != nil {
		panic(err)
	}

	for _, pdu := range result {
		fmt.Println(pdu.Name, pdu.Value)
	}

	// Set the outlet to on.
	log.Printf("Setting outlet to on")
	err = conn.Set(".1.3.6.1.4.1.318.1.1.12.3.3.1.1.4.3", 2)
	if err != nil {
		panic(err)
	}
}
