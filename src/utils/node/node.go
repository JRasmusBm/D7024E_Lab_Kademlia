package node

import (
	"fmt"
	"sort"
    "utils/hashing"
)

// Node definition
// stores the KademliaID, the ip address and the distance
type Node struct {
	ID       *hashing.KademliaID
	Address  string
	distance *hashing.KademliaID
}

// NewNode returns a new instance of a Node
func NewNode(id *hashing.KademliaID, ip string) Node {
	return Node{id, address, nil}
}

// CalcDistance calculates the distance to the target and 
// fills the nodes distance field
func (node *Node) CalcDistance(target *hashing.KademliaID) {
	node.distance = node.ID.CalcDistance(target)
}

// Less returns true if node.distance < otherNode.distance
func (node *Node) Less(otherNode *Node) bool {
	return node.distance.Less(otherNode.distance)
}

// String returns a simple string representation of a Node
func (node *Node) String() string {
	return fmt.Sprintf(`node("%s", "%s")`, node.ID, node.Address)
}

// NodeCandidates definition
// stores an array of Nodes
type NodeCandidates struct {
	nodes []Node
}

// Append an array of Nodes to the NodeCandidates
func (candidates *NodeCandidates) Append(nodes []Node) {
	candidates.nodes = append(candidates.nodes, nodes...)
}

// GetNodes returns the first count number of Nodes
func (candidates *NodeCandidates) GetNodes(count int) []Node {
	return candidates.nodes[:count]
}

// Sort the Nodes in NodeCandidates
func (candidates *NodeCandidates) Sort() {
	sort.Sort(candidates)
}

// Len returns the length of the NodeCandidates
func (candidates *NodeCandidates) Len() int {
	return len(candidates.nodes)
}

// Swap the position of the Nodes at i and j
// WARNING does not check if either i or j is within range
func (candidates *NodeCandidates) Swap(i, j int) {
	candidates.nodes[i], candidates.nodes[j] = candidates.nodes[j], candidates.nodes[i]
}

// Less returns true if the Node at index i is smaller than 
// the Node at index j
func (candidates *NodeCandidates) Less(i, j int) bool {
	return candidates.nodes[i].Less(&candidates.nodes[j])
}
