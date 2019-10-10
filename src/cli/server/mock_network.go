package cli

import "net"

type MockNetwork struct {
	ListenResult net.Listener
	ListenErr    error
}

func (m *MockNetwork) Listen(network, address string) (net.Listener, error) {
	return m.ListenResult, nil
}
