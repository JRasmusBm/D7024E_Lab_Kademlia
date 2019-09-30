package api

import (
	"errors"
	"strconv"
	"testing"
	"utils/constants"
	"utils/hashing"
	nodeutils "utils/node"
)

func TestPingSuccessful(t *testing.T) {
	var nw Network = &MockNetwork{ping: true}
	node := nodeutils.Node{IP: ""}
	expected := true
	actual := ping(&nw, &node)
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			strconv.FormatBool(expected),
			strconv.FormatBool(actual))
	}
}

func TestPingFailed(t *testing.T) {
	var nw Network = &MockNetwork{ping: false}
	node := nodeutils.Node{IP: ""}
	expected := false
	actual := ping(&nw, &node)
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			strconv.FormatBool(expected),
			strconv.FormatBool(actual))
	}
}

func TestPingError(t *testing.T) {
	var nw Network = &MockNetwork{pingErr: errors.New("Random error")}
	node := nodeutils.Node{IP: ""}
	expected := false
	actual := ping(&nw, &node)
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			strconv.FormatBool(expected),
			strconv.FormatBool(actual))
	}
}

func TestStoreSuccessful(t *testing.T) {
	expected := hashing.NewKademliaID("def")
	var nw Network = &MockNetwork{store: expected}
	content := "abc"
	actual := store(&nw, content)
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
}

func TestStoreFailed(t *testing.T) {
	var expected *hashing.KademliaID = nil
	var nw Network = &MockNetwork{store: expected}
	content := "abc"
	actual := store(&nw, content)
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
}

func TestFindNodeSuccessful(t *testing.T) {
	node1 := nodeutils.Node{IP: "1", ID: hashing.NewKademliaID("1")}
	node2 := nodeutils.Node{IP: "2", ID: hashing.NewKademliaID("2")}
	node3 := nodeutils.Node{IP: "3", ID: hashing.NewKademliaID("3")}
	expected := &[constants.CLOSESTNODES]nodeutils.Node{node1, node2, node3}

	var nw Network = &MockNetwork{findNode: expected}
	actual, _ := findNode(&nw, &node1, hashing.NewKademliaID("3"))
	if expected != actual {
		t.Errorf("Expected %#v but got %#v", expected, actual)
	}
}

func TestFindNodeFailed(t *testing.T) {
	node := nodeutils.Node{IP: "1", ID: hashing.NewKademliaID("1")}
	var nw Network = &MockNetwork{findNodeErr: errors.New("Random error")}
	_, err := findNode(&nw, &node, hashing.NewKademliaID("3"))

	if err == nil {
		t.Errorf("Expected findNode to return an error but it didn't")
	}
}

func TestFindValueSuccessful(t *testing.T) {
	expected := "def"
	var nw Network = &MockNetwork{findValue: expected}
	key := hashing.NewKademliaID("abc")
	actual := findValue(&nw, key)
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
}

func TestFindValueFailed(t *testing.T) {
	var expected string
	var nw Network = &MockNetwork{findValue: expected}
	key := hashing.NewKademliaID("abc")
	actual := findValue(&nw, key)
	if expected != "" {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
}

func TestJoin(t *testing.T) {
	var nw Network = &MockNetwork{}
	join(&nw)
}
