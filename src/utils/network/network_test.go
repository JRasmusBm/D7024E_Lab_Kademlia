package network

import (
	"errors"
	"testing"
)

func TestGetAddress(t *testing.T) {
	addrs := []Address{
		&MockAddress{ip: &MockIP{str: "1.2.3.4", loopback: false, hasipv4: true}},
	}
	ifces := []Interface{
		&MockInterface{addrs: addrs, addrErr: nil, flagUp: 1, flagLoopback: 0},
	}
	var e error = nil
	expected := "1.2.3.4"
	actual, err := extractIP(ifces, e)
	if err != nil || expected != actual {
		t.Errorf("Expected %s got %s", expected, actual)
	}
}

func TestEmptyAddressList(t *testing.T) {
	addrs := []Address{}
	ifces := []Interface{
		&MockInterface{addrs: addrs, addrErr: nil, flagUp: 1, flagLoopback: 0},
	}
	var e error = nil
	_, err := extractIP(ifces, e)
	if err == nil {
		t.Errorf("Expected to return an error")
	}
}

func TestWithLoopbackFlag(t *testing.T) {
	addrs := []Address{
		&MockAddress{ip: &MockIP{str: "1.2.3.4", loopback: false, hasipv4: true}},
	}
	ifces := []Interface{
		&MockInterface{addrs: addrs, addrErr: nil, flagUp: 1, flagLoopback: 1},
	}
	var e error = nil
	_, err := extractIP(ifces, e)
	if err == nil {
		t.Errorf("Expected to return an error")
	}
}

func TestIPIsLoopback(t *testing.T) {
	addrs := []Address{
		&MockAddress{ip: &MockIP{str: "1.2.3.4", loopback: true, hasipv4: true}},
	}
	ifces := []Interface{
		&MockInterface{addrs: addrs, addrErr: nil, flagUp: 1, flagLoopback: 0},
	}
	var e error = nil
	_, err := extractIP(ifces, e)
	if err == nil {
		t.Errorf("Should return error")
	}
}

func TestNoIPv4(t *testing.T) {
	addrs := []Address{
		&MockAddress{ip: &MockIP{str: "1.2.3.4", loopback: false, hasipv4: false}},
	}
	ifces := []Interface{
		&MockInterface{addrs: addrs, addrErr: nil, flagUp: 1, flagLoopback: 0},
	}
	var e error = nil
	_, err := extractIP(ifces, e)
	if err == nil {
		t.Errorf("Should throw an error")
	}
}

func TestNotUp(t *testing.T) {
	addrs := []Address{
		&MockAddress{ip: &MockIP{str: "1.2.3.4", loopback: false, hasipv4: true}},
	}
	ifces := []Interface{
		&MockInterface{addrs: addrs, addrErr: nil, flagUp: 0, flagLoopback: 0},
	}
	var e error = nil
	_, err := extractIP(ifces, e)
	if err == nil {
		t.Errorf("Expected to return an error")
	}
}

func TestWithError(t *testing.T) {
	addrs := []Address{
		&MockAddress{ip: &MockIP{str: "1.2.3.4", loopback: false, hasipv4: true}},
	}
	ifces := []Interface{
		&MockInterface{addrs: addrs, addrErr: nil, flagUp: 1, flagLoopback: 0},
	}
	var e error = errors.New("Hello")
	expected := "1.2.3.4"
	actual, err := extractIP(ifces, e)
	if err == nil {
		t.Errorf("Expected %s got %s", expected, actual)
	}
}

func TestAddrsReturnsError(t *testing.T) {
	addrs := []Address{
		&MockAddress{ip: &MockIP{str: "1.2.3.4", loopback: false, hasipv4: true}},
	}
	ifces := []Interface{
		&MockInterface{
			addrs:        addrs,
			addrErr:      errors.New("Hello"),
			flagUp:       1,
			flagLoopback: 0,
		},
	}
	var e error = nil
	_, err := extractIP(ifces, e)
	if err == nil {
		t.Errorf("Expected to throw")
	}
}
