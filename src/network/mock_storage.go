package network

type MockStorage struct {
	ReadResult string
	ReadErr    error
}

func (m *MockStorage) Read(key string, ch chan string, errCh chan error) {
	if m.ReadErr != nil {
		errCh <- m.ReadErr
	} else {
		ch <- m.ReadResult
	}
	return
}

func (m *MockStorage) Write(key, data string) {
	return
}
