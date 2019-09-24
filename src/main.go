package main

import (
	cli "cli/server"
	"fmt"
	"time"
	networkutils "utils/network"
	nodeutils "utils/node"
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
	node = new(nodeutils.Node)
	node.IP = "172.20.0.2"

	//while true-loop.
	for {
		cliVar := <- cliChannel
		if cliVar == "exit" {
			break
		}
		time.Sleep(time.Second)
	}
}
