package main

import (
	"time"
)

// Alias time.Duration so that we can use a custom unmarshaler.
type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

type SNMP struct {
	Address  string
	Username string

	// Use SHA authentication with this password.
	Password string

	// Use AES encryption with this key.
	Key string
}

type Outlet struct {
	Name  string
	OID   string
	UDP   string
	Delay duration
}

type Configuration struct {
	SNMP   SNMP
	Outlet []Outlet
}
