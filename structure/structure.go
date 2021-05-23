package structure

import "github.com/angelsolaorbaiceta/inkfem/contracts"

/*
Structure is a group of resistant elements joined together designed to withstand the application of
external loads.
*/
type Structure struct {
	Metadata StrMetadata
	Nodes    map[contracts.StrID]*Node
	Elements []*Element
}

// NodesCount returns the number of nodes in the structure.
func (s Structure) NodesCount() int {
	return len(s.Nodes)
}

// ElementsCount returns the number of elements in the structure.
func (s Structure) ElementsCount() int {
	return len(s.Elements)
}
