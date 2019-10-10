package node

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"utils/constants"
	"utils/hashing"
)

// Node definition
// stores the KademliaID, the ip address and the distance
type Node struct {
	ID       *hashing.KademliaID
	IP       string
	distance *hashing.KademliaID
}

// NewNode returns a new instance of a Node
func NewNode(id *hashing.KademliaID, ip string) Node {
	return Node{id, ip, nil}
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
	return fmt.Sprintf(`node("%s","%s")`, node.ID, node.IP)
}

// Returns a Node instance from a string representation
func FromString(str string) (Node, error) {
	// TODO: Properly ensure the string is correctly formatted and test this.
	var node Node

	if !strings.HasPrefix(str, "node(\"") || !strings.HasSuffix(str, "\")") {
		return node, fmt.Errorf("Not a valid string interpretation of node")
	}

	// Remove prefix and suffix
	str = strings.Split(str, "node(\"")[1]
	str = strings.Split(str, "\")")[0]

	// Split string into id and ip
	node_data := strings.Split(str, "\",\"")
	if len(node_data) < 2 {
		return Node{}, errors.New(
			"Node string should have the form `node(\"<id>\",\"<ip>\")`",
		)
	}
	idString := node_data[0]
	ip := node_data[1]

	id, err := hashing.ToKademliaID(idString)
	if err != nil {
		return node, err
	}

	node = NewNode(id, ip)

	return node, nil
}

// Returns an array of nodes from string representation of
// nodes separated by spaces.
func FromStrings(str string) ([constants.CLOSESTNODES]*Node, error) {
	nodes_string := strings.Split(str, " ")
	found_nodes := [constants.CLOSESTNODES]*Node{}
	var node Node
	var err error
	for i, str := range nodes_string {
		node, err = FromString(str)
		if err != nil {
			return found_nodes, err
		}
		found_nodes[i] = &node
	}

	return found_nodes, nil
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
