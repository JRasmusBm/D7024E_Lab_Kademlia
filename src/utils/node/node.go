package node

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"sync"
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

// ToStrings returns a string representation of an array of Nodes
func ToStrings(nodes []Node) string {
	nodeStr := ""
	first := true
	for _, node := range nodes {
		if &node != nil {
			if !first {
				nodeStr += " "
			}
			nodeStr += node.String()
			first = false
		}
	}
	return nodeStr
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
func FromStrings(str string) ([]Node, error) {
	nodes_string := strings.Split(str, " ")
	found_nodes := []Node{}
	var node Node
	var err error
	for _, str := range nodes_string {
		node, err = FromString(str)
		if err != nil {
			fmt.Printf("\nTried to parse: %s\nGot Error: %s", str, err.Error())
			continue
		}
		found_nodes = append(found_nodes, node)
	}

	fmt.Printf("%v,\n ", found_nodes)

	return found_nodes, nil
}

// NodeCandidates definition
// stores an array of Nodes
type NodeCandidates struct {
	sync.RWMutex
	Nodes []Node
}

// Append an array of Nodes to the NodeCandidates
func (candidates *NodeCandidates) Append(nodes []Node) {
	candidates.Lock()
  for _, node := range nodes {
    if !NodeInArr(node, candidates.Nodes) {
      candidates.Nodes = append(candidates.Nodes, node)
    }
  }
	// candidates.Nodes = append(candidates.Nodes, nodes...)
	candidates.Unlock()
}

// GetNodePointers returns the first count number of Nodes
func (candidates *NodeCandidates) GetNodePointers(count int) [constants.CLOSESTNODES]*Node {
	candidates.RLock()
	var result [constants.CLOSESTNODES]*Node
	for i, node := range candidates.GetNodes(count) {
		result[i] = &node
	}
	candidates.RUnlock()
	return result
}

// GetNodes returns the first count number of Nodes
func (candidates *NodeCandidates) GetNodes(count int) []Node {
	candidates.RLock()
	var nodes []Node
	if count > len(candidates.Nodes) {
		nodes = candidates.Nodes
	} else {
		nodes = candidates.Nodes[:count]
	}
	candidates.RUnlock()
	return nodes
}

// Sort the Nodes in NodeCandidates
func (candidates *NodeCandidates) Sort() {
	candidates.RLock()
	tempCandidates := NodeCandidates{
		Nodes: candidates.Nodes,
	}
	candidates.RUnlock()

	sort.Sort(&tempCandidates)
	candidates.Lock()
	candidates.Nodes = tempCandidates.Nodes
	candidates.Unlock()

}

// Len returns the length of the NodeCandidates
func (candidates *NodeCandidates) Len() int {
	candidates.RLock()
	length := len(candidates.Nodes)
	candidates.RUnlock()
	return length
}

// Swap the position of the Nodes at i and j
// WARNING does not check if either i or j is within range
func (candidates *NodeCandidates) Swap(i, j int) {
	candidates.Lock()
	candidates.Nodes[i], candidates.Nodes[j] = candidates.Nodes[j], candidates.Nodes[i]
	candidates.Unlock()
}

// Less returns true if the Node at index i is smaller than
// the Node at index j
func (candidates *NodeCandidates) Less(i, j int) bool {
	candidates.RLock()
	isLess := candidates.Nodes[i].Less(&candidates.Nodes[j])
	candidates.RUnlock()
	return isLess
}

func NodeInArr(nodeA Node, nodes []Node) bool {
	for _, nodeB := range nodes {
		if nodeA.ID.String() == nodeB.ID.String() {
			return true
		}
	}
	return false
}
