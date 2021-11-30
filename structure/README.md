# Structure package

The structure package defines the structural model.

A structure is defined as a [`Structure`](./structure.go) struct, which contains a map of nodes by id and a slice of elements.

## Nodes

Nodes are represented by the [`Node`](./node.go) struct.
A node is identified by a unique id, has a position in the plane and an optional external constraint.

A `Node` can be created using one of the following functions:

- `MakeNode`: requires an id, a position and external constraint
- `MakeNodeAtPosition`: requires an id, x and y coordinates and an external constraint
- `MakeFreeNodeAtPosition`: creates a non-constrained node with the given id and position.

## Elements
