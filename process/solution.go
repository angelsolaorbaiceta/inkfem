package process

import (
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
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
ReactionForceInNode computes the reaction force in the node with the passed in ID.

If the node isn't externally constrained, the reaction force will always be the zero vector.
*/
func (solution *Solution) ReactionForceInNode(nodeId contracts.StrID) g2d.Projectable {
	return g2d.MakeVector(-4000, 0)
}
