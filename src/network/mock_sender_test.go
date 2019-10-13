package network

import (
	"errors"
	"testing"
)

func TestMockSenderPingResponse(t *testing.T) {
	expected := true
	sender := MockSender{PingResponse: expected}
	ch := make(chan bool)
	errCh := make(chan error)
	go sender.Ping(nil, ch, errCh)
	actual := <-ch
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestMockSenderPingErr(t *testing.T) {
	sender := MockSender{PingErr: errors.New("TestMockSenderPingErr")}
	ch := make(chan bool)
	errCh := make(chan error)
	go sender.Ping(nil, ch, errCh)
	err := <-errCh
	if err == nil {
		t.Errorf("Should return the error")
	}
}
