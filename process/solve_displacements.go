package process

import (
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkmath/lineq"
	"github.com/angelsolaorbaiceta/inkmath/mat"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

// ComputeGlobalDisplacements computes the structure's global displacements given the
// preprocessed structure.
//
// The process involves generating the structure's system of equations and solving it using the
// Preconditioned Conjugate Gradiend numerical procedure.
func computeGlobalDisplacements(
	structure *preprocess.Structure,
	options SolveOptions,
) vec.ReadOnlyVector {
	log.StartAssembleSysEqs()
	sysMatrix, sysVector := makeSystemOfEquations(structure)
	log.EndAssembleSysEqs(sysVector.Length())

	if options.SaveSysMatrixImage {
		go mat.ToImage(sysMatrix, options.OutputPath)
	}

	log.StartSolveSysEqs()
	solver := lineq.PreconditionedConjugateGradientSolver{
		MaxError:       options.MaxDisplacementsError,
		MaxIter:        sysVector.Length(),
		Preconditioner: computePreconditioner(sysMatrix),
	}
	if options.SafeChecks && !solver.CanSolve(sysMatrix, sysVector) {
		panic("Solver cannot solve system!")
	}

	globalDispSolution := solver.Solve(sysMatrix, sysVector)
	log.EndSolveSysEqs(globalDispSolution.IterCount, globalDispSolution.MinError)

	return globalDispSolution.Solution
}

func computePreconditioner(m mat.ReadOnlyMatrix) mat.ReadOnlyMatrix {
	precond := mat.MakeSparse(m.Rows(), m.Cols())
	for i := 0; i < m.Rows(); i++ {
		precond.SetValue(i, i, 1.0/m.Value(i, i))
	}

	return precond
}

// MakeSystemOfEquations generates the system of equations matrix and vector from the
// preprocessed structure.
//
// It computes each of the sliced element's stiffness matrices and assembles them into one
// global matrix. It also assembles the global loads vector from the sliced element nodes.
func makeSystemOfEquations(str *preprocess.Structure) (mat.ReadOnlyMatrix, vec.ReadOnlyVector) {
	var (
		sysMatrix = mat.MakeSparse(str.DofsCount(), str.DofsCount())
		sysVector = vec.Make(str.DofsCount())
	)

	for _, element := range str.Elements {
		element.SetEquationTerms(sysMatrix, sysVector)
	}

	addDispConstraints(sysMatrix, sysVector, str.GetAllNodes())

	return sysMatrix, sysVector
}

// Sets the node's external constraints in the system of equations matrix and vector.
//
// A constrained degree of freedom is enforced by setting the corresponding matrix row as the
// identity, and the associated free value as zero. This yields a trivial equation of the form
// x = 0, where x is the constrained degree of freedom.
func addDispConstraints(
	matrix mat.MutableMatrix,
	vector vec.MutableVector,
	nodes []*structure.Node,
) {
	var (
		constraint *structure.Constraint
		dofs       [3]int
	)

	addConstraintAtDof := func(dof int) {
		matrix.SetZeroCol(dof)
		matrix.SetIdentityRow(dof)
		vector.SetZero(dof)
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
