package main

import (
	cli "cli"
	"fmt"
	"time"
	networkutils "utils/network"
	nodeutils "utils/node"
)

func main() {
	go cli.CliApp()

	ip, err := networkutils.GetIP()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ip)

	var node *nodeutils.Node
	node = new(nodeutils.Node)
	node.IP = "172.20.0.2"

	//while true-loop.
	for {
		time.Sleep(time.Second)
	}
}
