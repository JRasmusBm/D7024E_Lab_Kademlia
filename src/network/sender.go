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

func Ping(node *nodeutils.Node, ch chan bool, errCh chan error) {
	conn, err := dial(node)
	if err != nil {
		errCh <- err
		return
	}
	fmt.Fprintf(conn, "PING;")
	msg, err := bufio.NewReader(conn).ReadString(';')
	if err != nil {
		errCh <- err
		return
	}
	msg = strings.TrimRight(msg, ";")

	ch <- (msg == "PONG")
}

func Store(content string, ch chan *hashing.KademliaID) {
	hash := hashing.NewKademliaID(content)
	ch <- hash
	return
}

func FindNode(
	node *nodeutils.Node,
	id *hashing.KademliaID,
	ch chan *[constants.CLOSESTNODES]nodeutils.Node,
	errCh chan error,
) {
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

	ch <- nodes
}

func FindValue(key *hashing.KademliaID, ch chan string) {
	fmt.Printf("Finding Value %s", key)
	ch <- "Random value"
	return
}

func Join(node *nodeutils.Node) {
	fmt.Printf("Joining Kademlia")
	conn := dial(node)
	fmt.Fprintf(conn, "JOIN " + node.String())

	bufio.NewReader(conn).ReadString(';')
	return
}
