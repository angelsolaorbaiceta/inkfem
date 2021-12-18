package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

const (
	elementWithLoadsSlices    = 10
	elementWithoutLoadsSlices = 6
)

// elementModel preprocesses the given structural element subdividing it as corresponds.
// The result is sent through a channel.
func elementModel(element *structure.Element, c chan<- *Element) {
	if element.IsAxialMember() {
		c <- sliceAxialElement(element)
	} else if element.HasLoadsApplied() {
		c <- sliceLoadedElement(element, elementWithLoadsSlices)
	} else {
		c <- sliceElementWithoutLoads(element, elementWithoutLoadsSlices)
	}
}
