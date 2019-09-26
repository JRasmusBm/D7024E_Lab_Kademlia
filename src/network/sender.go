package network

import (
	"fmt"
	hashing "utils/hashing"
	nodeutils "utils/node"
	"net"
	"utils/constants"
	"bufio"
	"os"
	"strings"
)

func dial(node *nodeutils.Node) net.Conn {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", node.IP, constants.KADEMLIA_PORT))
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	return conn
}

func Ping(node *nodeutils.Node, ch chan bool) {
	conn := dial(node)
	fmt.Fprintf(conn, "PING;")

	msg, err := bufio.NewReader(conn).ReadString(';')
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
	msg = strings.TrimRight(msg, ";")
	
	ch <- (msg == "PONG")
}

func Store(content string, ch chan *hashing.KademliaID) {
	hash := hashing.NewKademliaID(content)
	ch <- hash
	return
}

func FindNode(node *nodeutils.Node, id *hashing.KademliaID, ch chan *[constants.CLOSESTNODES]nodeutils.Node) {
	fmt.Printf("Finding Node %s", id)
	conn := dial(node)

	fmt.Fprintf(conn, "FIND_NODE %s;", id)
	msg, err := bufio.NewReader(conn).ReadString(';')
	if err != nil {
		fmt.Println("Error", err.Error())
		os.Exit(1)
	}

	msg = strings.TrimRight(msg, ";")

	nodes := nodeutils.FromStrings(msg)

	ch <- nodes
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
