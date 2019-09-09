package network

type MockIP struct {
	str      string
	loopback bool
	hasipv4  bool
}

func (m *MockIP) String() string {
	return m.str
}

func (m *MockIP) IsLoopback() bool {
	return m.loopback
}

func (m *MockIP) To4() IP {
	if m.hasipv4 {
		return m
	}
	return nil
}
