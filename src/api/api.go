package api

import (
	"utils/constants"
	hashing "utils/hashing"
	nodeutils "utils/node"
    "network"
)

type API struct {
    Sender network.Sender
}

func (api API) ping(node *nodeutils.Node) bool {
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

func (api API) Ping(node *nodeutils.Node) bool {
	return api.ping(node)
}

func (api API) store(content string) *hashing.KademliaID {
	ch := make(chan *hashing.KademliaID)
	go api.Sender.Store(content, ch)
	key := <-ch
	return key
}

func (api API) Store(content string) *hashing.KademliaID {
	return api.store(content)
}

func (api API) findNode(node *nodeutils.Node, id *hashing.KademliaID) (*[constants.CLOSESTNODES]nodeutils.Node, error) {
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

func (api API) FindNode(node *nodeutils.Node, id *hashing.KademliaID) (*[constants.CLOSESTNODES]nodeutils.Node, error) {
	return api.findNode(node, id)
}

func (api API) findValue(key *hashing.KademliaID) (string, error) {
	ch := make(chan string)
    errCh := make(chan error)
    var err error
    // TODO: Implement FindValue
	//go api.sender.FindNode(node, id, ch, errCh)
    select {
        case value := <- ch:
            return value, err
        case err = <- errCh:
            return "", err
    }
}

func (api API) FindValue(key *hashing.KademliaID) (string, error) {
	return api.findValue(key)
}

func (api API) join(ip string) {
	// TODO: Implement this
	//sender.Join()
	return
}

func (api API) Join(ip string) {
	api.join(ip)
	return
}
