package node

import (
    hashing "utils/hashing"
)

type Node struct {
	ID *hashing.KademliaID
    IP string
    Dist int
    RTable RoutingTable
}
