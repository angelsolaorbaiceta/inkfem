/*
Package structure defines the structure model used for the
Finite Element Method analysis.
*/
package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath/mat"
)

// Element represents s resistant element defined between two structural nodes, a section and a material.
// An element can have loads applied to it.
type Element struct {
	Id, StartNodeId, EndNodeId int
	Geometry                   inkgeom.Segment
	StartLink, EndLink         *Constraint
	material                   Material
	section                    Section
	Loads                      []load.Load
}

/* ::::::::::::::: Construction ::::::::::::::: */

// MakeElement creates a new element with all properties initialized.
func MakeElement(
	id int,
	startNode, endNode Node,
	startLink, endLink *Constraint,
	material Material,
	section Section,
	loads []load.Load) *Element {
	return &Element{
		id, startNode.Id, endNode.Id,
		inkgeom.MakeSegment(startNode.Position, endNode.Position),
		startLink, endLink,
		material, section, loads}
}

/* ::::::::::::::: Properties ::::::::::::::: */

// StartPoint returns the position of the start node of this element's geometry.
func (e Element) StartPoint() inkgeom.Projectable {
	return e.Geometry.Start
}

// EndPoint returns the position of the end node of this element's geometry.
func (e Element) EndPoint() inkgeom.Projectable {
	return e.Geometry.End
}

// PointAt returns the position of an intermediate point in this element's geometry.
func (e Element) PointAt(t inkgeom.TParam) inkgeom.Projectable {
	return e.Geometry.PointAt(t)
}

/* ::::::::::::::: Methods ::::::::::::::: */

/*
IsAxialMember returns true if this element is pinned in both ends and, in case of having loads
applied, they are always in the end positions of the directrix and do not include moments about Z,
but just forces in X and Y directions.
*/
func (e Element) IsAxialMember() bool {
	for _, ld := range e.Loads {
		if !ld.IsNodal() || ld.Term == load.MZ {
			return false
		}
	}

	return e.StartLink.AllowsRotation() && e.EndLink.AllowsRotation()
}

// HasLoadsApplied returns true if any load of any type is applied to the element.
func (e Element) HasLoadsApplied() bool {
	return len(e.Loads) > 0
}

/*
StiffnessGlobalMat generates the local stiffness matrix for the element and applies
the rotation defined by the elements' geometry reference frame.
*/
func (e Element) StiffnessGlobalMat(startT, entT inkgeom.TParam) mat.Matrixable {
	k := mat.MakeSquareDense(6)
	// TODO: implement
	return k
}

/* ::::::::::::::: Stringer ::::::::::::::: */
func (e Element) String() string {
	return fmt.Sprintf(
		"%d -> %d%s %d%s %s %s",
		e.Id,
		e.StartNodeId, e.StartLink.String(),
		e.EndNodeId, e.EndLink.String(),
		e.material.Name, e.section.Name)
}
