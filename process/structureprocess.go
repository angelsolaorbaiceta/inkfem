/*
Package process defines the Finite Element Method computation.
It starts from the sliced structure, assembles the global system of equations, solves it
and creates a solution.
*/
package process

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkmath/mat"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

/*
Solve ...
*/
func Solve(s *preprocess.Structure) {
	c := make(chan preprocess.Element)

	for _, element := range s.Elements {
		fmt.Printf("*** Element [%d]\n", element.ID())
		go element.ComputeStiffnessMatrices(c)
	}

	sysMatrix, sysVector := mat.MakeSparse(s.DofsCount, s.DofsCount), vec.Make(s.DofsCount)
	for i := 0; i < len(s.Elements); i++ {
		element := <-c
		fmt.Printf("-- Got element: %d\n", element.ID())
		addTermsToStiffnessMatrix(sysMatrix, sysVector, &element)
	}
}

func addTermsToStiffnessMatrix(m mat.Matrixable, v *vec.Vector, e *preprocess.Element) {
	// fmt.Printf("Adding terms to matrix for element %d\n", e.ID())
}
