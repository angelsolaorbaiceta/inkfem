package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkmath/mat"
)

/*
An Element represents a resistant element defined between two structural nodes, a section and
a material.

An Element can have distributed and concentrated loads applied to it.

TODO: choose the bending axis
TODO: buckling analysis
*/
type Element struct {
	Id, StartNodeID, EndNodeID contracts.StrID
	Geometry                   g2d.Segment
	StartLink, EndLink         *Constraint
	material                   *Material
	section                    *Section
	ConcentratedLoads          []*load.ConcentratedLoad
	DistributedLoads           []*load.DistributedLoad
}

// MakeElement creates a new element with all properties initialized.
func MakeElement(
	id contracts.StrID,
	startNode, endNode *Node,
	startLink, endLink *Constraint,
	material *Material,
	section *Section,
	concentratedLoads []*load.ConcentratedLoad,
	distributedLoads []*load.DistributedLoad,
) *Element {
	return &Element{
		Id:                id,
		StartNodeID:       startNode.Id,
		EndNodeID:         endNode.Id,
		Geometry:          g2d.MakeSegment(startNode.Position, endNode.Position),
		StartLink:         startLink,
		EndLink:           endLink,
		material:          material,
		section:           section,
		ConcentratedLoads: concentratedLoads,
		DistributedLoads:  distributedLoads,
	}
}

// MakeElementWithoutLoads creates a new element with no external loads.
func MakeElementWithoutLoads(
	id contracts.StrID,
	startNode, endNode *Node,
	startLink, endLink *Constraint,
	material *Material,
	section *Section,
) *Element {
	return MakeElement(
		id,
		startNode, endNode,
		startLink, endLink,
		material, section,
		[]*load.ConcentratedLoad{},
		[]*load.DistributedLoad{},
	)
}

// StartPoint returns the position of the start node of this element's geometry.
func (e Element) StartPoint() g2d.Projectable {
	return e.Geometry.Start
}

// EndPoint returns the position of the end node of this element's geometry.
func (e Element) EndPoint() g2d.Projectable {
	return e.Geometry.End
}

// PointAt returns the position of a middle point in this element's geometry.
func (e Element) PointAt(t inkgeom.TParam) g2d.Projectable {
	return e.Geometry.PointAt(t)
}

// Material returns the material for the element.
func (e Element) Material() *Material {
	return e.material
}

// Section returns the section for the element.
func (e Element) Section() *Section {
	return e.section
}

// HasLoadsApplied returns true if any load of any type is applied to the element.
func (e Element) HasLoadsApplied() bool {
	return len(e.ConcentratedLoads) > 0 || len(e.DistributedLoads) > 0
}

/*
IsAxialMember returns true if this element is pinned in both ends and, in case of having loads
applied, they are always in the end positions of the directrix and does not include moments about Z,
but just forces in X and Y directions.

FIXME: is axial member a good name? a distributed Fx load would mean this is not an axiam member,
which seems weird...
*/
func (e Element) IsAxialMember() bool {
	if len(e.DistributedLoads) > 0 {
		return false
	}

	for _, ld := range e.ConcentratedLoads {
		if !ld.IsNodal() || ld.Term == load.MZ {
			return false
		}
	}

	return e.StartLink.AllowsRotation() && e.EndLink.AllowsRotation()
}

/*
StiffnessGlobalMat generates the local stiffness matrix for the element and applies the rotation
defined by the elements' geometry reference frame.

It returns the element's stiffness matrix in the global reference frame.
*/
func (e Element) StiffnessGlobalMat(startT, endT inkgeom.TParam) mat.ReadOnlyMatrix {
	var (
		l    = e.Geometry.LengthBetween(startT, endT)
		c    = e.Geometry.RefFrame().Cos()
		s    = e.Geometry.RefFrame().Sin()
		ea   = e.material.YoungMod * e.section.Area
		ei   = e.material.YoungMod * e.section.IStrong
		c2   = c * c
		s2   = s * s
		cs   = c * s
		eal  = ea / l
		eil3 = 12.0 * ei / (l * l * l)
		eil2 = 6.0 * ei / (l * l)
		eil  = ei / l
		k    = mat.MakeSquareDense(6)
	)

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

// Equals tests whether this element is equal to other.
func (e *Element) Equals(other *Element) bool {
	return e.StartNodeID == other.StartNodeID &&
		e.EndNodeID == other.EndNodeID &&
		e.StartLink.Equals(other.StartLink) &&
		e.EndLink.Equals(other.EndLink) &&
		e.material.Name == other.material.Name &&
		e.section.Name == other.section.Name
}

// GetID returns the element's id. Implements Identifiable interface.
func (e Element) GetID() contracts.StrID {
	return e.Id
}

func (e Element) String() string {
	return fmt.Sprintf(
		"%s -> %s %s %s %s %s %s",
		e.Id,
		e.StartNodeID, e.StartLink.String(),
		e.EndNodeID, e.EndLink.String(),
		e.material.Name, e.section.Name,
	)
}
