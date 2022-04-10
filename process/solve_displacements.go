package process

import (
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
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
	sysMatrix, sysVector := structure.MakeSystemOfEquations()
	log.EndAssembleSysEqs(sysVector.Length())

	if options.SaveSysMatrixImage {
		go mat.ToImage(sysMatrix, options.OutputPath)
	}

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

	return globalDispSolution.Solution
}

func computePreconditioner(m mat.ReadOnlyMatrix) mat.ReadOnlyMatrix {
	precond := mat.MakeSparse(m.Rows(), m.Cols())
	for i := 0; i < m.Rows(); i++ {
		precond.SetValue(i, i, 1.0/m.Value(i, i))
	}

	return precond
}

func logProgress(ch <-chan lineq.IterativeSolverProgress) {
	for progress := range ch {
		log.SolveSysProgress(progress)
	}
}
