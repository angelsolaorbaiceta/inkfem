/*
Structure package defines the structure model used for the
Finite Element Method analysis.
*/
package structure

import (
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
)

// Resistant element defined between two structural nodes, a section and a material.
// An element can have loads applied to it.
type Element struct {
	Id, StartNodeId, EndNodeId int
	Geometry inkgeom.Segment
	StartLink, EndLink Constraint
	material Material
	section Section
	Loads []load.Load
}

/* Construction */
func MakeElement(
	id int,
	startNode, endNode Node,
	startLink, endLink Constraint,
	material Material,
	section Section,
	loads []load.Load) Element {
	return Element{
		id, startNode.Id, endNode.Id,
		inkgeom.MakeSegment(startNode.Position, endNode.Position),
		startLink, endLink,
		material, section, loads}
}

/* Properties */
func (e Element) StartPoint() inkgeom.Projectable {
	return e.Geometry.Start
}

func (e Element) EndPoint() inkgeom.Projectable {
	return e.Geometry.End
}

/* Methods */
func (e Element) IsAxialMember() bool {
	for _, ld := range(e.Loads) {
		if !ld.IsNodal() || ld.Term == load.MZ {
			return false
		}
	}

	return e.StartLink.AllowsRotation() && e.EndLink.AllowsRotation()
}

func (e Element) HasLoadsApplied() bool {
	return len(e.Loads) > 0
}

// func (e Element) StiffnessValue(actionDof, effectDof int, startT, entT inkgeom.TParam) float64 {
// 	return 0.0
// }
