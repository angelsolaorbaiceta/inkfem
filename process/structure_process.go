package process

import (
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
)

/*
Solve assembles the system of equations for the structure and solves it using
the Preconditioned Conjugate Gradient numerical procedure.

Using the displacements obtained from the solution of the system's solution,
the local stresses are computed.
*/
func Solve(structure *preprocess.Structure, options SolveOptions) *Solution {
	globalDisplacements := computeGlobalDisplacements(structure, options)

	var (
		elementSolution  *ElementSolution
		elementSolutions = make([]*ElementSolution, len(structure.Elements))
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
