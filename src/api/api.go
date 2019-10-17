package api

import (
	"network"
	hashing "utils/hashing"
	nodeutils "utils/node"
)

type API struct {
	Sender network.Sender
}

func (api API) Ping(node nodeutils.Node) bool {
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

	go api.Sender.Store(content, ch)
	sent := <-ch

	return key, sent
}

func (api API) FindNode(id *hashing.KademliaID) ([]nodeutils.Node, error) {
	return api.Sender.LookUp(id), nil
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
