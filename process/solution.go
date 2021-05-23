package process

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Solution is the group of all element solutions with the structure metadata.
type Solution struct {
	Metadata *structure.StrMetadata
	Elements []*ElementSolution
}

/*
ElementCount returns the number of total bars in the structure's solution, which is the same number
as in the original definition of the structure.
*/
func (solution *Solution) ElementCount() int {
	return len(solution.Elements)
}
