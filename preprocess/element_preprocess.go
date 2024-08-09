package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

const (
	elementWithLoadsSlices    = 10
	elementWithoutLoadsSlices = 6
	// Minimum distance between two consecutive t values in the slices.
	minDistBetweenTSlices = 1e-3
)

// sliceElement slices the given bar into finite elements.
// The algorithm ensures that between two nodes, there's always a minimum
// distance of 0.003 between the t values of the slices.
// The result is sent through a channel.
//
// Depending on the nature of the bar, it is sliced differently:
//
//   - Axial bars aren't sliced at all.
//   - Bars without loads are sliced into 6 elements.
//   - Bars with loads are sliced into 10 elements.
//
// Intermediate points (not end nodes) where a concentrated load is applied,
// also generate intermediate nodes for the load to be included.
func sliceElement(element *structure.Element, c chan<- *Element) {
	if element.IsAxialMember() {
		c <- sliceAxialElement(element)
	} else if element.HasLoadsApplied() {
		c <- sliceLoadedElement(element, elementWithLoadsSlices)
	} else {
		c <- sliceElementWithoutLoads(element, elementWithoutLoadsSlices)
	}
}
