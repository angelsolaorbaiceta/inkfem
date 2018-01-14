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

// SetDegreesOfFreedomNums assigns numbers to the degress of freedom of the node.
func (n *Node) SetDegreesOfFreedomNums(dx, dy, rz int) {
	n.globalDof[0] = dx
	n.globalDof[1] = dy
	n.globalDof[2] = rz
}

// DegreesOfFreedomNums returns the degrees of freedom numbers assigned to the node.
func (n Node) DegreesOfFreedomNums() [3]int {
	return n.globalDof
}

/* ::::::::::::::: Stringer ::::::::::::::: */
func (n Node) String() string {
	return fmt.Sprintf(
		"%d -> %f %f %s | DOF: %v",
		n.Id, n.Position.X, n.Position.Y,
		n.ExternalConstraint.String(),
		n.DegreesOfFreedomNums(),
	)
}
