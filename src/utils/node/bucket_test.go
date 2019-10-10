package node

import (
	"container/list"
	"testing"
	"utils/hashing"
)

func TestAddTwoNodes(t *testing.T) {
	id1 := hashing.NewKademliaID("abc")
	node1 := Node{ID: id1}
	id2 := hashing.NewKademliaID("def")
	node2 := Node{ID: id2}
	theList := list.List{}
	theBucket := bucket{list: &theList}
	theBucket.AddNode(node1)
	theBucket.AddNode(node2)
	actual := theBucket.list.Front().Value.(Node).ID
	if actual != node2.ID {
		t.Errorf(
			"Expected the second node to be in front\nfront: %v\nnode2: %v",
			actual,
			node2.ID,
		)
	}
	if theBucket.Len() != 2 {
		t.Errorf(
			"Expected there to only be two values in the bucket, %v",
			theBucket.list,
		)
	}
}

func TestAddNodeTwice(t *testing.T) {
	id1 := hashing.NewKademliaID("abc")
	id2 := hashing.NewKademliaID("def")
	node1 := Node{ID: id1}
	node2 := Node{ID: id2}
	theList := list.List{}
	theBucket := bucket{list: &theList}
	theBucket.AddNode(node1)
	theBucket.AddNode(node2)
	theBucket.AddNode(node1)
	actual := theBucket.list.Front().Value.(Node).ID
	if actual != node1.ID {
		t.Errorf(
			"Expected the re-added node to be in front\nfront: %v\nnode1: %v",
			actual,
			node1.ID,
		)
	}
	if theBucket.Len() != 2 {
		t.Errorf(
			"Expected there to only be two values in the bucket, %v",
			theBucket.list,
		)
	}
}
