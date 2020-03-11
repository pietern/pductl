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

type PDU struct {
	Address  string
	Username string
	Password string
	Timeout  duration
}

type Outlet struct {
	Name  string
	UDP   string
	Delay duration
}

type Configuration struct {
	PDU    PDU
	Outlet []Outlet
}
