package process

import (
	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

// ElementSolution is the displacements, stresses, forces and moments for a given
// preprocessed element.
//
// The displacements are stored in both local and global coordinates. Stresses, forces
// and moments are referred only to the local reference frame of the bar.
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

// MakeElementSolution creates a solution element with all solution values for the
// preprocessed element.
//
// It sets the element's global and local displacements given the structure's
// system of equations solution vector (the global node displacements) and computes
// the axial stress, shear force and bending moment in each of the slices of the
// preprocessed element.
//
// To compare whether two stresses or forces are the same in both sides of a node,
// it is necessary to compare the values of the solution values using the maximum
// displacement error used for the displacements calculation. In case a node has the
// same values in both sides (at the same T position), only one is kept.
func MakeElementSolution(
	element *preprocess.Element,
	globalDisp *GlobalDisplacementsVector,
) *ElementSolution {
	var (
		nOfNodes          = element.NodesCount()
		nOfSolutionValues = 2*nOfNodes - 2
	)

	solution := &ElementSolution{
		Element: element,

		GlobalXDispl: make([]PointSolutionValue, nOfNodes),
		GlobalYDispl: make([]PointSolutionValue, nOfNodes),
		GlobalZRot:   make([]PointSolutionValue, nOfNodes),

		LocalXDispl: make([]PointSolutionValue, nOfNodes),
		LocalYDispl: make([]PointSolutionValue, nOfNodes),
		LocalZRot:   make([]PointSolutionValue, nOfNodes),

		AxialStress:                      make([]PointSolutionValue, 0, nOfSolutionValues),
		ShearForce:                       make([]PointSolutionValue, 0, nOfSolutionValues),
		BendingMoment:                    make([]PointSolutionValue, 0, nOfSolutionValues),
		BendingMomentTopFiberAxialStress: make([]PointSolutionValue, 0, nOfSolutionValues),
	}

	solution.setDisplacements(globalDisp.Vector)
	solution.computeStresses(globalDisp.MaxError)

	return solution
}

// RefFrame returns the element's reference frame.
func (es *ElementSolution) RefFrame() *g2d.RefFrame {
	return es.Element.RefFrame()
}

// setDisplacements sets the global and local displacements given the structure's system
// of equations solution vector (the global node displacements).
func (es *ElementSolution) setDisplacements(globalDisp vec.ReadOnlyVector) {
	var (
		nodeDofs               [3]int
		localDisplacementsProj *g2d.Vector
		elementFrame           *g2d.RefFrame
	)

	for j, node := range es.Element.Nodes() {
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
		elementFrame = es.Element.RefFrame()
		localDisplacementsProj = elementFrame.ProjectProjections(
			es.GlobalXDispl[j].Value,
			es.GlobalYDispl[j].Value,
		)

		es.LocalXDispl[j] = PointSolutionValue{
			node.T,
			localDisplacementsProj.X(),
		}
		es.LocalYDispl[j] = PointSolutionValue{
			node.T,
			localDisplacementsProj.Y(),
		}
		es.LocalZRot[j] = PointSolutionValue{
			node.T,
			es.GlobalZRot[j].Value,
		}
	}
}

// computeStresses use the displacements to compute the stress in each of the slices
// of the preprocessed structure.
//
// This method should be called after setDisplacements, as it requires the displacements
// to compute the corresponding stresses and forces.
//
// It uses the maximum displacement error used for the displacements calculation to
// compare the values of the solution values and avoid having the same force/stress value
// in both sides of the node.
func (es *ElementSolution) computeStresses(maxDispError float64) {
	var (
		trailNode, leadNode *preprocess.Node
		youngMod            = es.Element.Material().YoungMod
		iStrong             = es.Element.Section().IStrong
		sStrong             = es.Element.Section().SStrong
		section             = es.Section().Area
		ei                  = youngMod * iStrong
		nodesCount          = es.Element.NodesCount()

		trailDx, leadDx, trailDy, leadDy, trailRz, leadRz float64
		length, length2, length3, eil, eil2, eil3         float64
	)

	for i := 1; i < nodesCount; i++ {
		trailNode, leadNode = es.Element.NodeAt(i-1), es.Element.NodeAt(i)
		length = es.Element.LengthBetween(trailNode.T, leadNode.T)
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
		es.AxialStress = appendIfNotSameAsLast(
			es.AxialStress,
			PointSolutionValue{trailNode.T, trailAxial},
			maxDispError,
		)
		es.AxialStress = append(es.AxialStress, PointSolutionValue{leadNode.T, leadAxial})

		/* <-- Shear --> */
		var (
			shearDispTerm = 12.0 * eil3 * (trailDy - leadDy)
			shearRotTerm  = 6.0 * eil2 * (trailRz + leadRz)
			shear         = shearDispTerm + shearRotTerm
			trailShear    = shear - trailNode.LocalLeftFy()
			leadShear     = shear + leadNode.LocalRightFy()
		)
		es.ShearForce = appendIfNotSameAsLast(
			es.ShearForce,
			PointSolutionValue{trailNode.T, trailShear},
			maxDispError,
		)
		es.ShearForce = append(es.ShearForce, PointSolutionValue{leadNode.T, leadShear})

		/* <-- Bending --> */
		var (
			bendStartDispTerm = 6.0 * eil2 * (leadDy - trailDy)
			bendStartRotTerm  = 2.0 * eil * (leadRz + 2.0*trailRz)
			bendEndDispTerm   = 6.0 * eil2 * (trailDy - leadDy)
			bendEndRotTerm    = 2.0 * eil * (trailRz + 2.0*leadRz)
			trailBending      = bendStartDispTerm - bendStartRotTerm + trailNode.LocalLeftMz()
			leadBending       = bendEndDispTerm + bendEndRotTerm - leadNode.LocalRightMz()
		)
		es.BendingMoment = appendIfNotSameAsLast(
			es.BendingMoment,
			PointSolutionValue{trailNode.T, trailBending},
			maxDispError,
		)
		es.BendingMomentTopFiberAxialStress = appendIfNotSameAsLast(
			es.BendingMomentTopFiberAxialStress,
			PointSolutionValue{trailNode.T, trailBending / sStrong},
			maxDispError,
		)
		es.BendingMoment = append(es.BendingMoment, PointSolutionValue{leadNode.T, leadBending})
		es.BendingMomentTopFiberAxialStress = append(
			es.BendingMomentTopFiberAxialStress,
			PointSolutionValue{leadNode.T, leadBending / sStrong},
		)
	}
}

// GlobalStartTorsor returns the forces and moment torsor {fx, fy, mz} at the start node
// in global coordinates.
//
// Sign convention:
// - a tensile stress (positive) yields a negative force value
// - a positive shear force yields a positive force value
// - a positive bending moment yields a negative moment value
func (es *ElementSolution) GlobalStartTorsor() *math.Torsor {
	return math.MakeTorsor(
		-es.AxialStress[0].Value*es.Section().Area,
		es.ShearForce[0].Value,
		-es.BendingMoment[0].Value,
	).ProjectedToGlobal(es.RefFrame())
}

// GlobalEndTorsor returns the forces and moment torsor {fx, fy, mz} at the end node
// in global coordinates.
//
// Sign convention:
// - a tensile stress (positive) yields a positive force value
// - a positive shear force yields a negative force value
// - a positive bending moment yields a positive moment value
func (es *ElementSolution) GlobalEndTorsor() *math.Torsor {
	var (
		axialIndex   = len(es.AxialStress) - 1
		shearIndex   = len(es.ShearForce) - 1
		bendingIndex = len(es.BendingMoment) - 1
	)

	return math.MakeTorsor(
		es.AxialStress[axialIndex].Value*es.Section().Area,
		-es.ShearForce[shearIndex].Value,
		es.BendingMoment[bendingIndex].Value,
	).ProjectedToGlobal(es.RefFrame())
}
