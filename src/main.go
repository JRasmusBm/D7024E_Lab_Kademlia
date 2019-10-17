package main

import (
	api_p "api"
	cli "cli/server"
	"fmt"
	"net"
	"network"
	"utils/constants"
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
	node := nodeutils.NewNode(hashing.NewKademliaID(ip), ip)
	table := nodeutils.NewRoutingTable(node)
	var store storage.Storage = &storage.RealStorage{Data: make(map[string]string)}
	var dialer network.Dialer = &network.RealDialer{}

	listener, _ := net.Listen("tcp", fmt.Sprintf("%s:%d", ip, constants.KADEMLIA_PORT))
	var realListener network.Listener = &network.RealListener{Listener: listener}
	addNode := make(chan nodeutils.AddNodeOp, 1000)
	findClosestNodes := make(chan nodeutils.FindClosestNodesOp, 1000)
	sender := network.RealSender{
		AddNode:          addNode,
		FindClosestNodes: findClosestNodes,
		Storage:          &store,
		Me:               &node,
		Dialer:           &dialer,
	}
	var receiver network.Receiver = &network.RealReceiver{
		Sender:           sender,
		AddNode:          addNode,
		FindClosestNodes: findClosestNodes,
		Storage:          &store,
		Me:               &node,
		IP:               ip,
		Listener:         realListener,
	}
	api := api_p.API{Sender: sender}

	go nodeutils.TableSynchronizer(table, addNode, findClosestNodes)

	// Start node receiver.
	go receiver.Server()

	// Start CLI
	cliChannel := make(chan string)
	go cli.CliServerInit(api, &networkUtils, cliChannel)

	if ip == "172.19.1.2" {
		// TODO: Handle case when bootstrap node
	} else {
		// TODO: Handle case when "normal" node
		go api.Join("172.19.1.2")
	}

	// Busy wait in main thread until "exit" is sent by CLI
	for {
		cliVar := <-cliChannel
		if cliVar == "exit" {
			break
		}
	}
}
