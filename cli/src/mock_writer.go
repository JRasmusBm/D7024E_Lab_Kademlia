package main

type MockWriter struct {
  WriteErr error
	ch chan []byte
}

func (m *MockWriter) Write(val []byte) (int, error) {
	m.ch <- val
	return 0, m.WriteErr
}
