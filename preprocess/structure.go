package preprocess

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Structure result of preprocessing original structure, ready to be solved.
// The elements of a preprocessed structure are already sliced.
type Structure struct {
	Metadata  structure.StrMetadata
	nodes     map[contracts.StrID]*structure.Node
	Elements  []*Element
	DofsCount int
}

// NodesCount returns the number of nodes in the original structure.
func (s *Structure) NodesCount() int {
	return len(s.nodes)
}

// GetNodeById returns the node with the given id. Panics if the node is doesn't exist in the structure.
func (s *Structure) GetNodeById(id contracts.StrID) *structure.Node {
	if node, exists := s.nodes[id]; exists {
		return node
	}

	panic(fmt.Sprintf("Can't find node with id: %s", id))
}

// GetElementNodes returns the element's start and end nodes.
func (s *Structure) GetElementNodes(element *Element) (*structure.Node, *structure.Node) {
	return s.GetNodeById(element.StartNodeID()), s.GetNodeById(element.EndNodeID())
}

// GetAllNodes returns a slice containing all of the structure nodes.
func (s *Structure) GetAllNodes() []*structure.Node {
	nodes := make([]*structure.Node, 0, s.NodesCount())

	for _, node := range s.nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

// NodesById is a map where the nodes of the structure can be accessed by their id.
func (s *Structure) NodesById() map[contracts.StrID]*structure.Node {
	return s.nodes
}

// ElementsCount returns the number of elements in the original structure.
func (s *Structure) ElementsCount() int {
	return len(s.Elements)
}
