package api

import (
	"testing"
	"utils/constants"
	"utils/hashing"
	nodeutils "utils/node"
)

func TestPingNonExistantAddress(t *testing.T) {
	node := nodeutils.Node{IP: ""}
	value := Ping(&node)
	if value {
		t.Errorf("Should return an error")
	}
}

func TestFindNoNodes(t *testing.T) {
	node := nodeutils.Node{IP: ""}
	id := hashing.NewKademliaID("abc")
	val, err := FindNode(&node, id)
	if err == nil || len(val) != constants.CLOSESTNODES {
		t.Errorf("Expected list with %d spaces, got %#v", constants.CLOSESTNODES, val)
	}
}

func TestFindValueNotImplemented(t *testing.T) {
	key := hashing.NewKademliaID("abc")
	expected := "Random value"
	actual := FindValue(key)
	if actual != expected {
		t.Errorf("Expected %#v got %#v", expected, actual)
	}
}

func TestStoreNotImplemented(t *testing.T) {
	content := "abc"
	expected := hashing.NewKademliaID(content)
	actual := Store(content)
	if actual.String() != expected.String() {
		t.Errorf("Expected %#v got %#v", expected.String(), actual.String())
	}
}

func TestJoinNotImplemented(t *testing.T) {
  Join()
}
