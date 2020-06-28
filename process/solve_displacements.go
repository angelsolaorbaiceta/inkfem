package process

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath/lineq"
	"github.com/angelsolaorbaiceta/inkmath/mat"
	"github.com/angelsolaorbaiceta/inkmath/nums"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

/*
Computes the structure's global displacements given the preprocessed structure.

The process involves generating the structure's system of equations and solving it
using the Preconditioned Conjugate Gradiend numerical procedure.
*/
func computeGlobalDisplacements(
	structure *preprocess.Structure,
	options SolveOptions,
) *vec.Vector {
	if options.Verbose {
		fmt.Println("> assembling system of equations...")
	}

	sysMatrix, sysVector := makeSystemOfEquations(structure)
	if options.Verbose {
		fmt.Printf("[DONE] assembled system with %d equations\n", sysVector.Length())
	}

	if options.SaveSysMatrixImage {
		go mat.ToImage(sysMatrix, options.OutputPath)
	}

	if options.Verbose {
		fmt.Println("> solving sytem of equations for global displacements")
	}

	solver := lineq.PreconditionedConjugateGradientSolver{
		MaxError: options.MaxDisplacementsError,
		MaxIter:  sysVector.Length(),
	}
	if options.SafeChecks && !solver.CanSolve(sysMatrix, sysVector) {
		panic("Solver cannot solve system!")
	}

	globalDispSolution := solver.Solve(sysMatrix, sysVector)
	if options.Verbose {
		fmt.Printf(
			"[DONE] solved system in %d iterations, error = %f\n",
			globalDispSolution.IterCount,
			globalDispSolution.MinError,
		)
	}

	return globalDispSolution.Solution
}

/*
Generates the system of equations matrix and vector from the preprocessed structure.

It computes each of the sliced element's stiffness matrices and assembles them into one
global matrix. It also assembles the global loads vector from the sliced element nodes.
*/
func makeSystemOfEquations(
	structure *preprocess.Structure,
) (mat.ReadOnlyMatrix, *vec.Vector) {
	c := make(chan preprocess.Element)

	for _, element := range structure.Elements {
		go element.ComputeStiffnessMatrices(c)
	}

	var (
		sysMatrix = mat.MakeSparse(structure.DofsCount, structure.DofsCount)
		sysVector = vec.Make(structure.DofsCount)
	)

	for i := 0; i < len(structure.Elements); i++ {
		element := <-c
		addTermsToStiffnessMatrix(sysMatrix, &element)
		addTermsToLoadVector(sysVector, &element)
	}
	addDispConstraints(sysMatrix, sysVector, &structure.Nodes)

	return sysMatrix, sysVector
}

func addTermsToStiffnessMatrix(matrix mat.MutableMatrix, element *preprocess.Element) {
	var (
		stiffMat                    mat.ReadOnlyMatrix
		trailNodeDofs, leadNodeDofs [3]int
		dofs                        [6]int
		stiffVal                    float64
	)

	for i := 1; i < len(element.Nodes); i++ {
		stiffMat = element.GlobalStiffMatrixAt(i - 1)
		trailNodeDofs = element.Nodes[i-1].DegreesOfFreedomNum()
		leadNodeDofs = element.Nodes[i].DegreesOfFreedomNum()
		dofs = [6]int{
			trailNodeDofs[0], trailNodeDofs[1], trailNodeDofs[2],
			leadNodeDofs[0], leadNodeDofs[1], leadNodeDofs[2],
		}

		// TODO: this i is dangerous: shadows earlier def of i
		for i := 0; i < stiffMat.Rows(); i++ {
			for j := 0; j < stiffMat.Cols(); j++ {
				if stiffVal = stiffMat.Value(i, j); !nums.IsCloseToZero(stiffVal) {
					matrix.AddToValue(dofs[i], dofs[j], stiffVal)
				}
			}
		}
	}
}

func addDispConstraints(
	matrix mat.MutableMatrix,
	vector *vec.Vector,
	nodes *map[int]*structure.Node,
) {
	var (
		constraint structure.Constraint
		dofs       [3]int
	)

	addConstraintAtDof := func(dof int) {
		matrix.SetZeroCol(dof)
		matrix.SetIdentityRow(dof)
		vector.SetZero(dof)
	}

	for _, node := range *nodes {
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

func addTermsToLoadVector(vector *vec.Vector, element *preprocess.Element) {
	var (
		localActions [3]float64
		globalForces inkgeom.Projectable
		dofs         [3]int
		refFrame     = element.Geometry().RefFrame()
	)

	for _, node := range element.Nodes {
		localActions = node.LocalActions()
		globalForces = refFrame.ProjectionsToGlobal(localActions[0], localActions[1])
		dofs = node.DegreesOfFreedomNum()

		vector.SetValue(dofs[0], globalForces.X)
		vector.SetValue(dofs[1], globalForces.Y)
		vector.SetValue(dofs[2], localActions[2])
	}
}
