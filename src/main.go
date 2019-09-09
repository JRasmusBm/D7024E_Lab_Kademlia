package main

import (
	networkutils "utils/network"
	network "network"
	nodeutils "utils/node"
	"fmt"
)

func main() {
	ip, err := networkutils.GetIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip)

	var node *nodeutils.Node
	node = new(nodeutils.Node)
	node.IP = "127.0.0.1"
	network.PingNode(node)
}
