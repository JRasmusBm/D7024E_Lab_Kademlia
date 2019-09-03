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
