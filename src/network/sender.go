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
	Ping(node *nodeutils.Node, ch chan bool, errCh chan error)
	Store(content string, ch chan *hashing.KademliaID, errCh chan error)
	FindNode(node *nodeutils.Node, id *hashing.KademliaID, ch chan *[constants.CLOSESTNODES]nodeutils.Node, errCh chan error)
	FindValue(key *hashing.KademliaID, ch chan string, errCh chan error)
	Join(node *nodeutils.Node, ch chan bool, errCh chan error)
}

type RealSender struct {
	AddNode          chan nodeutils.AddNodeOp
	FindClosestNodes chan nodeutils.FindClosestNodesOp
}

func (sender RealSender) Ping(node *nodeutils.Node, ch chan bool, errCh chan error) {
	conn, err := dial(node)
	if err != nil {
		errCh <- err
		return
	}
	fmt.Fprintf(conn, "PING;")
	msg, err := bufio.NewReader(conn).ReadString(';')
	if err != nil {
		errCh <- err
	}

	msg = strings.TrimRight(msg, ";")

	ch <- (msg == "PONG")
}

func (sender RealSender) Store(content string, ch chan *hashing.KademliaID, errCh chan error) {
	hash, _ := hashing.NewKademliaID(content)
	ch <- hash
	return
}

func (sender RealSender) FindNode(node *nodeutils.Node, id *hashing.KademliaID, ch chan *[constants.CLOSESTNODES]nodeutils.Node, errCh chan error) {
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

	nodes, err := nodeutils.FromStrings(msg)
	if err != nil {
		errCh <- err
		return
	}

	// Add all given nodes to routing table
	var result chan bool
	for _, node := range nodes {
		sender.AddNode <- nodeutils.AddNodeOp{AddedNode: node, Resp: result}
	}

	ch <- nodes
}

func (sender RealSender) FindValue(key *hashing.KademliaID, ch chan string, errCh chan error) {
	fmt.Printf("Finding Value %s", key)
	ch <- "Random value"
	return
}

func (sender RealSender) Join(node *nodeutils.Node, ch chan bool, errCh chan error) {
	fmt.Printf("Joining Kademlia")
	conn, err := dial(node)
	if err != nil {
		errCh <- err
	}
	fmt.Fprintf(conn, "JOIN "+node.String()+";")

	msg, err := bufio.NewReader(conn).ReadString(';')
	if err != nil {
		ch <- false
		return
	}

	msg = strings.TrimRight(msg, ";")
	ch <- (msg == "SUCCESS")
}
