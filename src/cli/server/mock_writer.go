package cli

type MockWriter struct {
	WriteErr error
	ch       chan []byte
}

func (m *MockWriter) Write(p []byte) (int, error) {
	m.ch <- p
	return 0, m.WriteErr
}
