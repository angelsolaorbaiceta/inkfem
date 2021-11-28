package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

/*
Non axial elements which have no loads applied are sliced just by subdividing their
geometry into a given number of slices, so that the slices have the same length.
*/
func sliceElementWithoutLoads(element *structure.Element, slices int) *Element {
	if element.HasLoadsApplied() {
		panic("Expected an element without external loads")
	}

	var (
		tPos  = nums.SubTParamCompleteRangeTimes(slices)
		nodes = make([]*Node, len(tPos))
	)

	for i := 0; i < len(tPos); i++ {
		nodes[i] = MakeUnloadedNode(tPos[i], element.PointAt(tPos[i]))
	}

	return MakeElement(element, nodes)
}
