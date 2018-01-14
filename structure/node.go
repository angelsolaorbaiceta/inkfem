package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkgeom"
)

// Node is a point in the structure where one or more resistant elements meet.
type Node struct {
	Id                 int
	Position           inkgeom.Projectable
	ExternalConstraint *Constraint
	globalDof          [3]int
}

/* ::::::::::::::: Construction ::::::::::::::: */

/*
MakeNode creates a new node with the given id, position and external constraint.
*/
func MakeNode(id int, position inkgeom.Projectable, externalConstraint *Constraint) Node {
	return Node{id, position, externalConstraint, [3]int{0, 0, 0}}
}

/*
MakeNodeFromProjs creates a new node with the given id, position coordinates and
external constraint.
*/
func MakeNodeFromProjs(id int, x, y float64, externalConstraint *Constraint) Node {
	return Node{id, inkgeom.MakePoint(x, y), externalConstraint, [3]int{0, 0, 0}}
}

/*
MakeFreeNodeFromProjs creates a new node without external constraint, with the
given id and position by coordinates.
*/
func MakeFreeNodeFromProjs(id int, x, y float64) Node {
	return Node{id, inkgeom.MakePoint(x, y), MakeNilConstraint(), [3]int{0, 0, 0}}
}

/* ::::::::::::::: Properties ::::::::::::::: */

// SetDegreesOfFreedomNum assigns numbers to the degress of freedom of the node.
func (n *Node) SetDegreesOfFreedomNum(dx, dy, rz int) {
	n.globalDof[0] = dx
	n.globalDof[1] = dy
	n.globalDof[2] = rz
}

// DegreesOfFreedomNum returns the degrees of freedom numbers assigned to the node.
func (n Node) DegreesOfFreedomNum() [3]int {
	return n.globalDof
}

// HasDegreesOfFreedomNum returns true if the node has already been assigned degress of freedom.
func (n Node) HasDegreesOfFreedomNum() bool {
	return n.globalDof[0] != 0 || n.globalDof[1] != 0 || n.globalDof[2] != 0
}

/* ::::::::::::::: Stringer ::::::::::::::: */
func (n Node) String() string {
	return fmt.Sprintf(
		"%d -> %f %f %s | DOF: %v",
		n.Id, n.Position.X, n.Position.Y,
		n.ExternalConstraint.String(),
		n.DegreesOfFreedomNum(),
	)
}
