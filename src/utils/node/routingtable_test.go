package node

import (
	"fmt"
	"testing"
)

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewNode(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	rt.AddNode(NewNode(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddNode(NewNode(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddNode(NewNode(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddNode(NewNode(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddNode(NewNode(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002"))
	rt.AddNode(NewNode(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002"))

	contacts := rt.FindClosestNodes(NewKademliaID("2111111400000000000000000000000000000000"), 20)
	for i := range contacts {
		fmt.Println(contacts[i].String())
	}
}