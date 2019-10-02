package main

import (
	"net"
)

type MockNetwork struct {
	DialResult net.Conn
	DialErr    error
}

func (r *MockNetwork) Dial(network, address string) (net.Conn, error) {
	return r.DialResult, nil
}
