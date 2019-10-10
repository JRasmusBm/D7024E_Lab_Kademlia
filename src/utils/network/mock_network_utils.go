package network

type MockNetworkUtils struct {
	IP  string
	Err error
}

func (r *MockNetworkUtils) GetIP() (string, error) {
	if r.Err != nil {
		return "", r.Err
	}
	return r.IP, nil
}
