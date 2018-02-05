/*
Package process defines the Finite Element Method computation.
It starts from the sliced structure, assembles the global system of equations, solves it
and creates a solution.
*/
package process

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath"
	"github.com/angelsolaorbaiceta/inkmath/lineq"
	"github.com/angelsolaorbaiceta/inkmath/mat"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

/*
Solve ...
*/
func Solve(s *preprocess.Structure) {
	sysMatrix, sysVector := makeSystemOfEqs(s)
	solver := lineq.ConjugateGradientSolver{MaxError: 1e-5, MaxIter: sysVector.Length()}
	if !solver.CanSolve(sysMatrix, sysVector) {
		panic("Solver cannot solve system!")
	}
	displSolutions := solver.Solve(sysMatrix, sysVector)
	fmt.Println(displSolutions)
}

func makeSystemOfEqs(s *preprocess.Structure) (mat.Matrixable, *vec.Vector) {
	c := make(chan preprocess.Element)

	for _, element := range s.Elements {
		go element.ComputeStiffnessMatrices(c)
	}

	sysMatrix, sysVector := mat.MakeSparse(s.DofsCount, s.DofsCount), vec.Make(s.DofsCount)
	for i := 0; i < len(s.Elements); i++ {
		element := <-c
		addTermsToStiffnessMatrix(sysMatrix, &element)
		addTermsToLoadVector(sysVector, &element)
	}
	addDispConstraints(sysMatrix, sysVector, s.Nodes)

	return sysMatrix, sysVector
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

func addDispConstraints(m mat.Matrixable, v *vec.Vector, nodes map[int]structure.Node) {
	var (
		constraint *structure.Constraint
		dofs       [3]int
	)

	addConstraintAtDof := func(dof int) {
		m.SetZeroCol(dof)
		m.SetIdentityRow(dof)
		v.SetZero(dof)
	}

	for _, node := range nodes {
		if node.IsExternallyConstrained() {
			constraint = node.ExternalConstraint
			dofs = node.DegreesOfFreedomNum()

			if !constraint.AllowsDispX() {
				addConstraintAtDof(dofs[0])
			}
			if !constraint.AllowsDispY() {
				addConstraintAtDof(dofs[1])
			}
			if !constraint.AllowsRotation() {
				addConstraintAtDof(dofs[2])
			}
		}
	}
}

func addTermsToLoadVector(v *vec.Vector, e *preprocess.Element) {
	var (
		localActions [3]float64
		globalForces inkgeom.Projectable
		dofs         [3]int
		refFrame     = e.Geometry().RefFrame()
	)

	for _, node := range e.Nodes {
		localActions = node.LocalActions()
		globalForces = refFrame.ProjectionsToGlobal(localActions[0], localActions[1])
		dofs = node.DegreesOfFreedomNum()

		v.SetValue(dofs[0], globalForces.X)
		v.SetValue(dofs[1], globalForces.Y)
		v.SetValue(dofs[2], localActions[2])
	}
}
