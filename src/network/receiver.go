package network

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"utils/constants"
	"utils/hashing"
	nodeutils "utils/node"
	"utils/storage"
)

func Receiver(ip string, sender RealSender, store *storage.Storage) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, constants.KADEMLIA_PORT))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Receiver listening on: " + ip + ":" + strconv.Itoa(constants.KADEMLIA_PORT))

	for {
		// Will block until connection is made.
		conn, _ := ln.Accept()

		decoder := json.NewDecoder(conn)
		var msg Message
		err := decoder.Decode(&msg)
		if err != nil {
			log.Fatal(err)
		}
		author, _ := nodeutils.FromString(msg.Author)

		encoder := json.NewEncoder(conn)

		fmt.Println("Message received:", msg)

		switch msg.RPC {
		case "PING": // Return PONG to verify that the request succeded.
			encoder.Encode(PingMsg{Msg: "PONG"})

		case "FIND_NODE": // Return the x closest known nodes in a sequence separated by spaces.
			findNodeMsg := msg.Msg.(FindNodeMsg)
			target, _ := hashing.ToKademliaID(findNodeMsg.ID)
			var closest_nodes_ch chan []nodeutils.Node
			sender.FindClosestNodes <- nodeutils.FindClosestNodesOp{Target: target, Count: constants.CLOSESTNODES}
			closest_nodes := <-closest_nodes_ch
			response := ""
			for i, node := range closest_nodes {
				sender.AddNode <- nodeutils.AddNodeOp{AddedNode: node}

				if i != len(closest_nodes)-1 {
					response += node.String() + "  "
				} else {
					response += node.String()
				}
			}
			encoder.Encode(FindNodeRespMsg{Nodes: response})

		case "STORE": // Store the data and acknowledge.
			storeMsg := msg.Msg.(StoreMsg)
			kid := hashing.NewKademliaID(storeMsg.Data)
			key := kid.String()
			(*store).Write(key, storeMsg.Data)

			sender.AddNode <- nodeutils.AddNodeOp{AddedNode: author}

			encoder.Encode(AckMsg{Success: true})

		case "JOIN": // Add a node to routing table/bucket list (if possible) and acknowledge.
			joinMsg := msg.Msg.(JoinMsg)
			node, err := nodeutils.FromString(joinMsg.Msg)

			result := make(chan bool)
			sender.AddNode <- nodeutils.AddNodeOp{AddedNode: node}

			var success bool
			if <-result && err == nil {
				success = true
			} else {
				success = false
			}

			encoder.Encode(JoinRespMsg{Success: success, ID: sender.Me.ID.String()})

		case "FIND_VALUE": // Given a kademlia id for data, return data or return x closest nodes.
			findValueMsg := msg.Msg.(FindValueMsg)
			key, _ := hashing.ToKademliaID(findValueMsg.Key)

			ch := make(chan string)
			errCh := make(chan error)
			(*store).Read(key.String(), ch, errCh)

			select {
			case content := <-ch:
				encoder.Encode(FindValueRespMsg{Content: content, Nodes: ""})
			case <-errCh:
				resp := make(chan []*nodeutils.Node)
				sender.FindClosestNodes <- nodeutils.FindClosestNodesOp{
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
	}
}
