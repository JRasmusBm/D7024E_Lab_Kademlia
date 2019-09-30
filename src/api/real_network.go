package api

import (
	"network"
	"utils/constants"
	"utils/hashing"
	nodeutils "utils/node"
)

type RealNetwork struct{}

func (r *RealNetwork) Ping(node *nodeutils.Node, ch chan bool, errCh chan error) {
	network.Ping(node, ch, errCh)
	return
}

func (r *RealNetwork) Store(content string, ch chan *hashing.KademliaID) {
	network.Store(content, ch)
	return
}

func (r *RealNetwork) FindNode(
	node *nodeutils.Node,
	id *hashing.KademliaID,
	ch chan *[constants.CLOSESTNODES]nodeutils.Node,
	errCh chan error,
) {
	network.FindNode(node, id, ch, errCh)
	return
}

func (r *RealNetwork) FindValue(key *hashing.KademliaID, ch chan string) {
	network.FindValue(key, ch)
	return
}

func (r *RealNetwork) Join() {
	network.Join()
	return
}
