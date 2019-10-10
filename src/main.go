package main

import (
	api_p "api"
	cli "cli/server"
	"fmt"
	"network"
	"utils/hashing"
	networkutils "utils/network"
	nodeutils "utils/node"
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
	api := api_p.API{Sender: sender}

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
  go cli.CliServerInit(api, cliChannel)

	// Busy wait in main thread until "exit" is sent by CLI
	for {
		cliVar := <-cliChannel
		if cliVar == "exit" {
			break
		}
	}
}
