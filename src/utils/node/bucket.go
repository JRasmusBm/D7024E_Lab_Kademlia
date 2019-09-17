package node

import (
	"container/list"
    hashing "utils/hashing"
)

// bucket definition
// contains a List
type bucket struct {
	list *list.List
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	return bucket
}

// AddNode adds the Node to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddNode(node Node) {
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Node).ID

		if (node).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(node)
		}
	} else {
		bucket.list.MoveToFront(element)
	}
}

// GetNodeAndCalcDistance returns an array of Nodes where 
// the distance has already been calculated
func (bucket *bucket) GetNodeAndCalcDistance(target *hashing.KademliaID) []Node {
	var nodes []Node

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		node := elt.Value.(Node)
		node.CalcDistance(target)
		nodes = append(nodes, node)
	}

	return nodes
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
