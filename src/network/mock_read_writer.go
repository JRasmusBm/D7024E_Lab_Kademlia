package network

import (
	"io"
)

type MockReadWriter struct {
	Msg      []byte
	WriteErr error
	WriteCh  chan []byte
	ReadErr  error
	i        int64
}

func (m *MockReadWriter) Read(b []byte) (n int, err error) {
	if m.i >= int64(len(m.Msg)) {
		return 0, io.EOF
	}
	n = copy(b, m.Msg[m.i:])
	m.i += int64(n)
	return
}

func (m *MockReadWriter) Write(p []byte) (n int, err error) {
  go func() {
    m.WriteCh <- p
  }()
	return 0, m.WriteErr
}
