/*
Copyright 2020 Angel Sola Orbaiceta

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package process

import (
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
)

/*
Solve assembles the system of equations for the structure and solves it using
the Preconditioned Conjugate Gradient numerical procedure.

Using the displacements obtained from the solution of the system's solution,
the local stresses are computed.
*/
func Solve(structure *preprocess.Structure, options SolveOptions) *Solution {
	globalDisplacements := computeGlobalDisplacements(structure, options)

	var (
		elementSolution  *ElementSolution
		elementSolutions = make([]*ElementSolution, structure.ElementsCount())
	)

	log.StartComputeStresses()
	for i, element := range structure.Elements {
		elementSolution = MakeElementSolution(element)
		elementSolution.SolveUsingDisplacements(globalDisplacements)
		elementSolutions[i] = elementSolution
	}
	log.EndComputeStresses()

	return &Solution{Metadata: &structure.Metadata, Elements: elementSolutions}
}
