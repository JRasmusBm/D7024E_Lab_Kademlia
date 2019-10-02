package cli

import (
	"net"
)

type RealNetwork struct {
}

func (r *RealNetwork) Listen(network, address string) (net.Listener, error) {
	return net.Listen(network, address)
}
