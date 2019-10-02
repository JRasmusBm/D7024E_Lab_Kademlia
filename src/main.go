package main

import (
	cli "cli/server"
	"fmt"
	networkutils "utils/network"
	nodeutils "utils/node"
	"utils/hashing"
	"network"
	//"time"
)

func main() {
	ip, err := networkutils.GetIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip)

	node := nodeutils.NewNode(hashing.NewRandomKademliaID(), ip)
	table := nodeutils.NewRoutingTable(node)

	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp)
	sender := network.RealSender{AddNode: addNode, FindClosestNodes: findClosestNodes}

	go nodeutils.TableSynchronizer(table, addNode, findClosestNodes)

	if ip == "172.19.1.2" {
		// TODO: Handle case when bootstrap node
	} else {
		// TODO: Handle case when "normal" node
	}

	// Start node receiver.
	go network.Receiver(ip, sender)

	// Start CLI
	cliChannel := make(chan string)
	go cli.CliServer(cliChannel, sender)

	// Busy wait in main thread until "exit" is sent by CLI
	for {
		cliVar := <- cliChannel
		if cliVar == "exit" {
			break
		}
	}
}
