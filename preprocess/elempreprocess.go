package preprocess

import (
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
	} else {
		c <- sliceElement(e, 12)
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

/* <---------- Sliced ----------> */
func sliceElement(e structure.Element, times int) Element {
	tPos := inkgeom.SubTParamCompleteRangeTimes(times)
	loadTPos := tValsForLoadApplications(e.Loads)

	nodes := make([]Node, len(tPos))
	for i := 0; i < len(tPos); i++ { // TODO: add loads
		nodes[i] = MakeUnloadedNode(tPos[i], e.PointAt(tPos[i]))
	}

	return MakeElement(e, nodes)
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
