package network

import (
	"io"
	"time"
)

type MockListener struct {
	AcceptResult MockReadWriter
	AcceptErr    error
}

func (m MockListener) Accept() (io.ReadWriter, error) {
	time.Sleep(300 * time.Millisecond)
	var acceptResult io.ReadWriter = &m.AcceptResult
	return acceptResult, m.AcceptErr
}
