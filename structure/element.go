package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/angelsolaorbaiceta/inkmath/mat"
)

// An Element represents a resistant element defined between two structural nodes, a section and
// a material.
//
// An Element can have distributed and concentrated loads applied to it.
//
// To create an element, use the `ElementBuilder`.
//
// TODO: choose the bending axis
// TODO: buckling analysis
type Element struct {
	id, startNodeID, endNodeID contracts.StrID
	geometry                   *g2d.Segment
	startLink, endLink         *Constraint
	material                   *Material
	section                    *Section
	ConcentratedLoads          []*load.ConcentratedLoad
	DistributedLoads           []*load.DistributedLoad
}

func (e Element) GetID() contracts.StrID {
	return e.id
}

func (e Element) StartNodeID() contracts.StrID {
	return e.startNodeID
}

func (e Element) EndNodeID() contracts.StrID {
	return e.endNodeID
}

func (e Element) RefFrame() *g2d.RefFrame {
	return e.geometry.RefFrame()
}

func (e Element) DirectionVersor() *g2d.Vector {
	return e.geometry.DirectionVersor()
}

func (e Element) NormalVersor() *g2d.Vector {
	return e.geometry.NormalVersor()
}

func (e Element) Length() float64 {
	return e.geometry.Length()
}

func (e Element) LengthBetween(tStart, tEnd nums.TParam) float64 {
	return e.geometry.LengthBetween(tStart, tEnd)
}

// StartPoint returns the position of the start node of this element's geometry.
func (e Element) StartPoint() *g2d.Point {
	return e.geometry.Start()
}

// EndPoint returns the position of the end node of this element's geometry.
func (e Element) EndPoint() *g2d.Point {
	return e.geometry.End()
}

// PointAt returns the position of a middle point in this element's geometry.
func (e Element) PointAt(t nums.TParam) *g2d.Point {
	return e.geometry.PointAt(t)
}

func (e Element) StartLink() *Constraint {
	return e.startLink
}

func (e Element) EndLink() *Constraint {
	return e.endLink
}

// Material returns the elements's material.
func (e Element) Material() *Material {
	return e.material
}

// Section returns the element's section.
func (e Element) Section() *Section {
	return e.section
}

// LoadsCount is the total number of concentrated and distributed loads applied to the element.
func (e Element) LoadsCount() int {
	return len(e.ConcentratedLoads) + len(e.DistributedLoads)
}

// HasLoadsApplied returns true if any load, either concentrated of distributed,
// is applied to the element.
func (e Element) HasLoadsApplied() bool {
	return e.LoadsCount() > 0
}

// AddOwnWeight adds a distributed load to the element that represents the weight
// of the element. The load's intensity per unit length is the product of the
// material's density and the section's area. The load is applied in the negative
// direction of the global Y axis.
func (e *Element) AddOwnWeight() {
	var (
		value = -e.material.Density * e.section.Area
		load  = load.MakeDistributed(load.FY, false, nums.MinT, value, nums.MaxT, value)
	)

	e.DistributedLoads = append(e.DistributedLoads, load)
}

// IsAxialMember returns true if this element is pinned in both ends and, in case of having
// loads applied, they are always in the end positions of the directrix and does not include
// moments about Z, but just forces in X and Y directions.
//
// FIXME: is axial member a good name? a distributed Fx load would make this bar not an
// axial member, which seems weird...
func (e Element) IsAxialMember() bool {
	if len(e.DistributedLoads) > 0 {
		return false
	}

	for _, ld := range e.ConcentratedLoads {
		if !ld.IsNodal() || ld.Term == load.MZ {
			return false
		}
	}

	return e.startLink.AllowsRotation() && e.endLink.AllowsRotation()
}

// StiffnessGlobalMat generates the local stiffness matrix for the element and applies the rotation
// defined by the elements' geometry reference frame.
//
// It returns the element's stiffness matrix in the global reference frame.
func (e Element) StiffnessGlobalMat(startT, endT nums.TParam) mat.ReadOnlyMatrix {
	var (
		l    = e.geometry.LengthBetween(startT, endT)
		c    = e.geometry.RefFrame().Cos()
		s    = e.geometry.RefFrame().Sin()
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
// Loads aren't compared, two bars with different set of loads might therefore be equal.
func (e *Element) Equals(other *Element) bool {
	return e.startNodeID == other.startNodeID &&
		e.endNodeID == other.endNodeID &&
		e.startLink.Equals(other.startLink) &&
		e.endLink.Equals(other.endLink) &&
		e.material.Name == other.material.Name &&
		e.section.Name == other.section.Name
}

// String representation of the bar.
func (e Element) String() string {
	return fmt.Sprintf(
		"%s: %s %s %s %s '%s' '%s'",
		e.id,
		e.startNodeID,
		e.startLink.String(),
		e.endNodeID,
		e.endLink.String(),
		e.material.Name,
		e.section.Name,
	)
}
