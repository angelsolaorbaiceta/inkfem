package preprocess

import (
	"sort"

	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

// Minimum distance between two consecutive t values in the slices.
const minDistBetweenTSlices = 1e-3

// Elemets with loads applied are firts sliced a given number of times, and then, all t parameters
// derived from the positions of the applied loads are included.
//
// The positions where concentrated loads are applied are critical as there will be a discontinuity,
// so a node must be added.
//
// The positions where distributed loads start and end also introduce discontinuities, so we also
// include nodes in those positions.
func sliceLoadedElement(element *structure.Element, slices int) *Element {
	var (
		tPos  = sliceLoadedElementPositions(element.ConcentratedLoads, element.DistributedLoads, slices)
		nodes = makeNodesWithConcentratedLoads(element, tPos)
	)

	applyDistributedLoadsToNodes(nodes, element.DistributedLoads)

	return MakeElement(element, nodes)
}

// Computes all the t values where to slice an element with loads applied.
//
// It starts by slicing the element a given number of times, and then adds all the load start and
// end t values, removing any possible duplications.
func sliceLoadedElementPositions(
	concentratedLoads []*load.ConcentratedLoad,
	distributedLoads []*load.DistributedLoad,
	slices int,
) []nums.TParam {
	tPos := nums.SubTParamCompleteRangeTimes(slices)
	tPos = append(tPos, slicePositionsForConcentratedLoads(concentratedLoads)...)
	tPos = append(tPos, slicePositionsForDistributedLoads(distributedLoads)...)

	sort.Sort(nums.ByTParamValue(tPos))

	var correctedTPos []nums.TParam
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

// SlicePositionsForConcentratedLoads collects all the concentrated loads t parameter value, provided
// the value is not extreme, that is, `t != tMin` and `t != tMax`.
func slicePositionsForConcentratedLoads(loads []*load.ConcentratedLoad) []nums.TParam {
	var tVals []nums.TParam

	for _, load := range loads {
		if !load.T.IsExtreme() {
			tVals = append(tVals, load.T)
		}
	}

	return tVals
}

// SlicePositionsForDistributedLoads collects all the distibutd loads start and end position t values,
// provided these values are not extreme, that is, `t != tMin` and `t != tMax`.
func slicePositionsForDistributedLoads(loads []*load.DistributedLoad) []nums.TParam {
	var tVals []nums.TParam

	for _, load := range loads {
		if !load.StartT.IsExtreme() {
			tVals = append(tVals, load.StartT)
		}

		if !load.EndT.IsExtreme() {
			tVals = append(tVals, load.EndT)
		}
	}

	return tVals
}

// MakeNodesWithConcentratedLoads creates all the nodes for the given t positions and applies the
// concentrated loads on those t positions where one is defined.
//
// If the load is in global coordinates, its vector representation is projected into the element's
// local reference frame.
func makeNodesWithConcentratedLoads(element *structure.Element, tPos []nums.TParam) []*Node {
	var (
		nodes        = make([]*Node, len(tPos))
		elemRefFrame = element.RefFrame()
	)

	for i, t := range tPos {
		node := MakeUnloadedNode(t, element.PointAt(t))

		for _, load := range element.ConcentratedLoads {
			if t.Equals(load.T) {
				var localLoadTorsor *math.Torsor

				if load.IsInLocalCoords {
					localLoadTorsor = load.AsTorsor()
				} else {
					localLoadTorsor = load.AsTorsorProjectedTo(elemRefFrame)
				}

				node.AddLocalExternalLoad(localLoadTorsor)
			}
		}

		nodes[i] = node
	}

	return nodes
}
