package api

import (
	"utils/constants"
	hashing "utils/hashing"
	nodeutils "utils/node"
)

type Network interface {
	Ping(node *nodeutils.Node, sender network.Sender, ch chan bool, errCh chan error)
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

func ping(nw *Network, node *nodeutils.Node, sender network.Sender) bool {
	ch := make(chan bool)
	errCh := make(chan error)

	go sender.Ping(node, ch, errCh)
	select {
	case ok := <-ch:
		return ok
	case <-errCh:
		return false
	}
}

func Ping(node *nodeutils.Node, sender network.Sender) bool {
	var nw Network = &RealNetwork{}
	return ping(&nw, node, sender)
}

func store(nw *Network, content string, sender network.Sender) *hashing.KademliaID {
	ch := make(chan *hashing.KademliaID)
	go sender.Store(content, ch)
	key := <-ch
	return key
}

func Store(content string, sender network.Sender) *hashing.KademliaID {
	var nw Network = &RealNetwork{}
	return store(&nw, content)
}

func findNode(nw *Network, node *nodeutils.Node, id *hashing.KademliaID, sender network.Sender) (*[constants.CLOSESTNODES]nodeutils.Node, error) {
	ch := make(chan *[constants.CLOSESTNODES]nodeutils.Node)
	errCh := make(chan error)

	go sender.FindNode(node, id, ch, errCh)
	select {
		case nodes := <-ch:
			return nodes, nil
		case err := <-errCh:
			return nil, err
		}
}

func FindNode(node *nodeutils.Node, id *hashing.KademliaID, sender network.Sender) (*[constants.CLOSESTNODES]nodeutils.Node, error) {
	var nw Network = &RealNetwork{}
	return findNode(&nw, node, sender)
}

func findValue(nw *Network, key *hashing.KademliaID, sender network.Sender) string {
	ch := make(chan string)
	// TODO: Use channel instead of network
	go sender.FindNode(node, id, ch)
	value := <-ch
	return value
}

func FindValue(key *hashing.KademliaID, sender network.Sender) string {
	var nw Network = &RealNetwork{}
	return findValue(&nw, key, sender)
}

func join(nw *Network, sender network.Sender) {
	// TODO: Implement this
	//sender.Join()
	return
}

func Join(sender network.Sender) {
	var nw Network = &RealNetwork{}
	join(&nw, sender)
	return
}
