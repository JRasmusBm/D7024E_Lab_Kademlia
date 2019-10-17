package api

import (
	"errors"
	"network"
	"strconv"
	"testing"
	"utils/hashing"
	nodeutils "utils/node"
)

func TestPingSuccessful(t *testing.T) {
	var sender network.Sender = &network.MockSender{PingResponse: true}
	api := API{Sender: sender}
	node := nodeutils.Node{IP: ""}
	expected := true
	actual := api.Ping(node)
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			strconv.FormatBool(expected),
			strconv.FormatBool(actual))
	}
}

func TestPingFailed(t *testing.T) {
	var sender network.Sender = &network.MockSender{PingResponse: false}
	api := API{Sender: sender}
	node := nodeutils.Node{IP: ""}
	expected := false
	actual := api.Ping(node)
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			strconv.FormatBool(expected),
			strconv.FormatBool(actual))
	}
}

func TestPingError(t *testing.T) {
	var sender network.Sender = &network.MockSender{PingErr: errors.New("Random error")}
	api := API{Sender: sender}
	node := nodeutils.Node{IP: ""}
	expected := false
	actual := api.Ping(node)
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			strconv.FormatBool(expected),
			strconv.FormatBool(actual))
	}
}

func TestStoreSuccessful(t *testing.T) {
	expected := hashing.NewKademliaID("abc")
	var sender network.Sender = &network.MockSender{StoreSent: 2}
	api := API{Sender: sender}
	content := "abc"
	actual, sent := api.Store(content)
	if !expected.Equals(actual) {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
	if sent != 2 {
		t.Errorf("Expected %v got %v", 2, sent)
	}
}

func TestStoreFailedHash(t *testing.T) {
	expected := hashing.NewKademliaID("abc")
	var sender network.Sender = &network.MockSender{StoreSent: 2}
	api := API{Sender: sender}
	content := "abc"
	actual, sent := api.Store(content)
	if !expected.Equals(actual) {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
	if sent != 2 {
		t.Errorf("Expected %v got %v", 2, sent)
	}
}

func TestFindNodeSuccessful(t *testing.T) {
	id1 := hashing.NewKademliaID("1")
	id2 := hashing.NewKademliaID("2")
	id3 := hashing.NewKademliaID("3")
	node1 := nodeutils.Node{IP: "1", ID: id1}
	node2 := nodeutils.Node{IP: "2", ID: id2}
	node3 := nodeutils.Node{IP: "3", ID: id3}
	expected := []nodeutils.Node{node1, node2, node3}

	var sender network.Sender = &network.MockSender{LookUpResult: expected}
	api := API{Sender: sender}
	actual, _ := api.FindNode(id3)
	if nodeutils.ToStrings(expected) != nodeutils.ToStrings(actual) {
		t.Errorf("Expected %#v but got %#v", nodeutils.ToStrings(expected), nodeutils.ToStrings(actual))
	}
}

func TestFindValueSuccessful(t *testing.T) {
	expected := "def"
	var sender network.Sender = &network.MockSender{LookUpValueResult: expected}
	api := API{Sender: sender}
	key := hashing.NewKademliaID("abc")
	actual, _ := api.FindValue(key)
	if expected != actual {
		t.Errorf(
			"Expected %s got %s",
			expected,
			actual)
	}
}

func TestJoinSuccessful(t *testing.T) {
	var sender network.Sender = &network.MockSender{JoinResponse: true}
	api := API{Sender: sender}
	ok := api.Join("1.2.3.4")
	if !ok {
		t.Errorf("Should be able to join successfully")
	}
}

func TestJoinFailed(t *testing.T) {
	var sender network.Sender = &network.MockSender{JoinResponse: false}
	api := API{Sender: sender}
	ok := api.Join("1.2.3.4")
	if ok {
		t.Errorf("Should be able to join successfully")
	}
}

func TestJoinError(t *testing.T) {
	var sender network.Sender = &network.MockSender{JoinErr: errors.New("Random Error")}
	api := API{Sender: sender}
	ok := api.Join("1.2.3.4")
	if ok {
		t.Errorf("Should be able to join successfully")
	}
}
