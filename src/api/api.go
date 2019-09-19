package api

import (
	"fmt"
	"network"
	hashing "utils/hashing"
	nodeutils "utils/node"
)

func Ping(node *nodeutils.Node) (ret bool) {
	ch := make(chan bool)
	go network.Ping(node, ch)
	ok := <-ch
	fmt.Println("Ping successful: ", ok)
	return ok
}

func Store(content string) (ret bool) {
	ch := make(chan bool)
	go network.Store(content, ch)
	ok := <-ch
	fmt.Println("Store Success: ", ok)
	return ok
}

func FindNode(id *hashing.KademliaID) (ret *nodeutils.Node) {
	ch := make(chan *nodeutils.Node)
	go network.FindNode(id, ch)
	node := <-ch
	fmt.Println("Found Node: ", node.IP)
	return node
}

func FindValue(key *hashing.KademliaID) (ret string) {
	ch := make(chan string)
	go network.FindValue(key, ch)
	value := <-ch
	fmt.Println("Found Value: ", value)
	return value
}

func Join() {
	go network.Join()
	return
}
