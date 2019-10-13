package network

import (
	"io"
)

type MockDialer struct {
	DialResults []io.ReadWriter
	DialErrors  []error
	index       int
}

func (m *MockDialer) DialIP(ip string) (io.ReadWriter, error) {
	result := m.DialResults[m.index]
	err := m.DialErrors[m.index]
	m.index += 1
	return result, err
}
