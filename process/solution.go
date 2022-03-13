package process

import (
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Solution is the group of all element solutions with the structure metadata.
type Solution struct {
	Metadata structure.StrMetadata
	structure.NodesById
	Elements []*ElementSolution
}

// ElementCount returns the number of total bars in the structure's solution, which is the same number
// as in the original definition of the structure.
func (solution *Solution) ElementCount() int {
	return len(solution.Elements)
}

// NodeReactions returns the map of externally constrained nodes with their reaction.
func (solution *Solution) NodeReactions() map[contracts.StrID]*math.Torsor {
	nodeReactions := make(map[contracts.StrID]*math.Torsor)

	for _, node := range solution.GetAllNodes() {
		if node.IsExternallyConstrained() {
			nodeReactions[node.GetID()] = solution.reactionInNode(node.GetID())
		}
	}

	return nodeReactions
}

// ReactionInNode computes the reaction torsor {fx, fy, mz} in the node with the passed
// in ID in global coordinates.
//
// If the node isn't externally constrained, the reaction will always be a nil torsor {0, 0, 0}.
// If the structure contains no node with the given ID, it'll panic.
func (solution *Solution) reactionInNode(nodeId contracts.StrID) *math.Torsor {
	var (
		node     = solution.GetNodeById(nodeId)
		reaction = math.MakeNilTorsor()
	)

	if !node.IsExternallyConstrained() {
		return reaction
	}

	for _, element := range solution.Elements {
		if element.StartNodeID() == nodeId {
			reaction = reaction.Plus(element.GlobalStartTorsor())
		} else if element.EndNodeID() == nodeId {
			reaction = reaction.Plus(element.GlobalEndTorsor())
		}
	}

	return reaction
}
