package network

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"utils/constants"
	hashing "utils/hashing"
	nodeutils "utils/node"
)

func dial(node *nodeutils.Node) (net.Conn, error) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", node.IP, constants.KADEMLIA_PORT))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

type Sender interface {
	Ping(node *nodeutils.Node, ch chan bool)
	Store(content string, ch chan *hashing.KademliaID)
	FindNode(node *nodeutils.Node, id *hashing.KademliaID, ch chan *[constants.CLOSESTNODES]nodeutils.Node)
	FindValue(key *hashing.KademliaID, ch chan string)
	Join(node *nodeutils.Node, ch chan bool)
}

type RealSender struct {
	AddNode chan nodeutils.AddNodeOp
	FindClosestNodes chan nodeutils.FindClosestNodesOp
}

func (sender RealSender) Ping(node *nodeutils.Node, ch chan bool) {
	conn := dial(node)
	fmt.Fprintf(conn, "PING;")
	msg, err := bufio.NewReader(conn).ReadString(';')
	if err != nil {
		errCh <- err
		return
	}
	msg = strings.TrimRight(msg, ";")

	ch <- (msg == "PONG")
}

func (sender RealSender) Store(content string, ch chan *hashing.KademliaID) {
	hash := hashing.NewKademliaID(content)
	ch <- hash
	return
}

func (sender RealSender) FindNode(node *nodeutils.Node, id *hashing.KademliaID, ch chan *[constants.CLOSESTNODES]nodeutils.Node) {
	fmt.Printf("Finding Node %s", id)
	conn, err := dial(node)
	if err != nil {
		errCh <- err
		return
	}

	fmt.Fprintf(conn, "FIND_NODE %s;", id)
	msg, err := bufio.NewReader(conn).ReadString(';')
	if err != nil {
		errCh <- err
		return
	}

	msg = strings.TrimRight(msg, ";")

	nodes := nodeutils.FromStrings(msg)

	// Add all given nodes to routing table
	var result chan bool
	for _, node := range nodes {
		sender.AddNode <- nodeutils.AddNodeOp{AddedNode: node, Resp: result}
	}

	ch <- nodes
}

func (sender RealSender) FindValue(key *hashing.KademliaID, ch chan string) {
	fmt.Printf("Finding Value %s", key)
	ch <- "Random value"
	return
}

func (sender RealSender) Join(node *nodeutils.Node, ch chan bool) {
	fmt.Printf("Joining Kademlia")
	conn := dial(node)
	fmt.Fprintf(conn, "JOIN " + node.String() + ";")

	msg, err := bufio.NewReader(conn).ReadString(';')
	if err != nil {
		ch <- false
		return
	}
	
	msg = strings.TrimRight(msg, ";")
	ch <- (msg == "SUCCESS")
}
