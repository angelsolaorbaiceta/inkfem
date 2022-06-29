package process

import (
	"github.com/angelsolaorbaiceta/inkfem/build"
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Solve assembles the system of equations for the structure and solves it using the Preconditioned
// Conjugate Gradient numerical procedure.
//
// The local element stresses are computed using the displacements obtained in the first step.
func Solve(str *preprocess.Structure, options SolveOptions) *Solution {
	var (
		globalDisplacements = computeGlobalDisplacements(str, options)
		elementSolutions    = make([]*ElementSolution, str.ElementsCount())
		metadata            = structure.StrMetadata{
			MajorVersion: build.Info.MajorVersion,
			MinorVersion: build.Info.MinorVersion,
		}
	)

	log.StartComputeStresses()
	for i, element := range str.Elements() {
		elementSolutions[i] = MakeElementSolution(element, globalDisplacements, options.MaxDisplacementsError)
	}
	log.EndComputeStresses()

	return MakeSolution(metadata, str.NodesById, elementSolutions)
}
