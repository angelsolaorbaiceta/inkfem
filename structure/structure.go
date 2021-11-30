package structure

import "github.com/angelsolaorbaiceta/inkfem/contracts"

// A Structure is a group of linear resistant elements joined together designed to withstand the
// application of external loads –concentrated and distributed.
type Structure struct {
	Metadata StrMetadata
	Nodes    map[contracts.StrID]*Node
	Elements []*Element
}

// NodesCount is the number of nodes in the structure.
func (s *Structure) NodesCount() int {
	return len(s.Nodes)
}

// ElementsCount is the number of elements in the structure.
func (s *Structure) ElementsCount() int {
	return len(s.Elements)
}
