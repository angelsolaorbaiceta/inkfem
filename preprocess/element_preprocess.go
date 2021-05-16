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
	"sort"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

const (
	elementWithLoadsSlices    = 10
	elementWithoutLoadsSlices = 7
)

/*
DoElement preprocesses the given structural element subdividing it as corresponds.
The result is sent through a channel.
*/
func DoElement(e *structure.Element, c chan<- *Element) {
	if e.IsAxialMember() {
		c <- sliceAxialElement(e)
	} else if e.HasLoadsApplied() {
		c <- sliceLoadedElement(e, elementWithLoadsSlices)
	} else {
		c <- sliceElementWithoutLoads(e, elementWithoutLoadsSlices)
	}
}

/* <-- Non Sliced --> */

/*
An axial element is an element which:

- is pinned in both ends and
- if it has loads, they are concentrated and applied to the ends of the element, and
never include a moment about Z

Axial elements can be sliced using only it's end nodes. Axial elements deformation only
happens in their axial direction: tension or compression.
*/
func sliceAxialElement(element *structure.Element) *Element {
	if element.HasLoadsApplied() {
		sFx, sFy, eFx, eFy := netNodalLoadValues(element.Loads, element.Geometry.RefFrame())

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
	loads []load.Load,
	localRefFrame g2d.RefFrame,
) (sFx, sFy, eFx, eFy float64) {
	var localForcesVector g2d.Projectable

	for _, ld := range loads {
		if ld.IsInLocalCoords {
			localForcesVector = ld.ForcesVector()
		} else {
			localForcesVector = localRefFrame.ProjectVector(ld.ForcesVector())
		}

		if ld.T().IsMin() {
			sFx += localForcesVector.X
			sFy += localForcesVector.Y
		} else if ld.T().IsMax() {
			eFx += localForcesVector.X
			eFy += localForcesVector.Y
		}
	}

	return
}

/* <-- Sliced : Loaded --> */

/*
Elemets with loads applied are firts sliced a given number of times, and then, all t
parameters derived from the positions of the applied loads are included.

The positions where concentrated loads are applied are critical as there will be a
discontinuity, so a node must be added.

The positions where distributed loads start and end also introduce discontinuities, so we
also include nodes in those positions.
*/
func sliceLoadedElement(element *structure.Element, slices int) *Element {
	tPos := sliceLoadedElementPositions(element.Loads, slices)
	nodes := makeNodesWithConcentratedLoads(element, tPos)
	applyDistributedLoadsToNodes(nodes, element)

	return MakeElement(element, nodes)
}

/*
Computes all the t values where to slice an element with loads applied.

It starts by slicing the element a given number of times, and then adds all the load start
and end t values, removing any possible duplications.
*/
func sliceLoadedElementPositions(loads []load.Load, slices int) []inkgeom.TParam {
	tPos := append(
		inkgeom.SubTParamCompleteRangeTimes(slices),
		tValsForLoadApplications(loads)...,
	)

	sort.Sort(inkgeom.ByTParamValue(tPos))

	var correctedTPos []inkgeom.TParam
	correctedTPos = append(correctedTPos, tPos[0])
	for i := 1; i < len(tPos); i++ {
		if tPos[i-1].DistanceTo(tPos[i]) > 1e-3 {
			correctedTPos = append(correctedTPos, tPos[i])
		}
	}

	return correctedTPos
}

func tValsForLoadApplications(loads []load.Load) []inkgeom.TParam {
	var tVals []inkgeom.TParam

	for _, ld := range loads {
		if ld.IsConcentrated() && !ld.T().IsExtreme() {
			tVals = append(tVals, ld.T())
		} else if ld.IsDistributed() {
			if !ld.StartT().IsExtreme() {
				tVals = append(tVals, ld.StartT())
			}

			if !ld.EndT().IsExtreme() {
				tVals = append(tVals, ld.EndT())
			}
		}
	}

	return tVals
}

/*
Creates all the nodes for the given t positions and applies the concentrated loads
on them.
*/
func makeNodesWithConcentratedLoads(element *structure.Element, tPos []inkgeom.TParam) []*Node {
	nodes := make([]*Node, len(tPos))
	elemRefFrame := element.Geometry.RefFrame()

	for i, t := range tPos {
		node := MakeUnloadedNode(t, element.Geometry.PointAt(t))

		for _, load := range element.Loads {
			if load.IsConcentrated() && t.Equals(load.T()) {
				var localLoadForces g2d.Projectable
				
				if load.IsInLocalCoords {
					localLoadForces = load.ForcesVector()
				} else {
					localLoadForces = elemRefFrame.ProjectVector(load.ForcesVector())
				}

				node.AddLoad(
					[3]float64{localLoadForces.X, localLoadForces.Y, load.VectorValue()[2]},
				)
			}
		}

		nodes[i] = node
	}

	return nodes
}

/* <-- Sliced : Unloaded --> */

/*
Non axial elements which have no loads applied are sliced just by subdividing their
geometry into a given number of slices, so that the slices have the same length.
*/
func sliceElementWithoutLoads(e *structure.Element, slices int) *Element {
	tPos := inkgeom.SubTParamCompleteRangeTimes(slices)
	nodes := make([]*Node, len(tPos))

	for i := 0; i < len(tPos); i++ {
		nodes[i] = MakeUnloadedNode(tPos[i], e.PointAt(tPos[i]))
	}

	return MakeElement(e, nodes)
}
