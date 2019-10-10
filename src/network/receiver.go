package network

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"utils/constants"
	"utils/hashing"
	nodeutils "utils/node"
	"encoding/json"
	"log"
	"utils/storage"
)

func Receiver(ip string, sender RealSender, store storage.Storage) {
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

		encoder := json.NewEncoder(conn)

		fmt.Println("Message received:", msg)

		switch msg.RPC {
		case "PING": // Return PONG to verify that the request succeded.
			encoder.Encode(PingMsg{Msg: "PONG"})

		case "FIND_NODE": // Return the x closest known nodes in a sequence separated by spaces.
			findNodeMsg := msg.Msg.(FindNodeMsg)
			target, _ := hashing.NewKademliaID(findNodeMsg.ID)
			var closest_nodes_ch chan []nodeutils.Node
			sender.FindClosestNodes <- nodeutils.FindClosestNodesOp{Target: target, Count: constants.CLOSESTNODES, Resp: closest_nodes_ch}
			closest_nodes := <-closest_nodes_ch
			response := ""
			for i, node := range closest_nodes {
				var result chan bool
				sender.AddNode <- nodeutils.AddNodeOp{AddedNode: node, Resp: result}

				if i != len(closest_nodes)-1 {
					response += node.String() + "  "
				} else {
					response += node.String()
				}
			}
			encoder.Encode(FindNodeRespMsg{Nodes: response})

		case "STORE": // Store the data and acknowledge.
			storeMsg := msg.Msg.(StoreMsg)
			kid, _ := hashing.NewKademliaID(storeMsg.Data)
			key := kid.String()
			store.Write(key, storeMsg.Data)
			encoder.Encode(AckMsg{Success: true})

		case "JOIN": // Add a node to routing table/bucket list (if possible) and acknowledge.
			joinMsg := msg.Msg.(JoinMsg)
			node, err := nodeutils.FromString(joinMsg.Msg)
			var result chan bool
			sender.AddNode <- nodeutils.AddNodeOp{AddedNode: node, Resp: result}

			var response string
			if <-result && err == nil {
				response = "SUCCESS;"
			} else {
				response = "FULL;"
			}

			conn.Write([]byte(response))
		}
	}
}
