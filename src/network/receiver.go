package network

import (
	"net"
	"strings"
	"bufio"
)

func receiver() {
	ln, _ := net.Listen("tcp", ":6000")
	conn, _ := ln.Accept()

	for {
		// Will block until message ending with newline (\n) is received.
		msg, _ := bufio.NewReader(conn).ReadString('\n')

		// Split string around spaces.
		msg_split := strings.Split(msg, " ")

		switch msg_split[0] {
			case "PING":
				// Syntax: PING
				conn.Write([]byte("PONG"))
			case "FIND_NODE":
				// Syntax: FIND_NODE <id>
				// TODO: Handle FIND_NODE.
				
			case "STORE":
				// Syntax: STORE <key>,<data>
				// TODO
			case "JOIN":
				// Syntax: JOIN 
				// TODO
		}
	}
}
