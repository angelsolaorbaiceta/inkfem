package process

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Solution is the group of all element solutions with the structure metadata.
type Solution struct {
	Metadata *structure.StrMetadata
	Nodes    map[contracts.StrID]*structure.Node
	Elements []*ElementSolution
}

/*
ElementCount returns the number of total bars in the structure's solution, which is the same number
as in the original definition of the structure.
*/
func (solution *Solution) ElementCount() int {
	return len(solution.Elements)
}

/*
ReactionInNode computes the reaction torsor {fx, fy, mz} in the node with the passed in ID
in global coordinates.

If the node isn't externally constrained, the reaction will always be a nil torsor {0, 0, 0}.
If the structure contains no node with the given ID, it'll panic.
*/
func (solution *Solution) ReactionInNode(nodeId contracts.StrID) *math.Torsor {
	node, hasNode := solution.Nodes[nodeId]
	if !hasNode {
		panic(fmt.Sprintf("Structure doesn't contain a node with id: '%v'", nodeId))
	}

	reaction := math.MakeNilTorsor()

	if !node.IsExternallyConstrained() {
		return reaction
	}

	for _, element := range solution.Elements {
		if element.StartNodeID == nodeId {
			reaction = reaction.Minus(element.GlobalStartTorsor())
		} else if element.EndNodeID == nodeId {
			reaction = reaction.Minus(element.GlobalEndTorsor())
		}
	}

	return reaction
}
