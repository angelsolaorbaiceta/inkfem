package process

import (
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

/*
ElementSolution is the displacements and stresses for a given preprocessed element.

Displacements are stored in both local and global coordinates. Stresses are referred only to the
local reference frame.
*/
type ElementSolution struct {
	*preprocess.Element

	GlobalXDispl []PointSolutionValue
	GlobalYDispl []PointSolutionValue
	GlobalZRot   []PointSolutionValue

	LocalXDispl []PointSolutionValue
	LocalYDispl []PointSolutionValue
	LocalZRot   []PointSolutionValue

	AxialStress                      []PointSolutionValue
	ShearForce                       []PointSolutionValue
	BendingMoment                    []PointSolutionValue
	BendingMomentTopFiberAxialStress []PointSolutionValue
}

// MakeElementSolution creates an empty solution for the given element.
func MakeElementSolution(element *preprocess.Element) *ElementSolution {
	nOfNodes := len(element.Nodes)

	return &ElementSolution{
		Element: element,

		GlobalXDispl: make([]PointSolutionValue, nOfNodes),
		GlobalYDispl: make([]PointSolutionValue, nOfNodes),
		GlobalZRot:   make([]PointSolutionValue, nOfNodes),

		LocalXDispl: make([]PointSolutionValue, nOfNodes),
		LocalYDispl: make([]PointSolutionValue, nOfNodes),
		LocalZRot:   make([]PointSolutionValue, nOfNodes),

		AxialStress:                      make([]PointSolutionValue, 2*nOfNodes-2),
		ShearForce:                       make([]PointSolutionValue, 2*nOfNodes-2),
		BendingMoment:                    make([]PointSolutionValue, 2*nOfNodes-2),
		BendingMomentTopFiberAxialStress: make([]PointSolutionValue, 2*nOfNodes-2),
	}
}

/*
SolveUsingDisplacements sets the element's global and local displacements given the structure's
system of equations solution vector (the global node displacements) and computes the stresses in
each of the slices of the preprocessed element.
*/
func (es *ElementSolution) SolveUsingDisplacements(globalDisp *vec.Vector) {
	es.setDisplacements(globalDisp)
	es.computeStresses()
}

/*
setDisplacements sets the global and local displacements given the structure's system of equations
solution vector (the global node displacements).
*/
func (es *ElementSolution) setDisplacements(globalDisp *vec.Vector) {
	var (
		nodeDofs               [3]int
		localDisplacementsProj g2d.Projectable
		elementFrame           g2d.RefFrame
	)

	for j, node := range es.Element.Nodes {
		nodeDofs = node.DegreesOfFreedomNum()

		// global displacements
		es.GlobalXDispl[j] = PointSolutionValue{
			node.T,
			globalDisp.Value(nodeDofs[0]),
		}
		es.GlobalYDispl[j] = PointSolutionValue{
			node.T,
			globalDisp.Value(nodeDofs[1]),
		}
		es.GlobalZRot[j] = PointSolutionValue{
			node.T,
			globalDisp.Value(nodeDofs[2]),
		}

		// local displacements
		elementFrame = es.Element.Geometry.RefFrame()
		localDisplacementsProj = elementFrame.ProjectProjections(
			es.GlobalXDispl[j].Value,
			es.GlobalYDispl[j].Value,
		)

		es.LocalXDispl[j] = PointSolutionValue{
			node.T,
			localDisplacementsProj.X,
		}
		es.LocalYDispl[j] = PointSolutionValue{
			node.T,
			localDisplacementsProj.Y,
		}
		es.LocalZRot[j] = PointSolutionValue{
			node.T,
			es.GlobalZRot[j].Value,
		}
	}
}

/*
computeStresses use the displacements to compute the stress in each of the slices of the
preprocessed structure.

This method should be called after SetDisplacements, as it depends on the displacements.
*/
func (es *ElementSolution) computeStresses() {
	var (
		trailNode, leadNode                               *preprocess.Node
		youngMod                                          = es.Element.Material().YoungMod
		iStrong                                           = es.Element.Section().IStrong
		sStrong                                           = es.Element.Section().SStrong
		section                                           = es.Section().Area
		ei                                                = youngMod * iStrong
		trailDx, leadDx, trailDy, leadDy, trailRz, leadRz float64
		length, length2, length3, eil, eil2, eil3         float64
		j                                                 int
	)

	for i := 1; i < len(es.Element.Nodes); i++ {
		j = 2 * (i - 1)
		trailNode, leadNode = es.Element.Nodes[i-1], es.Element.Nodes[i]
		length = es.Element.Geometry.LengthBetween(trailNode.T, leadNode.T)
		length2 = length * length
		length3 = length2 * length
		eil = ei / length
		eil2 = ei / length2
		eil3 = ei / length3
		trailDx = es.LocalXDispl[i-1].Value
		leadDx = es.LocalXDispl[i].Value
		trailDy = es.LocalYDispl[i-1].Value
		leadDy = es.LocalYDispl[i].Value
		trailRz = es.LocalZRot[i-1].Value
		leadRz = es.LocalZRot[i].Value

		/* <-- Axial --> */
		var (
			axial      = (leadDx - trailDx) * youngMod / length
			trailAxial = axial + (trailNode.LocalLeftFx() / section)
			leadAxial  = axial - (leadNode.LocalRightFx() / section)
		)
		es.AxialStress[j] = PointSolutionValue{trailNode.T, trailAxial}
		es.AxialStress[j+1] = PointSolutionValue{leadNode.T, leadAxial}

		/* <-- Shear --> */
		var (
			shearDispTerm = 12.0 * eil3 * (trailDy - leadDy)
			shearRotTerm  = 6.0 * eil2 * (trailRz + leadRz)
			shear         = shearDispTerm + shearRotTerm
			trailShear    = shear - trailNode.LocalLeftFy()
			leadShear     = shear + leadNode.LocalRightFy()
		)
		es.ShearForce[j] = PointSolutionValue{trailNode.T, trailShear}
		es.ShearForce[j+1] = PointSolutionValue{leadNode.T, leadShear}

		/* <-- Bending --> */
		var (
			bendStartDispTerm = 6.0 * eil2 * (leadDy - trailDy)
			bendStartRotTerm  = 2.0 * eil * (leadRz + 2.0*trailRz)
			bendEndDispTerm   = 6.0 * eil2 * (trailDy - leadDy)
			bendEndRotTerm    = 2.0 * eil * (trailRz + 2.0*leadRz)
			trailBending      = bendStartDispTerm - bendStartRotTerm + trailNode.LocalLeftMz()
			leadBending       = bendEndDispTerm + bendEndRotTerm - leadNode.LocalRightMz()
		)
		es.BendingMoment[j] = PointSolutionValue{trailNode.T, trailBending}
		es.BendingMomentTopFiberAxialStress[j] = PointSolutionValue{trailNode.T, trailBending / sStrong}

		es.BendingMoment[j+1] = PointSolutionValue{leadNode.T, leadBending}
		es.BendingMomentTopFiberAxialStress[j+1] = PointSolutionValue{leadNode.T, leadBending / sStrong}
	}
}
