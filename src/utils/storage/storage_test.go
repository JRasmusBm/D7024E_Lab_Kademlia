package storage

import (
	"testing"
	"utils/hashing"
)

func TestLookupSucceeds(t *testing.T) {
	storage := RealStorage{Data: make(map[string]string)}
	expected := "abc"
	storage.Write(hashing.NewKademliaID(expected).String(), expected)
	ch := make(chan string)
	errCh := make(chan error)
	go storage.Read(hashing.NewKademliaID(expected).String(), ch, errCh)
	actual := <-ch
	if expected != actual {
		t.Errorf("Expected %s got %s", expected, actual)
	}
}

func TestLookupFails(t *testing.T) {
	storage := RealStorage{Data: make(map[string]string)}
	expected := "abc"
	ch := make(chan string)
	errCh := make(chan error)
	go storage.Read(hashing.NewKademliaID(expected).String(), ch, errCh)
  var err error
	select {
	case err = <-errCh:
	// Do nothing
	case <-ch:
		// Do nothing
	}
	if err == nil {
		t.Errorf("Should throw error")
	}
}
