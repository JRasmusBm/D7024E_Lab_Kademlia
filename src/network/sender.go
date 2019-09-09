package network

import (
	"github.com/sparrc/go-ping"
    nodeutils "utils/node"
)

func PingNode(node *nodeutils.Node) {
	pinger, err := ping.NewPinger(node.IP)
	if err != nil {
		panic(err)
	}
	pinger.Count = 3
	pinger.Run()
}
