package network

import (
    nodeutils "utils/node"
	"os/exec"
)

func PingNode(node *nodeutils.Node) (ret bool) {
	cmd := exec.Command("ping", node.Address, "-c", "3")
	err := cmd.Run()
	return err == nil
}
