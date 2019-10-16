package network

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"utils/constants"
	"utils/hashing"
	nodeutils "utils/node"
	"utils/storage"
)

type Listener interface {
	Accept() (conn io.ReadWriter, err error)
}

type RealListener struct {
	Listener net.Listener
}

func (r *RealListener) Accept() (io.ReadWriter, error) {
	var conn io.ReadWriter
	var err error
	conn, err = r.Listener.Accept()
	return conn, err
}

type Receiver interface {
	Server()
	PingReply(msg Message, rw io.ReadWriter)
	FindNodeReply(msg Message, rw io.ReadWriter)
	FindValueReply(msg Message, rw io.ReadWriter)
	StoreReply(msg Message, rw io.ReadWriter)
	JoinReply(msg Message, rw io.ReadWriter)
}

type RealReceiver struct {
	Me               *nodeutils.Node
	Listener         Listener
	IP               string
	Sender           Sender
	Storage          *storage.Storage
	AddNode          chan nodeutils.AddNodeOp
	FindClosestNodes chan nodeutils.FindClosestNodesOp
}

func (receiver RealReceiver) Server() {
	fmt.Println("Receiver listening on: " + receiver.IP + ":" + strconv.Itoa(constants.KADEMLIA_PORT))

	for {
		// Will block until connection is made.
		conn, err := receiver.Listener.Accept()
		if err != nil {
			fmt.Printf(err.Error())
			continue
		}
		go receiver.handleRequest(conn)
	}
}

func (receiver RealReceiver) handleRequest(conn io.ReadWriter) {
	var tmp struct {
		RPC    string
		Author string
		Msg    json.RawMessage
	}

	decoder := json.NewDecoder(conn)
	err := decoder.Decode(&tmp)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	fmt.Println("Message received:", tmp)

	switch tmp.RPC {
	case "PING": // Return PONG to verify that the request succeded.
		var pingMsg PingMsg
		json.Unmarshal(tmp.Msg, &pingMsg)
		go receiver.PingReply(
			Message{
				Author: tmp.Author,
				RPC:    tmp.RPC,
				Msg:    pingMsg,
			},
			conn,
		)

	case "FIND_NODE": // Return the x closest known nodes in a sequence separated by spaces.
		var findNodeMsg FindNodeMsg
		json.Unmarshal(tmp.Msg, &findNodeMsg)
		go receiver.FindNodeReply(
			Message{
				Author: tmp.Author,
				RPC:    tmp.RPC,
				Msg:    findNodeMsg,
			},
			conn)

	case "STORE": // Store the data and acknowledge.
		var storeMsg StoreMsg
		json.Unmarshal(tmp.Msg, &storeMsg)
		go receiver.StoreReply(
			Message{
				Author: tmp.Author,
				RPC:    tmp.RPC,
				Msg:    storeMsg,
			},
			conn,
		)

	case "JOIN": // Add a node to routing table/bucket list (if possible) and acknowledge.
		var joinMsg JoinMsg
		json.Unmarshal(tmp.Msg, &joinMsg)
		go receiver.JoinReply(
			Message{
				Author: tmp.Author,
				RPC:    tmp.RPC,
				Msg:    joinMsg,
			},
			conn,
		)

	case "FIND_VALUE": // Given a kademlia id for data, return data or return x closest nodes.
		var findValueMsg FindValueMsg
		json.Unmarshal(tmp.Msg, &findValueMsg)
		go receiver.FindValueReply(
			Message{
				Author: tmp.Author,
				RPC:    tmp.RPC,
				Msg:    findValueMsg,
			},
			conn,
		)
	}
}

func (receiver RealReceiver) PingReply(msg Message, conn io.ReadWriter) {
	encoder := json.NewEncoder(conn)
	encoder.Encode(PingMsg{Msg: "PONG"})
}

func (receiver RealReceiver) FindNodeReply(msg Message, conn io.ReadWriter) {
	fmt.Printf("%v", msg.Msg)
	fmt.Printf("%#v", msg.Msg)
	findNodeMsg := msg.Msg.(FindNodeMsg)
	target, _ := hashing.ToKademliaID(findNodeMsg.ID)
	var closest_nodes_ch chan []nodeutils.Node
	receiver.FindClosestNodes <- nodeutils.FindClosestNodesOp{Target: target, Count: constants.CLOSESTNODES}
	closest_nodes := <-closest_nodes_ch
	response := ""
	for i, node := range closest_nodes {
		receiver.AddNode <- nodeutils.AddNodeOp{AddedNode: node}

		if i != len(closest_nodes)-1 {
			response += node.String() + "  "
		} else {
			response += node.String()
		}
	}

	encoder := json.NewEncoder(conn)
	encoder.Encode(FindNodeRespMsg{Nodes: response})
}

func (receiver RealReceiver) FindValueReply(msg Message, conn io.ReadWriter) {
	findValueMsg := msg.Msg.(FindValueMsg)
	key, _ := hashing.ToKademliaID(findValueMsg.Key)

	ch := make(chan string)
	errCh := make(chan error)
	(*receiver.Storage).Read(key.String(), ch, errCh)

	encoder := json.NewEncoder(conn)
	select {
	case content := <-ch:
		encoder.Encode(FindValueRespMsg{Content: content, Nodes: ""})
	case <-errCh:
		resp := make(chan []*nodeutils.Node)
		receiver.FindClosestNodes <- nodeutils.FindClosestNodesOp{
			Target: key,
			Count:  constants.CLOSESTNODES,
			Resp:   resp,
		}
		closestNodes := <-resp
		encoder.Encode(FindValueRespMsg{
			Content: "",
			Nodes:   nodeutils.ToStrings(closestNodes)},
		)
	}
}

func (receiver RealReceiver) StoreReply(msg Message, conn io.ReadWriter) {
	storeMsg := msg.Msg.(StoreMsg)
	kid := hashing.NewKademliaID(storeMsg.Data)
	key := kid.String()
	(*receiver.Storage).Write(key, storeMsg.Data)

	author, err := nodeutils.FromString(msg.Author)
	encoder := json.NewEncoder(conn)
	if err != nil {
		receiver.AddNode <- nodeutils.AddNodeOp{AddedNode: author}
		encoder.Encode(AckMsg{Success: true})
		return
	}
	encoder.Encode(AckMsg{Success: false})
}

func (receiver RealReceiver) JoinReply(msg Message, conn io.ReadWriter) {
	joinMsg := msg.Msg.(JoinMsg)
	node, _ := nodeutils.FromString(joinMsg.Msg)

	receiver.AddNode <- nodeutils.AddNodeOp{AddedNode: node}

	encoder := json.NewEncoder(conn)
	encoder.Encode(JoinRespMsg{Success: true, ID: receiver.Me.ID.String(), IP: receiver.Me.IP})
}
