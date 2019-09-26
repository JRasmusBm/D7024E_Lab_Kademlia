package network

import (
	"fmt"
	"os/exec"
	hashing "utils/hashing"
	nodeutils "utils/node"
)

func Ping(node *nodeutils.Node, ch chan bool) (ret bool) {
	cmd := exec.Command("ping", node.IP, "-c", "3")
	err := cmd.Run()
	ch <- err == nil
	return
}

func Store(content string, ch chan *hashing.KademliaID) {
	hash := hashing.NewKademliaID(content)
	ch <- hash
	return
}

func FindNode(id *hashing.KademliaID, ch chan *nodeutils.Node) {
	fmt.Printf("Finding Node %s", id)
	node := nodeutils.Node{IP: "0.0.0.0"}
	ch <- &node
	return
}

func FindValue(key *hashing.KademliaID, ch chan string) {
	fmt.Printf("Finding Value %s", key)
	ch <- "Random value"
	return
}

func Join() {
	fmt.Printf("Joining Kademlia")
	return
}
