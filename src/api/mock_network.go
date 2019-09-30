package api

import (
	"utils/constants"
	"utils/hashing"
	nodeutils "utils/node"
)

type MockNetwork struct {
	ping        bool
	pingErr     error
	store       *hashing.KademliaID
	findNode    *[constants.CLOSESTNODES]nodeutils.Node
	findNodeErr error
	findValue   string
}

func (r *MockNetwork) Ping(
	node *nodeutils.Node,
	ch chan bool,
	errCh chan error,
) {
	if r.pingErr != nil {
		errCh <- r.pingErr
	}
	ch <- r.ping
	return
}

func (r *MockNetwork) Store(content string, ch chan *hashing.KademliaID) {
	ch <- r.store
	return
}

func (r *MockNetwork) FindNode(
	node *nodeutils.Node,
	id *hashing.KademliaID,
	ch chan *[constants.CLOSESTNODES]nodeutils.Node,
	errCh chan error,
) {
	if r.findNodeErr != nil {
		errCh <- r.findNodeErr
	}
	ch <- r.findNode
	return
}

func (r *MockNetwork) FindValue(key *hashing.KademliaID, ch chan string) {
	ch <- r.findValue
	return
}

func (r *MockNetwork) Join() {
	return
}
