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
func Solve(s *preprocess.Structure, options SolveOptions) *Solution {
	log.StartSolve()

	globalDisplacements := computeGlobalDisplacements(s, options)

	var (
		elementSolution  *ElementSolution
		elementSolutions = make([]*ElementSolution, len(s.Elements))
	)

	for i, element := range s.Elements {
		elementSolution = MakeElementSolution(element)
		elementSolution.SolveUsingDisplacements(globalDisplacements)
		elementSolutions[i] = elementSolution
	}

	return &Solution{Metadata: &s.Metadata, Elements: elementSolutions}
}
