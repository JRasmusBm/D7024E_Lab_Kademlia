package main

import (
	api_p "api"
	cli "cli/server"
	"fmt"
	"network"
	"utils/hashing"
	networkutils "utils/network"
	nodeutils "utils/node"
	"utils/storage"
)

func main() {
  var networkUtils networkutils.NetworkUtils = &networkutils.RealNetworkUtils{}
	ip, err := networkUtils.GetIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip)

	node := nodeutils.NewNode(hashing.NewRandomKademliaID(), ip)
	table := nodeutils.NewRoutingTable(node)
	var store storage.Storage = &storage.RealStorage{Data: make(map[string]string)}
	var dialer network.Dialer = &network.RealDialer{}

	addNode := make(chan nodeutils.AddNodeOp)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp)
	sender := network.RealSender{AddNode: addNode, FindClosestNodes: findClosestNodes, Storage: &store, Me: &node, Dialer: &dialer}
	api := api_p.API{Sender: sender}


	go nodeutils.TableSynchronizer(table, addNode, findClosestNodes)

	if ip == "172.19.1.2" {
		// TODO: Handle case when bootstrap node
	} else {
		// TODO: Handle case when "normal" node
		api.Join("172.19.1.2")
	}

	// Start node receiver.
	go network.Receiver(ip, sender, &store)

	// Start CLI
	cliChannel := make(chan string)
	go cli.CliServerInit(api, &networkUtils, cliChannel)

	// Busy wait in main thread until "exit" is sent by CLI
	for {
		cliVar := <-cliChannel
		if cliVar == "exit" {
			break
		}
	}
}
