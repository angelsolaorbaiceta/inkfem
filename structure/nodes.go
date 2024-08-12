package structure

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
)

// NodesById is a composable map of nodes with some useful methods.
type NodesById struct {
	nodes         NodesByIdMap
	allNodesCache []*Node
}

func MakeNodesById(nodes NodesByIdMap) NodesById {
	return NodesById{nodes: nodes, allNodesCache: nil}
}

func (e *NodesById) Copy() NodesById {
	nodes := make(NodesByIdMap)

	for id, node := range e.nodes {
		nodes[id] = node.Copy()
	}

	return MakeNodesById(nodes)
}

// NodesCount is the number of nodes in the structure.
func (n *NodesById) NodesCount() int {
	return len(n.nodes)
}

// ConstrainedNodesCount is the number of nodes with an external constraint.
func (n *NodesById) ConstrainedNodesCount() int {
	count := 0

	for _, node := range n.nodes {
		if node.IsExternallyConstrained() {
			count++
		}
	}

	return count
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
	if n.allNodesCache == nil {
		nodes := make([]*Node, 0, n.NodesCount())

		for _, node := range n.nodes {
			nodes = append(nodes, node)
		}

		n.allNodesCache = nodes
	}

	return n.allNodesCache
}

// NodesById is a map where the nodes of the structure can be accessed by their id.
func (n *NodesById) NodesById() NodesByIdMap {
	return n.nodes
}
