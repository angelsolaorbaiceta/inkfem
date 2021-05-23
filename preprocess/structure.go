package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

/*
Structure result of preprocessing original structure, ready to be solved.
The elements of a preprocessed structure are already sliced.
*/
type Structure struct {
	Metadata  structure.StrMetadata
	Nodes     map[contracts.StrID]*structure.Node
	Elements  []*Element
	DofsCount int
}

// NodesCount returns the number of nodes in the original structure.
func (s *Structure) NodesCount() int {
	return len(s.Nodes)
}

// ElementsCount returns the number of elements in the original structure.
func (s *Structure) ElementsCount() int {
	return len(s.Elements)
}
