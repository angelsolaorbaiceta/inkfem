package process

import (
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkmath/lineq"
	"github.com/angelsolaorbaiceta/inkmath/mat"
	"github.com/angelsolaorbaiceta/inkmath/nums"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

/*
ComputeGlobalDisplacements computes the structure's global displacements given the
preprocessed structure.

The process involves generating the structure's system of equations and solving it using the
Preconditioned Conjugate Gradiend numerical procedure.
*/
func computeGlobalDisplacements(structure *preprocess.Structure, options SolveOptions) *vec.Vector {
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

/*
MakeSystemOfEquations generates the system of equations matrix and vector from the
preprocessed structure.

It computes each of the sliced element's stiffness matrices and assembles them into one
global matrix. It also assembles the global loads vector from the sliced element nodes.
*/
func makeSystemOfEquations(structure *preprocess.Structure) (mat.ReadOnlyMatrix, *vec.Vector) {
	var (
		sysMatrix = mat.MakeSparse(structure.DofsCount, structure.DofsCount)
		sysVector = vec.Make(structure.DofsCount)
	)

	for _, element := range structure.Elements {
		element.ComputeStiffnessMatrices()
		addTermsToStiffnessMatrix(sysMatrix, element)
		addTermsToLoadVector(sysVector, element)
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

		for row := 0; row < stiffMat.Rows(); row++ {
			for col := 0; col < stiffMat.Cols(); col++ {
				if stiffVal = stiffMat.Value(row, col); !nums.IsCloseToZero(stiffVal) {
					matrix.AddToValue(dofs[row], dofs[col], stiffVal)
				}
			}
		}
	}
}

func addDispConstraints(
	matrix mat.MutableMatrix,
	vector *vec.Vector,
	nodes *map[contracts.StrID]*structure.Node,
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

func addTermsToLoadVector(sysVector *vec.Vector, element *preprocess.Element) {
	var (
		globalTorsor *math.Torsor
		dofs         [3]int
		refFrame     = element.RefFrame()
	)

	for _, node := range element.Nodes {
		globalTorsor = node.NetLocalLoadTorsor().ProjectedToGlobal(refFrame)
		dofs = node.DegreesOfFreedomNum()

		sysVector.SetValue(dofs[0], globalTorsor.Fx())
		sysVector.SetValue(dofs[1], globalTorsor.Fy())
		sysVector.SetValue(dofs[2], globalTorsor.Mz())
	}
}
