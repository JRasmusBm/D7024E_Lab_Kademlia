package network

import (
	"errors"
	"testing"
	nodeutils "utils/node"
)

func TestMockSenderPingResponse(t *testing.T) {
	expected := true
	sender := MockSender{PingResponse: expected}
	ch := make(chan bool)
	errCh := make(chan error)
	go sender.Ping(nodeutils.Node{IP: "0.0.0.0"}, ch, errCh)
	var err error
	var actual bool
	select {
	case err = <-errCh:
	// Do nothing
	case actual = <-ch:
		// Do nothing
	}
  if err != nil {
    t.Errorf(err.Error())
  }
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestMockSenderPingErr(t *testing.T) {
	sender := MockSender{PingErr: errors.New("TestMockSenderPingErr")}
	ch := make(chan bool)
	errCh := make(chan error)
	go sender.Ping(nodeutils.Node{IP: "0.0.0.0"}, ch, errCh)
	err := <-errCh
	if err == nil {
		t.Errorf("Should return the error")
	}
}
