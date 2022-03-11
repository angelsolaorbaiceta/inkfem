package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Structure result of preprocessing original structure, ready to be solved.
// The elements of a preprocessed structure are already sliced.
type Structure struct {
	Metadata structure.StrMetadata
	structure.NodesById
	Elements  []*Element
	DofsCount int
}

// GetElementNodes returns the element's start and end nodes.
func (s *Structure) GetElementNodes(element *Element) (*structure.Node, *structure.Node) {
	return s.GetNodeById(element.StartNodeID()), s.GetNodeById(element.EndNodeID())
}

// ElementsCount returns the number of elements in the original structure.
func (s *Structure) ElementsCount() int {
	return len(s.Elements)
}
