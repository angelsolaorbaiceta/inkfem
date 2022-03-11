package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
)

// NodesById is a composable map of nodes with some useful methods.
type NodesById struct {
	nodes map[contracts.StrID]*Node
}

func MakeNodesById(nodes map[contracts.StrID]*Node) NodesById {
	return NodesById{nodes: nodes}
}

// NodesCount is the number of nodes in the structure.
func (n *NodesById) NodesCount() int {
	return len(n.nodes)
}

// GetNodeById returns the node with the given id.
// Panics if the node is doesn't exist in the structure.
func (n *NodesById) GetNodeById(id contracts.StrID) *Node {
	if node, exists := n.nodes[id]; exists {
		return node
	}

	panic(fmt.Sprintf("Can't find node with id: %s", id))
}

// GetAllNodes returns a slice containing all of the structure nodes.
func (n *NodesById) GetAllNodes() []*Node {
	nodes := make([]*Node, 0, n.NodesCount())

	for _, node := range n.nodes {
		nodes = append(nodes, node)
	}

	return nodes
}

// NodesById is a map where the nodes of the structure can be accessed by their id.
func (n *NodesById) NodesById() map[contracts.StrID]*Node {
	return n.nodes
}
