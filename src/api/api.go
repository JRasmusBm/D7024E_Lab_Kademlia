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

func (api API) Store(content string) (*hashing.KademliaID, int) {
	ch := make(chan int)
	key, err := hashing.NewKademliaID(content)
	if err != nil {
		return nil, 0
	}
	
	nodes, err := api.FindNode(key)
	if err != nil {
		return nil, 0
	}

	go api.Sender.Store(content, nodes, ch)
	sent := <- ch

	return key, sent
}

func (api API) FindNode(id *hashing.KademliaID) ([constants.CLOSESTNODES]*nodeutils.Node, error) {
	ch := make(chan [constants.CLOSESTNODES]*nodeutils.Node)
	errCh := make(chan error)

	go api.Sender.FindNode(id, ch, errCh)
	var nodes [constants.CLOSESTNODES]*nodeutils.Node
	select {
	case nodes = <-ch:
		return nodes, nil
	case err := <-errCh:
		return nodes, err
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
