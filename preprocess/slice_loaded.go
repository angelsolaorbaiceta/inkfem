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
Minimum distance between two consecutive t values in the slices.
*/
const minDistBetweenTSlices = 1e-3

/*
Elemets with loads applied are firts sliced a given number of times, and then, all t
parameters derived from the positions of the applied loads are included.

The positions where concentrated loads are applied are critical as there will be a
discontinuity, so a node must be added.

The positions where distributed loads start and end also introduce discontinuities, so we
also include nodes in those positions.
*/
func sliceLoadedElement(element *structure.Element, slices int) *Element {
	var (
		tPos  = sliceLoadedElementPositions(element.ConcentratedLoads, element.DistributedLoads, slices)
		nodes = makeNodesWithConcentratedLoads(element, tPos)
	)

	applyDistributedLoadsToNodes(nodes, element.DistributedLoads)

	return MakeElement(element, nodes)
}

/*
Computes all the t values where to slice an element with loads applied.

It starts by slicing the element a given number of times, and then adds all the load start
and end t values, removing any possible duplications.
*/
func sliceLoadedElementPositions(
	concentratedLoads []*load.ConcentratedLoad, 
	distributedLoads []*load.DistributedLoad,
	slices int,
) []inkgeom.TParam {
	tPos := inkgeom.SubTParamCompleteRangeTimes(slices)
	tPos = append(tPos, slicePositionsForConcentratedLoads(concentratedLoads)...)
	tPos = append(tPos, slicePositionsForDistributedLoads(distributedLoads)...)

	sort.Sort(inkgeom.ByTParamValue(tPos))

	var correctedTPos []inkgeom.TParam
	correctedTPos = append(correctedTPos, tPos[0])

	// FIXME: this might remove positions where a cocentrated load is applied, then,
	// the load will never be applied by the makeNodesWithConcentratedLoads function.
	for i := 1; i < len(tPos); i++ {
		if tPos[i-1].DistanceTo(tPos[i]) > minDistBetweenTSlices {
			correctedTPos = append(correctedTPos, tPos[i])
		}
	}

	return correctedTPos
}

/*
Collects all the concentrated loads t parameter value, provided the value
is not extreme, that is, `t != tMin` and `t != tMax`.
*/
func slicePositionsForConcentratedLoads(loads []*load.ConcentratedLoad) []inkgeom.TParam {
	var tVals []inkgeom.TParam

	for _, l := range loads {
		if !l.T.IsExtreme() {
			tVals = append(tVals, l.T)
		}
	}

	return tVals
}

/*
Collects all the distibutd loads start and end position t values, provided these
values are not extreme, that is, `t != tMin` and `t != tMax`.
*/
func slicePositionsForDistributedLoads(loads []*load.DistributedLoad) []inkgeom.TParam {
	var tVals []inkgeom.TParam

	for _, l := range loads {
		if !l.StartT.IsExtreme() {
			tVals = append(tVals, l.StartT)
		}

		if !l.EndT.IsExtreme() {
			tVals = append(tVals, l.EndT)
		}
	}

	return tVals
}

/*
Creates all the nodes for the given t positions and applies the concentrated loads
on those t positions where one is defined.

If the load is in global coordinates, its vector representation is projected
into the element's local reference frame.
*/
func makeNodesWithConcentratedLoads(element *structure.Element, tPos []inkgeom.TParam) []*Node {
	var (
		nodes = make([]*Node, len(tPos))
		elemRefFrame = element.Geometry.RefFrame()
	)

	for i, t := range tPos {
		node := MakeUnloadedNode(t, element.Geometry.PointAt(t))

		for _, load := range element.ConcentratedLoads {
			if t.Equals(load.T) {
				var localForces [3]float64

				if load.IsInLocalCoords {
					localForces = load.AsVector()
				} else {
					localForces = load.ProjectedVectorValue(elemRefFrame)
				}

				node.AddLocalExternalLoad(
					localForces[0],
					localForces[1],
					localForces[2],
				)
			}
		}

		nodes[i] = node
	}

	return nodes
}

func applyDistributedLoadsToNodes(nodes []*Node, loads []*load.DistributedLoad) {
	var trailNode, leadNode *Node

	for i, j := 0, 1; j < len(nodes); i, j = i+1, j+1 {
		trailNode, leadNode = nodes[i], nodes[j]

		for _, load := range loads {
			applyDistributedLoadToNodes(load, trailNode, leadNode)
		}
	}
}

/*
Applies a distribute load to the trailing and leading nodes in a finite element.

TODO: distribute Fx loads
TODO: distribute Mz loads
*/
func applyDistributedLoadToNodes(load *load.DistributedLoad, trailNode, leadNode *Node) {
	var (
		startLoad, endLoad = forceVectorInLocalCoords(load, trailNode, leadNode)
		length             = trailNode.DistanceTo(leadNode)
		halfLength         = 0.5 * length
		length2            = length * length
		length3            = length2 * length
		loadValueSlopes    = computeLoadValueSlopes(startLoad, endLoad, length)
	)

	var (
		trailFy       = (startLoad[1] * halfLength) + (3.0 * length2 * loadValueSlopes[1] / 20.0)
		trailFyMoment = (startLoad[1] * length2 / 12.0) + (length3 * loadValueSlopes[1] / 30.0)
	)
	trailNode.AddLocalLeftLoad(
		startLoad[0]*halfLength,
		trailFy,
		(startLoad[2]*halfLength)+trailFyMoment,
	)

	var (
		leadFy       = (startLoad[1] * halfLength) + (7.0 * length2 * loadValueSlopes[1] / 20.0)
		leadFyMoment = -(startLoad[1] * length2 / 12.0) - (length3 * loadValueSlopes[1] / 20.0)
	)
	leadNode.AddLocalRightLoad(
		startLoad[0]*halfLength,
		leadFy,
		(startLoad[2]*halfLength)+leadFyMoment,
	)
}

func forceVectorInLocalCoords(load *load.DistributedLoad, trailNode, leadNode *Node) (startLoad, endLoad [3]float64) {
	if load.IsInLocalCoords {
		startLoad = load.AsVectorAt(trailNode.T)
		endLoad = load.AsVectorAt(leadNode.T)
	} else {
		elementReferenceFrame := g2d.MakeRefFrameWithIVersor(g2d.MakeVectorFromTo(trailNode.Position, leadNode.Position))
		startLoad = load.ProjectedVectorAt(trailNode.T, elementReferenceFrame)
		endLoad = load.ProjectedVectorAt(leadNode.T, elementReferenceFrame)
	}

	return
}

func computeLoadValueSlopes(startLoad, endLoad [3]float64, length float64) [3]float64 {
	return [3]float64{
		(endLoad[0] - startLoad[0]) / length,
		(endLoad[1] - startLoad[1]) / length,
		(endLoad[2] - startLoad[2]) / length,
	}
}
