package node

import (
	"testing"
	"utils/hashing"
)

func TestThreeClosest1(t *testing.T) {
	id, _ := hashing.ToKademliaID("FFFFFFFF00000000000000000000000000000000")
	root_node := NewNode(id, "localhost:8000")
	id1, _ := hashing.ToKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2, _ := hashing.ToKademliaID("1111111100000000000000000000000000000000")
	id3, _ := hashing.ToKademliaID("1111111200000000000000000000000000000000")
	id4, _ := hashing.ToKademliaID("1111111300000000000000000000000000000000")
	id5, _ := hashing.ToKademliaID("1111111400000000000000000000000000000000")
	id6, _ := hashing.ToKademliaID("2111111400000000000000000000000000000000")
	nodes := []Node{
		NewNode(id1, "localhost:8001"),
		NewNode(id2, "localhost:8002"),
		NewNode(id3, "localhost:8003"),
		NewNode(id4, "localhost:8004"),
		NewNode(id5, "localhost:8005"),
		NewNode(id6, "localhost:8006"),
	}
	rt := NewRoutingTable(root_node)
	addNodes := make(chan AddNodeOp)
	findClosestNodes := make(chan FindClosestNodesOp, 3)
	go TableSynchronizer(rt, addNodes, findClosestNodes)
	for _, node := range nodes {
		addNodes <- AddNodeOp{AddedNode: node}
	}

	resp := make(chan []Node)
	id7, _ := hashing.ToKademliaID("2111111400000000000000000000000000000000")
	findClosestNodes <- FindClosestNodesOp{
		Target: id7,
		Count:  3,
		Resp:   resp,
	}
	actual := <-resp
	expected := []Node{nodes[5], nodes[4], nodes[3]}
	if len(actual) != 3 ||
		actual[0].ID != expected[0].ID ||
		actual[1].ID != expected[1].ID ||
		actual[2].ID != actual[2].ID {
		t.Errorf("Position 0: Expected %v got %v", expected[0].ID, actual[0].ID)
		t.Errorf("Position 1: Expected %v got %v", expected[1].ID, actual[1].ID)
		t.Errorf("Position 2: Expected %v got %v", expected[2].ID, actual[2].ID)
	}
}

func TestThreeClosest2(t *testing.T) {
	id, _ := hashing.ToKademliaID("FFFFFFFF00000000000000000000000000000000")
	root_node := NewNode(id, "localhost:8000")
	id1, _ := hashing.ToKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2, _ := hashing.ToKademliaID("1111111100000000000000000000000000000000")
	id3, _ := hashing.ToKademliaID("1111111200000000000000000000000000000000")
	id4, _ := hashing.ToKademliaID("1111111300000000000000000000000000000000")
	id5, _ := hashing.ToKademliaID("1111111400000000000000000000000000000000")
	id6, _ := hashing.ToKademliaID("2111111400000000000000000000000000000000")
	nodes := []Node{
		NewNode(id1, "localhost:8001"),
		NewNode(id2, "localhost:8002"),
		NewNode(id3, "localhost:8002"),
		NewNode(id4, "localhost:8002"),
		NewNode(id5, "localhost:8002"),
		NewNode(id6, "localhost:8002"),
	}
	rt := NewRoutingTable(root_node)
	addNodes := make(chan AddNodeOp)
	findClosestNodes := make(chan FindClosestNodesOp, 3)
	go TableSynchronizer(rt, addNodes, findClosestNodes)
	for _, node := range nodes {
		addNodes <- AddNodeOp{AddedNode: node}
	}

	resp := make(chan []Node)
	id7, _ := hashing.ToKademliaID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")
	findClosestNodes <- FindClosestNodesOp{
		Target: id7,
		Count:  3,
		Resp:   resp,
	}
	actual := <-resp
	expected := []Node{nodes[5], nodes[4], nodes[3]}
	if len(actual) != 3 ||
		actual[0].ID != expected[0].ID ||
		actual[1].ID != expected[1].ID ||
		actual[2].ID != actual[2].ID {
		t.Errorf("Position 0: Expected %v got %v", expected[0].ID, actual[0].ID)
		t.Errorf("Position 1: Expected %v got %v", expected[1].ID, actual[1].ID)
		t.Errorf("Position 2: Expected %v got %v", expected[2].ID, actual[2].ID)
	}
}

func TestAllClosest(t *testing.T) {
	id, _ := hashing.ToKademliaID("FFFFFFFF00000000000000000000000000000000")
	root_node := NewNode(id, "localhost:8000")
	id1, _ := hashing.ToKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2, _ := hashing.ToKademliaID("1111111100000000000000000000000000000000")
	id3, _ := hashing.ToKademliaID("1111111200000000000000000000000000000000")
	id4, _ := hashing.ToKademliaID("1111111300000000000000000000000000000000")
	id5, _ := hashing.ToKademliaID("1111111400000000000000000000000000000000")
	id6, _ := hashing.ToKademliaID("2111111400000000000000000000000000000000")
	nodes := []Node{
		NewNode(id1, "localhost:8001"),
		NewNode(id2, "localhost:8002"),
		NewNode(id3, "localhost:8002"),
		NewNode(id4, "localhost:8002"),
		NewNode(id5, "localhost:8002"),
		NewNode(id6, "localhost:8002"),
	}
	rt := NewRoutingTable(root_node)
	addNodes := make(chan AddNodeOp)
	findClosestNodes := make(chan FindClosestNodesOp)
	go TableSynchronizer(rt, addNodes, findClosestNodes)
	for _, node := range nodes {
		addNodes <- AddNodeOp{AddedNode: node}
	}

	resp := make(chan []Node)
	id7, _ := hashing.ToKademliaID("2111111400000000000000000000000000000000")
	findClosestNodes <- FindClosestNodesOp{
		Target: id7,
		Count:  20,
		Resp:   resp,
	}
	actual := <-resp
	expected := []Node{nodes[5], nodes[4], nodes[1], nodes[2], nodes[3], nodes[0]}
	if len(actual) != 6 ||
		actual[0].ID != expected[0].ID ||
		actual[1].ID != actual[1].ID ||
		actual[2].ID != actual[2].ID ||
		actual[3].ID != actual[3].ID ||
		actual[4].ID != actual[4].ID {
		t.Errorf("Position 0: Expected %v got %v", expected[0].ID, actual[0].ID)
		t.Errorf("Position 1: Expected %v got %v", expected[1].ID, actual[1].ID)
		t.Errorf("Position 2: Expected %v got %v", expected[2].ID, actual[2].ID)
		t.Errorf("Position 3: Expected %v got %v", expected[3].ID, actual[3].ID)
		t.Errorf("Position 4: Expected %v got %v", expected[4].ID, actual[4].ID)
		t.Errorf("Position 5: Expected %v got %v", expected[5].ID, actual[5].ID)
	}
}

func TestAddNilNode(t *testing.T) {
	id, _ := hashing.ToKademliaID("FFFFFFFF00000000000000000000000000000000")
	root_node := NewNode(
		id,
		"localhost:8000",
	)
	rt := NewRoutingTable(root_node)
	rt.AddNode(Node{})
}

func TestGetInvalidBucketIndex(t *testing.T) {
	id, _ := hashing.ToKademliaID("FFFFFFFF00000000000000000000000000000000")
	root_node := NewNode(
		id,
		"localhost:8000",
	)
	rt := NewRoutingTable(root_node)
	expected := hashing.IDLength*8 - 1
	id1, _ := hashing.ToKademliaID("0000000000000000000000000000000000000000")
	actual := rt.getBucketIndex(
		id1,
		0,
	)
	if expected != actual {
		t.Errorf("Expected %d got %d", expected, actual)
	}
}
