package network

import (
	"net"
	"strings"
	"bufio"
	"utils/constants"
	"fmt"
	"utils/node"
	"utils/hashing"
)

func Receiver(table *node.RoutingTable) {
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", constants.KADEMLIA_PORT))

	for {
		// Will block until connection is made.
		conn, _ := ln.Accept()

		// Will block until message ending with newline (\n) is received.
		msg, _ := bufio.NewReader(conn).ReadString('\n')

		// Split string around spaces.
		msg_split := strings.Split(msg, " ")

		switch msg_split[0] {
			case "PING": // Return PONG to verify that the request succeded.
				// Syntax: PING
				conn.Write([]byte("PONG"))
			case "FIND_NODE": // Return the x closest known nodes in a sequence separated by spaces.
				// Syntax: FIND_NODE <id>
				target := hashing.NewKademliaID(msg_split[1])
				closest_nodes := table.FindClosestNodes(target, constants.CLOSESTNODES)
				response := ""
				for i, node := range closest_nodes {
					if i != len(closest_nodes)-1 {
						response += node.String() + "  "
					} else {
						response += node.String()
					}
				}
				conn.Write([]byte(response))
			case "STORE": // Store the data and acknowledge.
				// Syntax: STORE <data>
				// TODO
			case "JOIN": // Add a node to routing table/bucket list (if possible) and acknowledge.
				// Syntax: JOIN 
				// TODO
		}
	}
}
