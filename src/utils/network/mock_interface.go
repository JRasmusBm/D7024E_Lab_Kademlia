package network

import (
	"net"
)

type MockInterface struct {
	addrs        []Address
	addrErr      error
	flagUp       net.Flags
	flagLoopback net.Flags
}

func (m *MockInterface) Addrs() ([]Address, error) {
	return m.addrs, m.addrErr
}

func (m *MockInterface) FlagUp() net.Flags {
	return m.flagUp
}

func (m *MockInterface) FlagLoopback() net.Flags {
	return m.flagLoopback
}
