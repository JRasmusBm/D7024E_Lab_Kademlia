# Kademlia Implementation

This is the repository for an implementation of Kademlia, written in golang.

- The lab instructions can be found [here](./docs/instructions/lab.pdf)
- The paper on Kademlia can be found [here](./docs/instructions/paper.pdf)

# Authors

- August Eriksson
- Rasmus Bergström

# Algorithm Overview

The following information is our interpretation of the Kademlia research paper
linked above.

## Distance Metric

The novel feature of the Kademlia is its use of an XOR distance metric. The
main advantage of this approach is that it facilitates the use of the same
algorithm to perform the entire lookup of a certain key, as opposed to many
other implementations where lookup is split into lookup of the keyspace and
then (basically) linear search inside said space.

The distance is defined as the XOR value of the two nodes.

### Formal Properties

```text
d(x, x) = 0
d(x, y) > 0 if x ≠ y
d(x, y) = d(y, x)
d(x, y) + d(y, z) ≥ d(x, z)
∀ (x, t) ∃! y | d(x, y) = t
```

## BST Subtree Division

Kademlia treats nodes as leaves in a binary tree, where each nodes position is
determined by the shortest unique prefix of its ID. For any given node, the
binary tree is divided into a series of successively lower subtrees which
don't contain the node. The algorithm makes sure the node knows at least one
other node in each of the subtrees (if the subtree has a node). This ensures
that all nodes can reach each other by "asking" the known node in the relevant
subtree.

## _k_-bucket

Each node keeps a list of _k_ other nodes, this list is known as the
_k_-bucket.

- Whenever the node interacts with a known node, it moves it to the head of
  the list, resulting in a list sorted by last encounter.
- Whenever the node encounters an untracked node, it attempts to add it to its
  bucket.
  - If the bucket is full, the tail of the list (the last seen node) is
    pinged. If it is off-line, it is evicted and the new node is added,
    otherwise no change is made to the bucket.
  - If the bucket is not full, it is added to head of the list and will remain
    in the bucket indefinitely unless evicted.

> _k_ is chosen such that any collection of _k_ nodes are unlikely to fail
> within an hour
