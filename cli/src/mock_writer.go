package main

type MockWriter struct {
}

func (m *MockWriter) Write(val []byte) (int, error) {
	return 0, nil
}
