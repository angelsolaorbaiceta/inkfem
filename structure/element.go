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

	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkmath/mat"
)

/*
An Element represents s resistant element defined between two structural nodes,
a section and a material.

An Element can have loads applied to it.

TODO: choose the bending axis
*/
type Element struct {
	Id, StartNodeId, EndNodeId int
	Geometry                   g2d.Segment
	StartLink, EndLink         Constraint
	material                   *Material
	section                    *Section
	Loads                      []load.Load
	_ea, _ei                   float64
}

/* <-- Construction --> */

/*
MakeElement creates a new element with all properties initialized.
*/
func MakeElement(
	id int,
	startNode, endNode *Node,
	startLink, endLink Constraint,
	material *Material,
	section *Section,
	loads []load.Load,
) *Element {
	return &Element{
		Id:          id,
		StartNodeId: startNode.Id,
		EndNodeId:   endNode.Id,
		Geometry:    g2d.MakeSegment(startNode.Position, endNode.Position),
		StartLink:   startLink,
		EndLink:     endLink,
		material:    material,
		section:     section,
		Loads:       loads,
		_ea:         material.YoungMod * section.Area,
		_ei:         material.YoungMod * section.IStrong,
	}
}

/* <-- Properties --> */

/*
StartPoint returns the position of the start node of this element's geometry.
*/
func (e Element) StartPoint() g2d.Projectable {
	return e.Geometry.Start
}

/*
EndPoint returns the position of the end node of this element's geometry.
*/
func (e Element) EndPoint() g2d.Projectable {
	return e.Geometry.End
}

/*
PointAt returns the position of a middle point in this element's geometry.
*/
func (e Element) PointAt(t inkgeom.TParam) g2d.Projectable {
	return e.Geometry.PointAt(t)
}

/*
Material returns the material for the element.
*/
func (e Element) Material() *Material {
	return e.material
}

/*
Section returns the section for the element.
*/
func (e Element) Section() *Section {
	return e.section
}

/*
HasLoadsApplied returns true if any load of any type is applied to the element.
*/
func (e Element) HasLoadsApplied() bool {
	return len(e.Loads) > 0
}

/* <-- Methods --> */

/*
IsAxialMember returns true if this element is pinned in both ends and, in case
of having loads applied, they are always in the end positions of the directrix
and does not include moments about Z, but just forces in X and Y directions.
*/
func (e Element) IsAxialMember() bool {
	for _, ld := range e.Loads {
		if !ld.IsNodal() || ld.Term == load.MZ {
			return false
		}
	}

	return e.StartLink.AllowsRotation() && e.EndLink.AllowsRotation()
}

/*
StiffnessGlobalMat generates the local stiffness matrix for the element and
applies the rotation defined by the elements' geometry reference frame.

It returns the element's stiffness matrix in the global reference frame.
*/
func (e Element) StiffnessGlobalMat(startT, endT inkgeom.TParam) mat.ReadOnlyMatrix {
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

/*
Equals tests whether this element is equal to other.
*/
func (e Element) Equals(other Element) bool {
	return e.StartNodeId == other.StartNodeId &&
		e.EndNodeId == other.EndNodeId &&
		e.StartLink.Equals(other.StartLink) &&
		e.EndLink.Equals(other.EndLink) &&
		e.material.Name == other.material.Name &&
		e.section.Name == other.section.Name
}

/* <-- Stringer --> */

func (e Element) String() string {
	return fmt.Sprintf(
		"%d -> %d%s %d%s %s %s",
		e.Id,
		e.StartNodeId, e.StartLink.String(),
		e.EndNodeId, e.EndLink.String(),
		e.material.Name, e.section.Name,
	)
}
