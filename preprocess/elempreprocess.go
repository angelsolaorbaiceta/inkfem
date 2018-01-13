package preprocess

import (
	"sort"
	"sync"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
)

// DoElement preprocesses the given structural element subdividing it as corresponds.
// The result is sent through a channel.
func DoElement(e structure.Element, c chan Element, wg *sync.WaitGroup) {
	defer wg.Done()

	if e.IsAxialMember() {
		c <- sliceAxialElement(e)
	} else if e.HasLoadsApplied() {
		c <- sliceLoadedElement(e, 18)
	} else {
		c <- sliceUnloadedElement(e, 12)
	}
}

/* <---------- Non Sliced ----------> */

/*
An axial element is an element which:
    - is pinned in both ends and
    - if it has loads, they are concentrated and applied to the ends of the element, and never include a moment about Z

Axial elements can be sliced using only it's end nodes. Axial elements deformation
only happens in their axial direction: tension or compression.
*/
func sliceAxialElement(e structure.Element) Element {
	if e.HasLoadsApplied() {
		sFx, sFy, eFx, eFy := netNodalLoadValues(e.Loads, e.Geometry.RefFrame())

		return MakeElement(
			e,
			[]Node{
				MakeNode(inkgeom.MIN_T, e.StartPoint(), sFx, sFy, 0.0),
				MakeNode(inkgeom.MAX_T, e.EndPoint(), eFx, eFy, 0.0),
			})
	}

	return MakeElement(
		e,
		[]Node{
			MakeUnloadedNode(inkgeom.MIN_T, e.StartPoint()),
			MakeUnloadedNode(inkgeom.MAX_T, e.EndPoint()),
		})
}

/*
Assuming all loads are nodal (concentrated and applied to the ends of the element), computes
the net, locally projected loads at the start end (sFx & sFy) and at the end end (eFx & eFy).
*/
func netNodalLoadValues(loads []load.Load, localRefFrame inkgeom.RefFrame) (sFx, sFy, eFx, eFy float64) {
	var localForcesVector inkgeom.Projectable

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

/* <---------- Sliced : Loaded ----------> */

/*
Elemets with loads applied are firts sliced a given number of times, and then, all t parameters
derived from the positions of the applied loads are included.

The positions where concentrated loads are applied are critical as there will be a discontinuity,
so a node must be added.

The positions where distributed loads start and end also introduce discontinuities, so we also
include nodes in those positions.
*/
func sliceLoadedElement(e structure.Element, times int) Element {
	tPos := sliceLoadedElementPositions(e.Loads, times)
	var (
		trailNode, leadNode *Node
		length, halfLength  float64
		avgLoadValVect      [3]float64
		nodes               = make([]Node, len(tPos))
	)

	// Nodes creation & concentrated laod application
	for i, t := range tPos {
		// TODO: add concentrated loads after projecting them
		nodes[i] = MakeUnloadedNode(t, e.Geometry.PointAt(t))
	}

	// Distributed load application
	for i, j := 0, 1; j < len(tPos); i, j = i+1, j+1 {
		trailNode, leadNode = &nodes[i], &nodes[j]
		length = e.Geometry.LengthBetween(trailNode.T, leadNode.T)
		halfLength = 0.5 * length

		for _, load := range e.Loads {
			avgLoadValVect = load.AvgValueVectorBetween(trailNode.T, leadNode.T)
			// TODO: basis change

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

	return MakeElement(e, nodes)
}

/*
Computes all the t values where to slice an element with loads applied.
It starts by slicing the element a given number of times, and then adds all the load
start and end t values, removing any possible duplications.
*/
func sliceLoadedElementPositions(loads []load.Load, times int) []inkgeom.TParam {
	tPos := append(inkgeom.SubTParamCompleteRangeTimes(times), tValsForLoadApplications(loads)...)
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

/* <---------- Sliced : Unloaded ----------> */

/*
Non axial elements which have no loads applied are sliced just by subdividint their geometry
a given number of times, so that the slices have the same length.
*/
func sliceUnloadedElement(e structure.Element, times int) Element {
	tPos := inkgeom.SubTParamCompleteRangeTimes(times)
	nodes := make([]Node, len(tPos))

	for i := 0; i < len(tPos); i++ {
		nodes[i] = MakeUnloadedNode(tPos[i], e.PointAt(tPos[i]))
	}

	return MakeElement(e, nodes)
}
