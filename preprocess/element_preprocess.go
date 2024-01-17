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
// The result is sent through a channel.
//
// Depending on the nature of the bar, it is sliced differently:
// axial bars aren't sliced at all; bars without loads are sliced into 6 elements;
// bars with loads are sliced into 10 elements.
func sliceElement(element *structure.Element, c chan<- *Element) {
	if element.IsAxialMember() {
		c <- sliceAxialElement(element)
	} else if element.HasLoadsApplied() {
		c <- sliceLoadedElement(element, elementWithLoadsSlices)
	} else {
		c <- sliceElementWithoutLoads(element, elementWithoutLoadsSlices)
	}
}
