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
		c <- sliceUnloadedElement(e, 12)
	} else {
		c <- sliceLoadedElement(e, 18)
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
func sliceLoadedElement(e structure.Element, times int) Element {
	// tPos := sliceLoadedElementPositions(e.Loads, times)

	var (
		// trailT, leadT inkgeom.TParam
		nodes []Node
	)

	// for i, j := 0, 1; j < len(tPos); i, j = i+1, j+1 {
	// 	trailT, leadT = tPos[i], tPos[j]
	//
	// 	if inkmath.FuzzyEqualEps(trailT.Value(), leadT.Value(), 1e-5) {
	// 		continue
	// 	}
	//
	// 	if len(nodes) == 0 {
	// 		// apply concentrated  start node load here
	// 		nodes = append(nodes, MakeUnloadedNode(trailT, e.PointAt(trailT)))
	// 	}
	// 	nodes = append(nodes, MakeUnloadedNode(trailT, e.PointAt(leadT)))
	//
	// 	for _, ld := range e.Loads {
	// 		if ld.IsConcentrated() && inkmath.IsCloseToZero(ld.T().DistanceTo(leadT)) {
	// 			// TODO: projected in local
	// 			nodes[len(nodes)-1].AddLoad(ld.VectorValue())
	// 		} else {
	//
	// 		}
	// 	}
	// }

	return MakeElement(e, nodes)
}

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
func sliceUnloadedElement(e structure.Element, times int) Element {
	tPos := inkgeom.SubTParamCompleteRangeTimes(times)
	nodes := make([]Node, len(tPos))

	for i := 0; i < len(tPos); i++ {
		nodes[i] = MakeUnloadedNode(tPos[i], e.PointAt(tPos[i]))
	}

	return MakeElement(e, nodes)
}
