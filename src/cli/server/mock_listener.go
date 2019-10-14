package cli

import (
	"errors"
	"net"
)

type MockListener struct {
}

func (m *MockListener) Accept() (net.Conn, error) {
	return nil, errors.New("Mock Accept")
}

func (m *MockListener) Close() error {
	return errors.New("Mock Close")
}

func (m *MockListener) Addr() net.Addr {
	return nil
}
