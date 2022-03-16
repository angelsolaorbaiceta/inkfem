package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
)

// A Structure is a group of linear resistant elements joined together designed to withstand the
// application of external loads, concentrated and distributed.
type Structure struct {
	Metadata StrMetadata
	NodesById
	Elements []*Element
}

// Make creates a new structure model.
func Make(metadata StrMetadata, nodes map[contracts.StrID]*Node, elements []*Element) *Structure {
	return &Structure{
		Metadata:  metadata,
		NodesById: NodesById{nodes: nodes},
		Elements:  elements,
	}
}

// ElementsCount is the number of elements in the structure.
func (s *Structure) ElementsCount() int {
	return len(s.Elements)
}

// GetElementById returns the element with the given id or panics.
// This operation has an O(n) time complexity as it needs to iterate over all elements.
func (s *Structure) GetElementById(id contracts.StrID) *Element {
	for _, element := range s.Elements {
		if element.GetID() == id {
			return element
		}
	}

	panic(fmt.Sprintf("Can't find element with id %s", id))
}
