package node

import (
	"fmt"
	"testing"
	"utils/constants"
	"utils/hashing"
)

func TestToString(t *testing.T) {
	id, _ := hashing.NewKademliaID("456")
	node := Node{IP: "123", ID: id}
	expected := fmt.Sprintf(`node("%s","123")`, node.ID)
	actual := node.String()
	if expected != actual {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestFromString(t *testing.T) {
	node, err := FromString(
		"node(\"1231330000000000000000000000000000000001\",\"abc\")",
	)
	if err != nil {
		t.Errorf("Should not throw error: %v", err)
	}
	expected := "1231330000000000000000000000000000000001"
	actual := node.ID.String()
	if actual != expected {
		t.Errorf(
			"Incorrect id, expected %v got %v",
			expected,
			actual,
		)
	}
}

func TestFromStringInvalidSyntax(t *testing.T) {
	node, err := FromString(
		"(\"1231330000000000000000000000000000000001\",\"abc\")",
	)
	if err == nil {
		t.Errorf("Should throw error, instead got node: %#v", node)
	}
}

func TestFromStringInvalidSyntax2(t *testing.T) {
	node, err := FromString(
		"node(\"1231330000000000000000000000000000000001\")",
	)
	if err == nil {
		t.Errorf("Should throw error, instead got node: %#v", node)
	}
}

func TestFromStringInvalidIDLength(t *testing.T) {
	node, err := FromString(
		"node(\"123\",\"abc\")",
	)
	if err == nil {
		t.Errorf("Should throw error, instead got node: %#v", node)
	}
}

func TestFromStringInvalidIDCharacters(t *testing.T) {
	node, err := FromString(
		"node(\"123133g000000000000000000000000000000001\",\"abc\")",
	)
	if err == nil {
		t.Errorf("Should throw error, instead got node: %#v", node)
	}
}

func TestFromStringsInvalidNodes(t *testing.T) {
	nodes, err := FromStrings(
		"node(\"123\",\"abc\") (123)",
	)
	if err == nil {
		t.Errorf("Should throw error, instead got node: %#v", nodes)
	}
}
func TestFromStrings(t *testing.T) {
	nodes, err := FromStrings(
		"node(\"0000000000000000000000000000000000000000\",\"abc\") " +
			"node(\"1111111111111111111111111111111111111111\",\"def\")",
	)
	if err != nil {
		t.Errorf("Should not throw error, but got: %v", err)
	}
	id1, _ := hashing.ToKademliaID("0000000000000000000000000000000000000000")
	id2, _ := hashing.ToKademliaID("1111111111111111111111111111111111111111")
	expected := [constants.CLOSESTNODES]Node{
		Node{ID: id1, IP: "abc"},
		Node{ID: id2, IP: "def"},
	}
	if len(nodes) != len(expected) ||
		nodes[0].ID.String() != expected[0].ID.String() ||
		nodes[0].IP != expected[0].IP ||
		nodes[1].IP != expected[1].IP ||
		nodes[1].ID.String() != expected[1].ID.String() {
		t.Errorf("Expected len %v got len %v",
			len(expected),
			len(nodes),
		)
		t.Errorf("Position 0 ID: Expected %v got %v",
			expected[0].ID,
			nodes[0].ID,
		)
		t.Errorf("Position 0 IP: Expected %v got %v",
			expected[0].IP,
			nodes[0].IP,
		)
		t.Errorf("Position 1 ID: Expected %v got %v",
			expected[1].ID,
			nodes[1].ID,
		)
		t.Errorf("Position 1 IP: Expected %v got %v",
			expected[1].IP,
			nodes[1].IP,
		)
	}
}
