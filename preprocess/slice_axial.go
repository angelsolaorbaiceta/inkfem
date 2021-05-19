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

package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

/*
An axial element is an element which:

- is pinned in both ends and
- if it has loads, they are concentrated and applied to the ends of the element, and
never include a moment about Z

Axial elements can be sliced using only it's end nodes. Axial elements deformation only
happens in their axial direction: tension or compression.
*/
func sliceAxialElement(element *structure.Element) *Element {
	if !element.IsAxialMember() {
		panic("Expected an axial element")
	}

	if element.HasLoadsApplied() {
		sFx, sFy, eFx, eFy := netNodalLoadValues(element.ConcentratedLoads, element.Geometry.RefFrame())

		return MakeElement(
			element,
			[]*Node{
				MakeNode(inkgeom.MinT, element.StartPoint(), sFx, sFy, 0.0),
				MakeNode(inkgeom.MaxT, element.EndPoint(), eFx, eFy, 0.0),
			})
	}

	return MakeElement(
		element,
		[]*Node{
			MakeUnloadedNode(inkgeom.MinT, element.StartPoint()),
			MakeUnloadedNode(inkgeom.MaxT, element.EndPoint()),
		})
}

/*
Assuming all loads are nodal (concentrated and applied to the ends of the
element), computes the net, locally projected loads at the start end (sFx & sFy)
and at the end end (eFx & eFy).
*/
func netNodalLoadValues(
	loads []*load.ConcentratedLoad,
	localRefFrame g2d.RefFrame,
) (sFx, sFy, eFx, eFy float64) {
	var localForcesVector g2d.Projectable

	for _, ld := range loads {
		if ld.IsInLocalCoords {
			localForcesVector = ld.ForcesVector()
		} else {
			localForcesVector = localRefFrame.ProjectVector(ld.ForcesVector())
		}

		if ld.T.IsMin() {
			sFx += localForcesVector.X
			sFy += localForcesVector.Y
		} else if ld.T.IsMax() {
			eFx += localForcesVector.X
			eFy += localForcesVector.Y
		}
	}

	return
}
