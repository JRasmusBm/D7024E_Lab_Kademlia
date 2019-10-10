package network

import (
	"fmt"
	"net"
	"utils/constants"
	hashing "utils/hashing"
	nodeutils "utils/node"
	"encoding/json"
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
	Store(content string, nodes [constants.REPLICATION_FACTOR]*nodeutils.Node, ch chan int)
	FindNode(id *hashing.KademliaID, ch chan [constants.CLOSESTNODES]*nodeutils.Node, errCh chan error)
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

	decoder := json.NewDecoder(conn)

	// Send PING message
	encoder := json.NewEncoder(conn)
	encoder.Encode(Message{RPC: "PING", Msg: PingMsg{Msg: "PING"}})

	// Wait for PONG message
	var msg PingMsg
	err = decoder.Decode(&msg)

	if err != nil {
		errCh <- err
	}

	ch <- (msg.Msg == "PONG")
}

func (sender RealSender) Store(content string, nodes [constants.REPLICATION_FACTOR]*nodeutils.Node, ch chan int) {
	var conn net.Conn
	var err error
	sent := 0
	for _, node := range nodes {
		if node == nil {
			continue
		}

		conn, err = dial(node)
		if err == nil {
			continue
		}

		encoder := json.NewEncoder(conn)
		encoder.Encode(Message{RPC: "STORE", Msg: StoreMsg{Data: content}})
		sent += 1
	}
	ch <- sent
}

func (sender RealSender) FindNode(id *hashing.KademliaID, ch chan [constants.CLOSESTNODES]*nodeutils.Node, errCh chan error) {
	fmt.Printf("Finding Node %s", id)
	
	resp := make(chan []nodeutils.Node)
	sender.FindClosestNodes <- nodeutils.FindClosestNodesOp{Target: id, Count: 1, Resp: resp}
	node := (<- resp)[0]

	conn, err := dial(&node)
	if err != nil {
		errCh <- err
		return
	}
	decoder := json.NewDecoder(conn)

	encoder := json.NewEncoder(conn)
	encoder.Encode(Message{RPC: "FIND_NODE", Msg: FindNodeMsg{ID: id.String()}})

	var msg FindNodeRespMsg
	err = decoder.Decode(&msg)
	if err != nil {
		errCh <- err
		return
	}

	nodes, err := nodeutils.FromStrings(msg.Nodes)
	if err != nil {
		errCh <- err
		return
	}

	// Add all given nodes to routing table
	var result chan bool
	for _, node := range nodes {
		sender.AddNode <- nodeutils.AddNodeOp{AddedNode: *node, Resp: result}
	}

	ch <- nodes
}

func (sender RealSender) FindValue(key *hashing.KademliaID, ch chan string, errCh chan error) {
	fmt.Printf("Finding Value %s", key)
	ch <- "Random value"
	return
}

func (sender RealSender) Join(node *nodeutils.Node, ch chan bool, errCh chan error) {
	conn, err := dial(node)
	if err != nil {
		errCh <- err
	}

	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	encoder.Encode(Message{RPC: "JOIN", Msg: JoinMsg{Msg: node.String()}})

	var msg AckMsg
	err = decoder.Decode(&msg)
	if err != nil {
		ch <- false
		return
	}

	ch <- msg.Success
}
