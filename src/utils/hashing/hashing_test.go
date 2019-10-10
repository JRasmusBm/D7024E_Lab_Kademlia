package hashing

import (
	"testing"
)

func TestNewKademliaID(t *testing.T) {
	expected := "c412b37f8c0484e6db8bce177ae88c5443b26e92"
	actual := NewKademliaID("hej")
	if expected != actual.String() {
		t.Errorf("Expected %v got %v", expected, actual)
	}
}

func TestToKademliaID(t *testing.T) {
	expected := "c412b37f8c0484e6db8bce177ae88c5443b26e92"
	actual, err := ToKademliaID(expected)
	if err != nil {
		t.Error(err)
	}
	if expected != actual.String() {
		t.Errorf("Expected %v got %v", expected, actual.String())
	}
}

func TestToKademliaIDInvalidLength(t *testing.T) {
	expected := "c412b37f8c0484e6db8bce177ae88c5443b26e"
	_, err := ToKademliaID(expected)
	if err == nil {
		t.Errorf("Should throw an error")
	}
}

func TestToKademliaIDInvalidCharacters(t *testing.T) {
	expected := "c412b37f8c0484e6db8bcex77ae88c5443b26e9"
	_, err := ToKademliaID(expected)
	if err == nil {
		t.Errorf("Should throw an error")
	}
}

func TestRandomKademliaID(t *testing.T) {
	key := NewRandomKademliaID()
	_, err := ToKademliaID(key.String())
	if err != nil {
		t.Errorf("Should create a valid KademliaID")
	}
}

func TestLessThan(t *testing.T) {
	id1, _ := ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := ToKademliaID("0000000000000000000000000000000000000000")
	if id1.Less(id2) {
		t.Errorf("Expected %s to be less than %s", id1.String(), id2.String())
	}
}

func TestLessThanSameSize(t *testing.T) {
	id1, _ := ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := ToKademliaID("1111111111111111111111111111111111111111")
	if id2.Less(id1) {
		t.Errorf("Expected %s not to be less than %s", id2.String(), id1.String())
	}
}

func TestNotEqual(t *testing.T) {
	id1, _ := ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := ToKademliaID("0000000000000000000000000000000000000000")
	if id1.Equals(id2) {
		t.Errorf("Expected %s not to be equal to %s", id1.String(), id2.String())
	}
}

func TestEqual(t *testing.T) {
	id1, _ := ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := ToKademliaID("1111111111111111111111111111111111111111")
	if !id1.Equals(id2) {
		t.Errorf("Expected %s not to be equal to %s", id1.String(), id2.String())
	}
}

func TestDistanceToSelf(t *testing.T) {
	id1, _ := ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := ToKademliaID("1111111111111111111111111111111111111111")
	expected, _ := ToKademliaID("0000000000000000000000000000000000000000")
	actual := id1.CalcDistance(id2)
	if actual.String() != expected.String() {
		t.Errorf("Expected %v got %v", expected.String(), actual.String())
	}
}

func TestDistanceToAnotherNode(t *testing.T) {
	id1, _ := ToKademliaID("1111111111111111111111111111111111111111")
	id2, _ := ToKademliaID("0111111111111111111111111111111111111111")
	expected, _ := ToKademliaID("1000000000000000000000000000000000000000")
	actual := id1.CalcDistance(id2)
	if actual.String() != expected.String() {
		t.Errorf("Expected %v got %v", expected.String(), actual.String())
	}
}
