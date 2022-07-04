package process

import (
	"github.com/angelsolaorbaiceta/inkfem/build"
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Solve assembles the system of equations for the structure and solves it using the
// Preconditioned Conjugate Gradient numerical procedure and sets the bars local stresses,
// forces and moments.
func Solve(str *preprocess.Structure, options SolveOptions) *Solution {
	var (
		globalDispl      = computeGlobalDisplacements(str, options)
		elementSolutions = make([]*ElementSolution, str.ElementsCount())
		metadata         = structure.StrMetadata{
			MajorVersion: build.Info.MajorVersion,
			MinorVersion: build.Info.MinorVersion,
		}
	)

	log.StartComputeStresses()
	for i, element := range str.Elements() {
		elementSolutions[i] = MakeElementSolution(element, globalDispl)
	}
	log.EndComputeStresses()

	return MakeSolution(metadata, str.NodesById, elementSolutions)
}
