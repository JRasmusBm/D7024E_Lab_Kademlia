package api

import (
	"network"
	hashing "utils/hashing"
	nodeutils "utils/node"
	"utils/constants"
)

//func Ping(node *nodeutils.Node) (ret bool) {
//	ch := make(chan bool)
//	go network.Ping(node, ch)
//	ok := <-ch
//	fmt.Println("Ping successful: ", ok)
//	return ok
//}

func Store(content string) (ret *hashing.KademliaID) {
	ch := make(chan *hashing.KademliaID)
	go network.Store(content, ch)
	key := <-ch
	return key
}

func FindNode(node *nodeutils.Node, id *hashing.KademliaID) (ret *[constants.CLOSESTNODES]nodeutils.Node) {
	ch := make(chan *[constants.CLOSESTNODES]nodeutils.Node)
	go network.FindNode(node, id, ch)
	nodes := <-ch
	fmt.Println("Found closest nodes.")
	return nodes
}

func FindValue(key *hashing.KademliaID) (ret string) {
	ch := make(chan string)
	go network.FindValue(key, ch)
	value := <-ch
	return value
}

func Join() {
	go network.Join()
	return
}
