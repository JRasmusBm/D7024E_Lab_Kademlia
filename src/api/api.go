package api

import (
	"network"
	hashing "utils/hashing"
	nodeutils "utils/node"
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

func FindNode(id *hashing.KademliaID) (ret *nodeutils.Node) {
	ch := make(chan *nodeutils.Node)
	go network.FindNode(id, ch)
	node := <-ch
	return node
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
