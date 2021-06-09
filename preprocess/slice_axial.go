package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

/*
An axial element is an element which:

- is pinned in both ends and
- if it has loads, they are concentrated and applied to the ends of the element, and never include a moment about Z

Axial elements can be sliced using only it's end nodes. Axial elements deformation only happens in
their axial direction: tension or compression.
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
NetNodalLoadValues computes the net, locally projected loads at the start end (sFx & sFy) and at
the end end (eFx & eFy), assuming all loads are nodal (concentrated and applied to the ends of
the element).
*/
func netNodalLoadValues(
	loads []*load.ConcentratedLoad,
	localRefFrame g2d.RefFrame,
) (sFx, sFy, eFx, eFy float64) {
	var localLoadTorsor *math.Torsor

	for _, ld := range loads {
		if ld.IsInLocalCoords {
			localLoadTorsor = ld.AsTorsor()
		} else {
			localLoadTorsor = ld.AsTorsorProjectedTo(localRefFrame)
		}

		if ld.T.IsMin() {
			sFx += localLoadTorsor.Fx()
			sFy += localLoadTorsor.Fy()
		} else if ld.T.IsMax() {
			eFx += localLoadTorsor.Fx()
			eFy += localLoadTorsor.Fy()
		}
	}

	return
}
