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
	key := hashing.NewKademliaID(content)

	nodes, _ := api.FindNode(key)

	go api.Sender.Store(content, nodes, ch)
	sent := <-ch

	key, _ := hashing.NewKademliaID(content)

	return key, sent
}

func (api API) FindNode(id *hashing.KademliaID) ([constants.CLOSESTNODES]*nodeutils.Node, error) {
	ch := make(chan [constants.CLOSESTNODES]*nodeutils.Node)
	errCh := make(chan error)

	go api.Sender.FindNode(id, ch, errCh)
	nodes := [constants.CLOSESTNODES]*nodeutils.Node{}
	select {
	case nodes = <-ch:
		return nodes, nil
	case err := <-errCh:
		return nodes, err
	}
}

func (api API) FindValue(key *hashing.KademliaID) (string, error) {
	return api.Sender.LookUpValue(key), nil
}

func (api API) Join(ip string) bool {
	ch := make(chan bool)
	errCh := make(chan error)
	go api.Sender.Join(ip, ch, errCh)
	select {
	case ok := <-ch:
		return ok
	case <-errCh:
		return false
	}
}
