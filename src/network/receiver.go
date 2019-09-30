package network

import (
	"net"
	"strings"
	"bufio"
	"utils/constants"
	"fmt"
	"utils/node"
	"utils/hashing"
	"os"
	"strconv"
)

func Receiver(table *node.RoutingTable, ip string) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, constants.KADEMLIA_PORT))
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	fmt.Println("Receiver listening on: " + ip + ":" + strconv.Itoa(constants.KADEMLIA_PORT))

	for {
	// Will block until connection is made.
		conn, _ := ln.Accept()

		// Will block until message ending with semicolon (;) is received.
		msg, _ := bufio.NewReader(conn).ReadString(';')
		msg = strings.TrimRight(msg, ";")
		
		// Split string around spaces.
		msg_split := strings.Split(msg, " ")

		fmt.Println("Message received:", msg)

		switch msg_split[0] {
			case "PING": // Return PONG to verify that the request succeded.
				// Syntax: PING
				conn.Write([]byte("PONG;"))
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
				response += ";"
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
