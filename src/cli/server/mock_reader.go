package cli

type MockReader struct {
	ReadStringResult string
	ReadStringErr    error
}

func (m *MockReader) ReadString(delim byte) (string, error) {
	if m.ReadStringErr != nil {
		return "", m.ReadStringErr
	}
	return m.ReadStringResult, nil
}
