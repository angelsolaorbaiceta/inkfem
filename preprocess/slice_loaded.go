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
	"sort"
)

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

func applyDistributedLoadsToNodes(nodes []*Node, element *structure.Element) {
	var (
		trailNode, leadNode *Node
		length, halfLength  float64
		avgLoadValVect      [3]float64
		elemRefFrame        = element.Geometry.RefFrame()
	)

	for i, j := 0, 1; j < len(nodes); i, j = i+1, j+1 {
		trailNode, leadNode = nodes[i], nodes[j]
		length = element.Geometry.LengthBetween(trailNode.T, leadNode.T)
		halfLength = 0.5 * length

		for _, load := range element.Loads {
			avgLoadValVect = load.AvgValueVectorBetween(trailNode.T, leadNode.T)

			if !load.IsInLocalCoords {
				localForces := elemRefFrame.ProjectVector(
					g2d.MakeVector(avgLoadValVect[0], avgLoadValVect[1]),
				)
				avgLoadValVect[0] = localForces.X
				avgLoadValVect[1] = localForces.Y
			}

			trailNode.AddLoad([3]float64{
				avgLoadValVect[0] * halfLength,
				avgLoadValVect[1] * halfLength,
				(avgLoadValVect[2] * halfLength) + (avgLoadValVect[1] * length * length / 12.0),
			})
			leadNode.AddLoad([3]float64{
				avgLoadValVect[0] * halfLength,
				avgLoadValVect[1] * halfLength,
				(avgLoadValVect[2] * halfLength) - (avgLoadValVect[1] * length * length / 12.0),
			})
		}
	}
}
