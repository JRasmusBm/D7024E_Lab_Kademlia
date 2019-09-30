package api

import (
	"utils/constants"
	hashing "utils/hashing"
	nodeutils "utils/node"
)

type Network interface {
	Ping(node *nodeutils.Node, ch chan bool, errCh chan error)
	Store(content string, ch chan *hashing.KademliaID)
	FindNode(
		node *nodeutils.Node,
		id *hashing.KademliaID,
		ch chan *[constants.CLOSESTNODES]nodeutils.Node,
		errCh chan error,
	)
	FindValue(key *hashing.KademliaID, ch chan string)
	Join()
}

func ping(nw *Network, node *nodeutils.Node) bool {
	ch := make(chan bool)
	errCh := make(chan error)
	go (*nw).Ping(node, ch, errCh)
	select {
	case ok := <-ch:
		return ok
	case <-errCh:
		return false
	}
}

func Ping(node *nodeutils.Node) bool {
	var nw Network = &RealNetwork{}
	return ping(&nw, node)
}

func store(nw *Network, content string) *hashing.KademliaID {
	ch := make(chan *hashing.KademliaID)
	go (*nw).Store(content, ch)
	key := <-ch
	return key
}

func Store(content string) *hashing.KademliaID {
	var nw Network = &RealNetwork{}
	return store(&nw, content)
}

func findNode(nw *Network, node *nodeutils.Node, id *hashing.KademliaID) (*[constants.CLOSESTNODES]nodeutils.Node, error) {
	ch := make(chan *[constants.CLOSESTNODES]nodeutils.Node)
	errCh := make(chan error)
	go (*nw).FindNode(node, id, ch, errCh)
	select {
	case nodes := <-ch:
		return nodes, nil
	case err := <-errCh:
		return nil, err
	}
}

func FindNode(node *nodeutils.Node, id *hashing.KademliaID) (*[constants.CLOSESTNODES]nodeutils.Node, error) {
	var nw Network = &RealNetwork{}
	return findNode(&nw, node, id)
}

func findValue(nw *Network, key *hashing.KademliaID) string {
	ch := make(chan string)
	go (*nw).FindValue(key, ch)
	value := <-ch
	return value
}

func FindValue(key *hashing.KademliaID) string {
	var nw Network = &RealNetwork{}
	return findValue(&nw, key)
}

func join(nw *Network) {
	go (*nw).Join()
	return
}

func Join() {
	var nw Network = &RealNetwork{}
	join(&nw)
	return
}
