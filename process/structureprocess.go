/*
Package process defines the Finite Element Method computation.
It starts from the sliced structure, assembles the global system of equations, solves it
and creates a solution.
*/
package process

import (
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkmath"
	"github.com/angelsolaorbaiceta/inkmath/mat"
)

/*
Solve ...
*/
func Solve(s *preprocess.Structure) {
	c := make(chan preprocess.Element)

	for _, element := range s.Elements {
		go element.ComputeStiffnessMatrices(c)
	}

	sysMatrix /*, sysVector*/ := mat.MakeSparse(s.DofsCount, s.DofsCount) //, vec.Make(s.DofsCount)
	for i := 0; i < len(s.Elements); i++ {
		element := <-c
		addTermsToStiffnessMatrix(sysMatrix, &element)
	}
}

func addTermsToStiffnessMatrix(m mat.Matrixable, e *preprocess.Element) {
	var (
		stiffMat                    mat.Matrixable
		trailNodeDofs, leadNodeDofs [3]int
		dofs                        [6]int
		stiffVal                    float64
	)

	for i := 1; i < len(e.Nodes); i++ {
		stiffMat = e.GlobalStiffMatrixAt(i - 1)
		trailNodeDofs = e.Nodes[i-1].DegreesOfFreedomNum()
		leadNodeDofs = e.Nodes[i].DegreesOfFreedomNum()
		dofs = [6]int{
			trailNodeDofs[0], trailNodeDofs[1], trailNodeDofs[2],
			leadNodeDofs[0], leadNodeDofs[1], leadNodeDofs[2],
		}

		for i := 0; i < stiffMat.Rows(); i++ {
			for j := 0; j < stiffMat.Cols(); j++ {
				if stiffVal = stiffMat.Value(i, j); !inkmath.IsCloseToZero(stiffVal) {
					m.AddToValue(dofs[i], dofs[j], stiffVal)
				}
			}
		}
	}
}
