package node

import (
	"utils/hashing"
)

const bucketSize = 20

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type RoutingTable struct {
	me      Node
	buckets [hashing.IDLength * 8]*bucket
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me Node) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < hashing.IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()
	}
	routingTable.me = me
	return routingTable
}

// AddNode add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddNode(contact Node) {
  if contact.ID == nil {
    return
  }
	bucketIndex := routingTable.getBucketIndex(contact.ID, hashing.IDLength)
	bucket := routingTable.buckets[bucketIndex]
	bucket.AddNode(contact)
}

// FindClosestNodes finds the count closest Nodes to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestNodes(target *hashing.KademliaID, count int) []Node {
	var candidates NodeCandidates
	bucketIndex := routingTable.getBucketIndex(target, hashing.IDLength)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetNodeAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < hashing.IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetNodeAndCalcDistance(target))
		}
		if bucketIndex+i < hashing.IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetNodeAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetNodes(count)
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) getBucketIndex(id *hashing.KademliaID, idLength int) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < idLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return hashing.IDLength*8 - 1
}

type AddNodeOp struct {
	AddedNode Node
}

type FindClosestNodesOp struct {
	Target *hashing.KademliaID
	Count  int
	Resp   chan []*Node
}

func TableSynchronizer(rTable *RoutingTable, addNodes chan AddNodeOp, findClosestNodes chan FindClosestNodesOp) {
	for {
		select {
		case addNode := <-addNodes:
			rTable.AddNode(addNode.AddedNode)
		case findClosestNode := <-findClosestNodes:
			result := make([]*Node, 0)
			for _, node := range rTable.FindClosestNodes(findClosestNode.Target, findClosestNode.Count) {
				result = append(result, &node)
			}
			findClosestNode.Resp <- result
		}
	}
}
