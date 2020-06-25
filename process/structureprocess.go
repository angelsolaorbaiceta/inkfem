/*
Package process defines the Finite Element Method computation.
It starts from the sliced structure, assembles the global system of equations,
solves it and creates a solution.
*/
package process

import (
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkgeom"
)

/*
Solve assembles the system of equations for the structure and solves it using
the Preconditioned Conjugate Gradient numerical procedure.

Using the displacements obtained from the solution of the system's solution,
the local stresses are computed.
*/
func Solve(s *preprocess.Structure, options SolveOptions) *Solution {
	globalDisplacements := computeGlobalDisplacements(s, options)

	var (
		localDisplacementsProj inkgeom.Projectable
		elementSolution        *ElementSolution
		elementSolutions       = make([]*ElementSolution, len(s.Elements))
		nodeDofs               [3]int
	)

	for i, element := range s.Elements {
		elementSolution = MakeElementSolution(&element)

		for j, node := range element.Nodes {
			nodeDofs = node.DegreesOfFreedomNum()

			// global displacements
			elementSolution.GlobalXDispl[j] = PointSolutionValue{
				node.T,
				globalDisplacements.Value(nodeDofs[0]),
			}
			elementSolution.GlobalYDispl[j] = PointSolutionValue{
				node.T,
				globalDisplacements.Value(nodeDofs[1]),
			}
			elementSolution.GlobalZRot[j] = PointSolutionValue{
				node.T,
				globalDisplacements.Value(nodeDofs[2]),
			}

			// local displacements
			localDisplacementsProj = element.Geometry().RefFrame().ProjectProjections(
				elementSolution.GlobalXDispl[j].Value,
				elementSolution.GlobalYDispl[j].Value,
			)
			elementSolution.LocalXDispl[j] = PointSolutionValue{
				node.T,
				localDisplacementsProj.X,
			}
			elementSolution.LocalYDispl[j] = PointSolutionValue{
				node.T,
				localDisplacementsProj.X,
			}
			elementSolution.LocalZRot[j] = PointSolutionValue{
				node.T,
				elementSolution.GlobalZRot[j].Value,
			}
		}

		elementSolutions[i] = elementSolution
		computeStresses(elementSolutions[i])
	}

	return &Solution{Metadata: &s.Metadata, Elements: elementSolutions}
}
