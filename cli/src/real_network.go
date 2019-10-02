package main

import (
	"net"
)

type RealNetwork struct {
}

func (r *RealNetwork) Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}
