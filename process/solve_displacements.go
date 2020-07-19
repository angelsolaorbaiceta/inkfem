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
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/log"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
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
	log.StartAssembleSysEqs()
	sysMatrix, sysVector := makeSystemOfEquations(structure)
	log.EndAssembleSysEqs(sysVector.Length())

	if options.SaveSysMatrixImage {
		go mat.ToImage(sysMatrix, options.OutputPath)
	}

	log.StartSolveSysEqs()
	solver := lineq.PreconditionedConjugateGradientSolver{
		MaxError: options.MaxDisplacementsError,
		MaxIter:  sysVector.Length(),
	}
	if options.SafeChecks && !solver.CanSolve(sysMatrix, sysVector) {
		panic("Solver cannot solve system!")
	}

	globalDispSolution := solver.Solve(sysMatrix, sysVector)
	log.EndSolveSysEqs(globalDispSolution.IterCount, globalDispSolution.MinError)

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
	nodes *map[contracts.StrID]*structure.Node,
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
		globalForces g2d.Projectable
		dofs         [3]int
		refFrame     = element.Geometry.RefFrame()
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
