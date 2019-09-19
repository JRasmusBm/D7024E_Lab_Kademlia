package network

import (
	"fmt"
	hashing "utils/hashing"
	nodeutils "utils/node"
	"net"
	"utils/constants"
	"bufio"
)

func Ping(node *nodeutils.Node, ch chan bool) {
	conn, _ := net.Dial("tcp", fmt.Sprintf("%s:%d", node.IP, constants.KADEMLIA_PORT))
	fmt.Fprintf(conn, "PING")

	msg, _ := bufio.NewReader(conn).ReadString('\n')

	ch <- (msg == "PONG")
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
