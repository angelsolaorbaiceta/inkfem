package process

import (
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkmath/lineq"
	"github.com/angelsolaorbaiceta/inkmath/mat"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

// A GlobalDisplacementsVector is the solution of the structure's system of equations, that yield
// the structure's global displacements.
// It includes the error upper bound of the displacements calculation.
type GlobalDisplacementsVector struct {
	MaxError float64
	Vector   vec.ReadOnlyVector
}

// computeGlobalDisplacements computes the structure's global displacements given the
// preprocessed structure.
//
// The process involves generating the structure's system of equations and solving it using the
// Preconditioned Conjugate Gradiend numerical procedure.
func computeGlobalDisplacements(
	structure *preprocess.Structure,
	options SolveOptions,
) *GlobalDisplacementsVector {
	log.StartAssembleSysEqs()
	sysMatrix, sysVector := structure.MakeSystemOfEquations()
	log.EndAssembleSysEqs(sysVector.Length())

	log.StartSolveSysEqs()

	var (
		progressChan = make(chan lineq.IterativeSolverProgress)
		solutionChan = make(chan *lineq.Solution)
		solver       = lineq.PreconditionedConjugateGradientSolver{
			MaxError:       options.MaxDisplacementsError,
			MaxIter:        sysVector.Length(),
			Preconditioner: computePreconditioner(sysMatrix),
			ProgressChan:   progressChan,
		}
	)

	if options.SafeChecks && !solver.CanSolve(sysMatrix, sysVector) {
		panic("Solver can't solve system!")
	}

	go func() {
		solutionChan <- solver.Solve(sysMatrix, sysVector)
		close(solutionChan)
	}()

	logProgress(progressChan)
	globalDispSolution := <-solutionChan

	log.EndSolveSysEqs(globalDispSolution.IterCount, globalDispSolution.MinError)

	return &GlobalDisplacementsVector{
		Vector:   globalDispSolution.Solution,
		MaxError: options.MaxDisplacementsError,
	}
}

// computePreconditioner computes the preconditioner of the system matrix for the conjugate
// gradient method to converge faster.
func computePreconditioner(sysMat mat.ReadOnlyMatrix) mat.ReadOnlyMatrix {
	precond := mat.MakeSparse(sysMat.Rows(), sysMat.Cols())
	for i := 0; i < sysMat.Rows(); i++ {
		precond.SetValue(i, i, 1.0/sysMat.Value(i, i))
	}

	return precond
}

func logProgress(ch <-chan lineq.IterativeSolverProgress) {
	for progress := range ch {
		log.SolveSysProgress(progress)
	}
}
