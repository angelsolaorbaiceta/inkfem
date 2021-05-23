package process

import (
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
)

/*
Solve assembles the system of equations for the structure and solves it using the Preconditioned
Conjugate Gradient numerical procedure.

The local element stresses are computed using the displacements obtained in the first step.
*/
func Solve(structure *preprocess.Structure, options SolveOptions) *Solution {
	globalDisplacements := computeGlobalDisplacements(structure, options)

	var (
		elementSolution  *ElementSolution
		elementSolutions = make([]*ElementSolution, structure.ElementsCount())
	)

	log.StartComputeStresses()
	for i, element := range structure.Elements {
		elementSolution = MakeElementSolution(element)
		elementSolution.SolveUsingDisplacements(globalDisplacements)
		elementSolutions[i] = elementSolution
	}
	log.EndComputeStresses()

	return &Solution{Metadata: &structure.Metadata, Elements: elementSolutions}
}
