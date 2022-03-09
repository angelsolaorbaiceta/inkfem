package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
)

// A Structure is a group of linear resistant elements joined together designed to withstand the
// application of external loads, concentrated and distributed.
type Structure struct {
	Metadata StrMetadata
	nodes    map[contracts.StrID]*Node
	Elements []*Element
}

// Make creates a new structure model.
func Make(metadata StrMetadata, nodes map[contracts.StrID]*Node, elements []*Element) *Structure {
	return &Structure{
		Metadata: metadata,
		nodes:    nodes,
		Elements: elements,
	}
}

// NodesCount is the number of nodes in the structure.
func (s *Structure) NodesCount() int {
	return len(s.nodes)
}

// GetNodeById returns the node with the given id.
// Panics if the node is doesn't exist in the structure.
func (s *Structure) GetNodeById(id contracts.StrID) *Node {
	if node, exists := s.nodes[id]; exists {
		return node
	}

	panic(fmt.Sprintf("Can't find node with id: %s", id))
}

// GetAllNodes returns a slice containing all of the structure nodes.
func (s *Structure) GetAllNodes() []*Node {
	nodes := make([]*Node, 0, s.NodesCount())

	for _, node := range s.nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

// NodesById is a map where the nodes of the structure can be accessed by their id.
func (s *Structure) NodesById() map[contracts.StrID]*Node {
	return s.nodes
}

// ElementsCount is the number of elements in the structure.
func (s *Structure) ElementsCount() int {
	return len(s.Elements)
}
