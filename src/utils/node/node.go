package node

import "utils/hashing"

type Node struct {
	ID *hashing.KademliaID
    IP string
    Dist int
}
