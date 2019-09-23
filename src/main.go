package main

import (
	cli "cli/server"
	"fmt"
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

	node := nodeutils.NewNode(hashing.NewRandomKademliaID(), ip)
	table := nodeutils.NewRoutingTable(node)

	if ip == "172.19.1.2" {
		// TODO: Handle case when bootstrap node
	} else {
		// TODO: Handle case when "normal" node
	}

	// Receiver will be busy waiting in the main thread.
	network.Receiver(table)
}
