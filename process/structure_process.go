package process

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
)

/*
Solve assembles the system of equations for the structure and solves it using
the Preconditioned Conjugate Gradient numerical procedure.

Using the displacements obtained from the solution of the system's solution,
the local stresses are computed.
*/
func Solve(s *preprocess.Structure, options SolveOptions) *Solution {
	globalDisplacements := computeGlobalDisplacements(s, options)

	var (
		elementSolution  *ElementSolution
		elementSolutions = make([]*ElementSolution, len(s.Elements))
	)

	for i, element := range s.Elements {
		elementSolution = MakeElementSolution(&element)
		elementSolution.SetDisplacements(globalDisplacements)
		elementSolution.ComputeStresses()

		fmt.Printf("[elementSolution] -> %p, %s\n", elementSolution, elementSolution.OriginalElementString())

		elementSolutions[i] = elementSolution
	}
	fmt.Printf("[list] -> %v\n", elementSolutions)

	return &Solution{Metadata: &s.Metadata, Elements: elementSolutions}
}
