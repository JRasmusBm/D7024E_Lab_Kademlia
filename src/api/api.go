package api

import (
	"network"
	"utils/constants"
	hashing "utils/hashing"
	nodeutils "utils/node"
)

type API struct {
	Sender network.Sender
}

func (api API) Ping(node *nodeutils.Node) bool {
	ch := make(chan bool)
	errCh := make(chan error)

	go api.Sender.Ping(node, ch, errCh)
	select {
	case ok := <-ch:
		return ok
	case <-errCh:
		return false
	}
}

func (api API) Store(content string) (*hashing.KademliaID, error) {
	ch := make(chan *hashing.KademliaID)
	errCh := make(chan error)
	go api.Sender.Store(content, ch, errCh)
	select {
	case key := <-ch:
		return key, nil
	case err := <-errCh:
		return nil, err
	}
}

func (api API) FindNode(node *nodeutils.Node, id *hashing.KademliaID) (*[constants.CLOSESTNODES]nodeutils.Node, error) {
	ch := make(chan *[constants.CLOSESTNODES]nodeutils.Node)
	errCh := make(chan error)

	go api.Sender.FindNode(node, id, ch, errCh)
	select {
	case nodes := <-ch:
		return nodes, nil
	case err := <-errCh:
		return nil, err
	}
}

func (api API) FindValue(key *hashing.KademliaID) (string, error) {
	ch := make(chan string)
	errCh := make(chan error)
	// TODO: Implement FindValue

	go api.Sender.FindValue(key, ch, errCh)
	select {
	case value := <-ch:
		return value, nil
	case err := <-errCh:
		return "", err
	}
}

func (api API) Join(ip string) bool {
	// TODO: Implement this
	ch := make(chan bool)
	errCh := make(chan error)
	node := nodeutils.Node{IP: "172.19.1.2"}
	go api.Sender.Join(&node, ch, errCh)
	select {
	case ok := <-ch:
		return ok
	case <-errCh:
		return false
	}
}
