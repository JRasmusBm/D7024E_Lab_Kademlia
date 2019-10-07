package network

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"utils/constants"
	"utils/hashing"
	nodeutils "utils/node"
	"encoding/json"
)

func Receiver(ip string, sender RealSender) {
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
			var findNodeMsg FindNodeMsg
			findNodeMsg = msg.Msg
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
			// Syntax: STORE <data>;

		case "JOIN": // Add a node to routing table/bucket list (if possible) and acknowledge.
			// Syntax: JOIN <node>;
			node, err := nodeutils.FromString(msg_split[1])
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
