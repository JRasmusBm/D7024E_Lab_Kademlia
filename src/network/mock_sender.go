package network

import (
	"io"
	"utils/constants"
	"utils/hashing"
	nodeutils "utils/node"
)

type MockSender struct {
	PingResponse       bool
	PingErr            error
	StoreSent          int
	FindNodeResponse   []nodeutils.Node
	FindNodeErr        error
	FindValueErr       error
	JoinResponse       bool
	JoinErr            error
	IsFindValueCloser  bool
	IsFindValueSuccess bool
	FindValueCloser    [constants.CLOSESTNODES]*nodeutils.Node
	FindValueSuccess   string
	LookUpResult       [constants.CLOSESTNODES]*nodeutils.Node
	LookUpValueResult  string
}

func (m *MockSender) Dial(node *nodeutils.Node) (io.ReadWriter, error) {
	return nil, nil
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

func (r *MockSender) Store(content string, ch chan int) {
	ch <- r.StoreSent
	return
}

func (r *MockSender) FindNode(
	id *hashing.KademliaID, node *nodeutils.Node, ch chan []nodeutils.Node, errCh chan error) {
	if r.FindNodeErr != nil {
		errCh <- r.FindNodeErr
	}
	ch <- r.FindNodeResponse
	return
}

func (r *MockSender) FindValue(node *nodeutils.Node, key *hashing.KademliaID, successCh chan string, closerCh chan [constants.CLOSESTNODES]*nodeutils.Node, errCh chan error) {
	if r.FindValueErr != nil {
		errCh <- r.FindValueErr
	}
	if r.IsFindValueCloser {
		closerCh <- r.FindValueCloser
	}
	if r.IsFindValueSuccess {
		successCh <- r.FindValueSuccess
	}
	return
}

func (r *MockSender) Join(ip string, ch chan bool, errCh chan error) {
	if r.JoinErr != nil {
		errCh <- r.JoinErr
	}
	ch <- r.JoinResponse
	return
}

func (r *MockSender) LookUp(id *hashing.KademliaID) [constants.CLOSESTNODES]*nodeutils.Node {
	return r.LookUpResult
}

func (r *MockSender) LookUpValue(key *hashing.KademliaID) string {
	return r.LookUpValueResult
}
