package cli

type MockCloser struct {
}

func (m *MockCloser) Close() error {
	return nil
}
