package network

type MockAddress struct {
	ip IP
}

func (m *MockAddress) IP() IP {
	return m.ip
}
