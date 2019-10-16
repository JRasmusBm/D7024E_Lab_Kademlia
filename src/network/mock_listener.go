package network

import (
	"io"
)

type MockListener struct {
	AcceptResult MockReadWriter
	AcceptErr error
}

func (m MockListener) Accept() (io.ReadWriter, error) {
	var acceptResult io.ReadWriter = &m.AcceptResult
	return acceptResult, m.AcceptErr
}
