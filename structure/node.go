/*
Copyright 2020 Angel Sola Orbaiceta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

const unsetDOFNumber = -1

/*
Node is a point in the structure where one or more resistant elements meet.
*/
type Node struct {
	Id                 contracts.StrID
	Position           g2d.Projectable
	ExternalConstraint Constraint
	globalDof          [3]int
}

/* <-- Construction --> */

/*
MakeNode creates a new node with the given id, position and external constraint.
*/
func MakeNode(
	id contracts.StrID,
	position g2d.Projectable,
	externalConstraint Constraint,
) *Node {
	return &Node{
		id,
		position,
		externalConstraint,
		[3]int{unsetDOFNumber, unsetDOFNumber, unsetDOFNumber},
	}
}

/*
MakeNodeAtPosition creates a new node with the given id, position coordinates and
external constraint.
*/
func MakeNodeAtPosition(
	id contracts.StrID,
	x, y float64,
	externalConstraint Constraint,
) *Node {
	return &Node{
		id,
		g2d.MakePoint(x, y),
		externalConstraint,
		[3]int{unsetDOFNumber, unsetDOFNumber, unsetDOFNumber},
	}
}

/*
MakeFreeNodeAtPosition creates a new node without external constraint, with the
given id and position by coordinates.
*/
func MakeFreeNodeAtPosition(id contracts.StrID, x, y float64) *Node {
	return &Node{
		id,
		g2d.MakePoint(x, y),
		NilConstraint,
		[3]int{unsetDOFNumber, unsetDOFNumber, unsetDOFNumber},
	}
}

/* <-- Properties --> */

/*
IsExternallyConstrained returns true if this node is externally constrained.
*/
func (n Node) IsExternallyConstrained() bool {
	return n.ExternalConstraint != NilConstraint
}

/*
DegreesOfFreedomNum returns the degrees of freedom numbers assigned to the node.
*/
func (n Node) DegreesOfFreedomNum() [3]int {
	return n.globalDof
}

/*
HasDegreesOfFreedomNum returns true if the node has already been assigned
degress of freedom.
*/
func (n Node) HasDegreesOfFreedomNum() bool {
	return n.globalDof[0] != unsetDOFNumber &&
		n.globalDof[1] != unsetDOFNumber &&
		n.globalDof[2] != unsetDOFNumber
}

/* <-- Methods --> */

/*
SetDegreesOfFreedomNum assigns numbers to the degress of freedom of the node.
*/
func (n *Node) SetDegreesOfFreedomNum(dx, dy, rz int) {
	n.globalDof[0] = dx
	n.globalDof[1] = dy
	n.globalDof[2] = rz
}

/*
Equals tests whether this node and other are equal.
*/
func (n *Node) Equals(other *Node) bool {
	return n.Position.Equals(other.Position) &&
		n.ExternalConstraint.Equals(other.ExternalConstraint)
}

/* <-- Identifiable --> */

/*
GetID returns the node's id.
*/
func (n Node) GetID() contracts.StrID {
	return n.Id
}

/* <-- Stringer --> */

func (n Node) String() string {
	return fmt.Sprintf(
		"%s -> %f %f %s | DOF: %v",
		n.Id, n.Position.X, n.Position.Y,
		n.ExternalConstraint.String(),
		n.DegreesOfFreedomNum(),
	)
}
