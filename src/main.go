package main

import (
	networkutils "utils/network"
	network "network"
	nodeutils "utils/node"
	"fmt"
	"time"
)

func main() {
	go network.CliApp()

	ip, err := networkutils.GetIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip)

	var node *nodeutils.Node
	node = new(nodeutils.Node)
	node.Address = "172.20.0.2"
	
    if network.PingNode(node) {
        fmt.Println("Succesfully pinged!")
    }

	//while true-loop.
	for {
        time.Sleep(time.Second)
    }
}
