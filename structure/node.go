package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

const unsetDOFNumber = -1

// Node is a point in the structure where one or more resistant elements meet.
type Node struct {
	id                 contracts.StrID
	Position           *g2d.Point
	ExternalConstraint *Constraint
	globalDof          [3]int
}

// MakeNode creates a new node with the given id, position and external constraint.
func MakeNode(
	id contracts.StrID,
	position *g2d.Point,
	externalConstraint *Constraint,
) *Node {
	return &Node{
		id,
		position,
		externalConstraint,
		[3]int{unsetDOFNumber, unsetDOFNumber, unsetDOFNumber},
	}
}

// MakeNodeAtPosition creates a new node with the given id, position coordinates and external constraint.
func MakeNodeAtPosition(
	id contracts.StrID,
	x, y float64,
	externalConstraint *Constraint,
) *Node {
	return &Node{
		id,
		g2d.MakePoint(x, y),
		externalConstraint,
		[3]int{unsetDOFNumber, unsetDOFNumber, unsetDOFNumber},
	}
}

/*
MakeFreeNodeAtPosition creates a new node without external constraint, with the given id and
position by coordinates.
*/
func MakeFreeNodeAtPosition(id contracts.StrID, x, y float64) *Node {
	return &Node{
		id,
		g2d.MakePoint(x, y),
		&NilConstraint,
		[3]int{unsetDOFNumber, unsetDOFNumber, unsetDOFNumber},
	}
}

// IsExternallyConstrained returns true if this node is externally constrained.
func (n Node) IsExternallyConstrained() bool {
	return n.ExternalConstraint != &NilConstraint
}

// DegreesOfFreedomNum returns the degrees of freedom numbers assigned to the node.
func (n Node) DegreesOfFreedomNum() [3]int {
	return n.globalDof
}

// HasDegreesOfFreedomNum returns true if the node has already been assigned degress of freedom.
func (n Node) HasDegreesOfFreedomNum() bool {
	return n.globalDof[0] != unsetDOFNumber &&
		n.globalDof[1] != unsetDOFNumber &&
		n.globalDof[2] != unsetDOFNumber
}

// SetDegreesOfFreedomNum assigns numbers to the degress of freedom of the node.
func (n *Node) SetDegreesOfFreedomNum(dx, dy, rz int) {
	n.globalDof[0] = dx
	n.globalDof[1] = dy
	n.globalDof[2] = rz
}

// Equals tests whether this node and other are equal.
func (n *Node) Equals(other *Node) bool {
	return n.Position.Equals(other.Position) &&
		n.ExternalConstraint.Equals(other.ExternalConstraint) &&
		n.globalDof[0] == other.globalDof[0] &&
		n.globalDof[1] == other.globalDof[1] &&
		n.globalDof[2] == other.globalDof[2]
}

// GetID returns the node's id.
func (n Node) GetID() contracts.StrID {
	return n.id
}

func (n Node) String() string {
	return fmt.Sprintf(
		"%s -> %f %f %s | %v",
		n.id, n.Position.X(), n.Position.Y(),
		n.ExternalConstraint.String(),
		n.DegreesOfFreedomNum(),
	)
}
