/*
Package process defines the Finite Element Method computation.
It starts from the sliced structure, assembles the global system of equations, solves it
and creates a solution.
*/
package process

import (
	"sync"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkmath"
	"github.com/angelsolaorbaiceta/inkmath/lineq"
	"github.com/angelsolaorbaiceta/inkmath/mat"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

/*
Solve assembles the system of equations for the structure and solves it using a numerical
procedure. Using the displacements obtained from the solution of the system, the local
stresses are computed.
*/
func Solve(s *preprocess.Structure, options SolveOptions) *Solution {
	sysMatrix, sysVector := makeSystemOfEqs(s)

	if options.SaveSysMatrixImage {
		go mat.ToImage(sysMatrix, options.OutputPath)
	}

	solver := lineq.PreconditionedConjugateGradientSolver{MaxError: options.MaxDisplacementsError, MaxIter: sysVector.Length()}
	if options.SafeChecks && !solver.CanSolve(sysMatrix, sysVector) {
		panic("Solver cannot solve system!")
	}

	var (
		globalDisplacements     = solver.Solve(sysMatrix, sysVector)
		globalDisplacementsProj inkgeom.Projectable
		elementSolution         ElementSolution
		elementSolutions        = make([]ElementSolution, len(s.Elements))
		nodeDofs                [3]int
		wg                      sync.WaitGroup
	)

	for i, element := range s.Elements {
		elementSolution = MakeElementSolution(element)

		for j, node := range element.Nodes {
			nodeDofs = node.DegreesOfFreedomNum()

			// global displacements
			elementSolution.GlobalXDispl[j] = PointSolutionValue{node.T, globalDisplacements.Solution.Value(nodeDofs[0])}
			elementSolution.GlobalYDispl[j] = PointSolutionValue{node.T, globalDisplacements.Solution.Value(nodeDofs[1])}
			elementSolution.GlobalZRot[j] = PointSolutionValue{node.T, globalDisplacements.Solution.Value(nodeDofs[2])}

			// local displacements
			globalDisplacementsProj = element.Geometry().RefFrame().ProjectProjections(elementSolution.GlobalXDispl[j].Value, elementSolution.GlobalYDispl[j].Value)
			elementSolution.LocalXDispl[j] = PointSolutionValue{node.T, globalDisplacementsProj.X}
			elementSolution.LocalYDispl[j] = PointSolutionValue{node.T, globalDisplacementsProj.X}
			elementSolution.LocalZRot[j] = PointSolutionValue{node.T, elementSolution.GlobalZRot[j].Value}
		}

		wg.Add(1)
		elementSolutions[i] = elementSolution
		go computeStresses(&elementSolutions[i], &wg)
	}

	wg.Wait()
	return &Solution{Metadata: &s.Metadata, Elements: elementSolutions}
}

/* ::::::::::::::: Solve Displacements ::::::::::::::: */
func makeSystemOfEqs(s *preprocess.Structure) (mat.Matrixable, *vec.Vector) {
	c := make(chan preprocess.Element)

	for _, element := range s.Elements {
		go element.ComputeStiffnessMatrices(c)
	}

	sysMatrix, sysVector := mat.MakeSparse(s.DofsCount, s.DofsCount), vec.Make(s.DofsCount)
	for i := 0; i < len(s.Elements); i++ {
		element := <-c
		addTermsToStiffnessMatrix(sysMatrix, &element)
		addTermsToLoadVector(sysVector, &element)
	}
	addDispConstraints(sysMatrix, sysVector, s.Nodes)

	return sysMatrix, sysVector
}

func addTermsToStiffnessMatrix(m mat.Matrixable, e *preprocess.Element) {
	var (
		stiffMat                    mat.Matrixable
		trailNodeDofs, leadNodeDofs [3]int
		dofs                        [6]int
		stiffVal                    float64
	)

	for i := 1; i < len(e.Nodes); i++ {
		stiffMat = e.GlobalStiffMatrixAt(i - 1)
		trailNodeDofs = e.Nodes[i-1].DegreesOfFreedomNum()
		leadNodeDofs = e.Nodes[i].DegreesOfFreedomNum()
		dofs = [6]int{
			trailNodeDofs[0], trailNodeDofs[1], trailNodeDofs[2],
			leadNodeDofs[0], leadNodeDofs[1], leadNodeDofs[2],
		}

		for i := 0; i < stiffMat.Rows(); i++ {
			for j := 0; j < stiffMat.Cols(); j++ {
				if stiffVal = stiffMat.Value(i, j); !inkmath.IsCloseToZero(stiffVal) {
					m.AddToValue(dofs[i], dofs[j], stiffVal)
				}
			}
		}
	}
}

func addDispConstraints(m mat.Matrixable, v *vec.Vector, nodes map[int]structure.Node) {
	var (
		constraint *structure.Constraint
		dofs       [3]int
	)

	addConstraintAtDof := func(dof int) {
		m.SetZeroCol(dof)
		m.SetIdentityRow(dof)
		v.SetZero(dof)
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

func addTermsToLoadVector(v *vec.Vector, e *preprocess.Element) {
	var (
		localActions [3]float64
		globalForces inkgeom.Projectable
		dofs         [3]int
		refFrame     = e.Geometry().RefFrame()
	)

	for _, node := range e.Nodes {
		localActions = node.LocalActions()
		globalForces = refFrame.ProjectionsToGlobal(localActions[0], localActions[1])
		dofs = node.DegreesOfFreedomNum()

		v.SetValue(dofs[0], globalForces.X)
		v.SetValue(dofs[1], globalForces.Y)
		v.SetValue(dofs[2], localActions[2])
	}
}

/* ::::::::::::::: Compute Stresses ::::::::::::::: */
func computeStresses(sol *ElementSolution, wg *sync.WaitGroup) {
	defer wg.Done()

	var (
		trailNode, leadNode                    preprocess.Node
		youngMod                               = sol.Element.Material().YoungMod
		iStrong                                = sol.Element.Section().IStrong
		nIndex, vIndex, mIndex                 = 0, 0, 0
		incX, trailDy, leadDy, trailRz, leadRz float64
		length                                 float64
	)

	for i := 1; i < len(sol.Element.Nodes); i++ {
		trailNode, leadNode = sol.Element.Nodes[i-1], sol.Element.Nodes[i]
		length = sol.Element.Geometry().LengthBetween(trailNode.T, leadNode.T)
		incX = sol.LocalXDispl[i].Value - sol.LocalXDispl[i-1].Value
		trailDy = sol.LocalYDispl[i-1].Value
		leadDy = sol.LocalYDispl[i].Value
		trailRz = sol.LocalZRot[i-1].Value
		leadRz = sol.LocalZRot[i].Value

		/* Axial */
		n := incX * youngMod / length
		sol.AxialStress[nIndex] = PointSolutionValue{trailNode.T, n - trailNode.LocalFx()}
		sol.AxialStress[nIndex+1] = PointSolutionValue{leadNode.T, n + leadNode.LocalFx()}
		nIndex += 2

		/* Shear */
		v := (6.0 * youngMod * iStrong / (length * length * length)) * ((2.0 * (trailDy - leadDy)) + (length * (leadRz - trailRz)))
		sol.ShearStress[vIndex] = PointSolutionValue{trailNode.T, v - trailNode.LocalFy()}
		sol.ShearStress[vIndex+1] = PointSolutionValue{leadNode.T, v + leadNode.LocalFy()}
		vIndex += 2

		/* Bending */
		eil2 := youngMod * iStrong / (length * length)
		sol.BendingMoment[mIndex] =
			PointSolutionValue{
				trailNode.T,
				eil2*(-6.0*trailDy+2.0*length*trailRz-6.0*leadDy+4.0*length*leadRz) + trailNode.LocalMz(),
			}
		sol.BendingMoment[mIndex+1] =
			PointSolutionValue{
				inkgeom.AverageT(trailNode.T, leadNode.T),
				(youngMod * iStrong / length) * (leadRz - trailRz),
			}
		sol.BendingMoment[mIndex+2] =
			PointSolutionValue{
				leadNode.T,
				eil2*(-6.0*trailDy+4.0*length*trailRz+6.0*leadDy-2.0*length*leadRz) + leadNode.LocalMz(),
			}
		mIndex += 3
	}
}
