package structure

import (
	"github.com/angelsolaorbaiceta/inkgeom"
)

// Point where one or more resistant elements meet.
type Node struct {
	Id int
	Position inkgeom.Projectable
	ExternalConstraint Constraint
}

/* Construction */
func MakeNode(id int, position inkgeom.Projectable, externalConstraint Constraint) Node {
	return Node{id, position, externalConstraint}
}

func MakeNodeFromProjs(id int, x, y float64, externalConstraint Constraint) Node {
	return Node{id, inkgeom.MakePoint(x, y), externalConstraint}
}

func MakeFreeNodeFromProjs(id int, x, y float64) Node {
	return Node{id, inkgeom.MakePoint(x, y), MakeNilConstraint()}
}
