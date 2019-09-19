package main

import (
	cli "cli/server"
	"fmt"
	"time"
	networkutils "utils/network"
	nodeutils "utils/node"
	"utils/hashing"
	"network"
)

func main() {
	cliChannel := make(chan string)
	go cli.CliServer(cliChannel)

	ip, err := networkutils.GetIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip)

	var node *nodeutils.Node
	node = nodeutils.NewNode(hashing.NewRandomKademliaID(), networkutils.GetIP())

	var table nodeutils.RoutingTable
	table = nodeutils.NewRoutingTable(node)

	// Receiver will be busy waiting in the main thread.
	network.receiver(table)
}
