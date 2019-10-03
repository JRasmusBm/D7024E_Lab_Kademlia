package network

import (
	"utils/constants"
	"utils/hashing"
	nodeutils "utils/node"
)

type MockSender struct {
	PingResponse      bool
	PingErr           error
	StoreResponse     *hashing.KademliaID
	StoreErr          error
	FindNodeResponse  *[constants.CLOSESTNODES]nodeutils.Node
	FindNodeErr       error
	FindValueResponse string
	FindValueErr      error
	JoinResponse      bool
	JoinErr           error
}

func (r *MockSender) Ping(
	node *nodeutils.Node,
	ch chan bool,
	errCh chan error,
) {
	if r.PingErr != nil {
		errCh <- r.PingErr
	}
	ch <- r.PingResponse
	return
}

func (r *MockSender) Store(
	content string,
	ch chan *hashing.KademliaID,
	errCh chan error,
) {
	if r.StoreErr != nil {
		errCh <- r.StoreErr
	}
	ch <- r.StoreResponse
	return
}

func (r *MockSender) FindNode(
	node *nodeutils.Node,
	id *hashing.KademliaID,
	ch chan *[constants.CLOSESTNODES]nodeutils.Node,
	errCh chan error,
) {
	if r.FindNodeErr != nil {
		errCh <- r.FindNodeErr
	}
	ch <- r.FindNodeResponse
	return
}

func (r *MockSender) FindValue(key *hashing.KademliaID, ch chan string, errCh chan error) {
	if r.FindValueErr != nil {
		errCh <- r.FindValueErr
	}
	ch <- r.FindValueResponse
	return
}

func (r *MockSender) Join(node *nodeutils.Node, ch chan bool, errCh chan error) {
	if r.JoinErr != nil {
		errCh <- r.JoinErr
	}
	ch <- r.JoinResponse
	return
}
