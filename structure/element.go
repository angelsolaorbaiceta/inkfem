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
	_ea, _ei                   float64
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
		material, section, loads,
		material.YieldStrength * section.Area,
		material.YieldStrength * section.IStrong,
	}
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
func (e Element) StiffnessGlobalMat(startT, endT inkgeom.TParam) mat.Matrixable {
	var (
		l    = e.Geometry.LengthBetween(startT, endT)
		c    = e.Geometry.RefFrame().Cos()
		s    = e.Geometry.RefFrame().Sin()
		c2   = c * c
		s2   = s * s
		cs   = c * s
		eal  = e._ea / l
		eil3 = 12.0 * e._ei / (l * l * l)
		eil2 = 6.0 * e._ei / (l * l)
		eil  = e._ei / l
	)

	k := mat.MakeSquareDense(6)

	// First Row
	k.SetValue(0, 0, (c2*eal + s2*eil3))
	k.SetValue(0, 1, (cs*eal - cs*eil3))
	k.SetValue(0, 2, -s*eil2)
	k.SetValue(0, 3, -c2*eal-s2*eil3)
	k.SetValue(0, 4, (-cs*eal + cs*eil3))
	k.SetValue(0, 5, -s*eil2)

	// Second Row
	k.SetValue(1, 0, (cs*eal - cs*eil3))
	k.SetValue(1, 1, (s2*eal + c2*eil3))
	k.SetValue(1, 2, c*eil2)
	k.SetValue(1, 3, (-cs*eal + cs*eil3))
	k.SetValue(1, 4, (-s2*eal - c2*eil3))
	k.SetValue(1, 5, c*eil2)

	// Third Row
	k.SetValue(2, 0, -s*eil2)
	k.SetValue(2, 1, c*eil2)
	k.SetValue(2, 2, 4.0*eil)
	k.SetValue(2, 3, s*eil2)
	k.SetValue(2, 4, -c*eil2)
	k.SetValue(2, 5, 2.0*eil)

	// Fourth Row
	k.SetValue(3, 0, (-c2*eal - s2*eil3))
	k.SetValue(3, 1, (-cs*eal + cs*eil3))
	k.SetValue(3, 2, s*eil2)
	k.SetValue(3, 3, (c2*eal + s2*eil3))
	k.SetValue(3, 4, (cs*eal - cs*eil3))
	k.SetValue(3, 5, s*eil2)

	// Fifth Row
	k.SetValue(4, 0, (-cs*eal + cs*eil3))
	k.SetValue(4, 1, (-s2*eal - c2*eil3))
	k.SetValue(4, 2, -c*eil2)
	k.SetValue(4, 3, (cs*eal - cs*eil3))
	k.SetValue(4, 4, (s2*eal + c2*eil3))
	k.SetValue(4, 5, -c*eil2)

	// Sixth Row
	k.SetValue(5, 0, -s*eil2)
	k.SetValue(5, 1, c*eil2)
	k.SetValue(5, 2, 2.0*eil)
	k.SetValue(5, 3, s*eil2)
	k.SetValue(5, 4, -c*eil2)
	k.SetValue(5, 5, 4.0*eil)

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
