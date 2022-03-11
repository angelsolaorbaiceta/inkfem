package structure

import (
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
