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
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

/*
ElementSolution is the displacements and stresses for a given preprocessed element.

Displacements are stored in both local and global coordinates. Stresses are referred
only to the local reference frame.
*/
type ElementSolution struct {
	*preprocess.Element

	GlobalXDispl []PointSolutionValue
	GlobalYDispl []PointSolutionValue
	GlobalZRot   []PointSolutionValue

	LocalXDispl []PointSolutionValue
	LocalYDispl []PointSolutionValue
	LocalZRot   []PointSolutionValue

	AxialStress   []PointSolutionValue
	ShearStress   []PointSolutionValue
	BendingMoment []PointSolutionValue
}

/*
MakeElementSolution creates an empty solution for the given element.
*/
func MakeElementSolution(element *preprocess.Element) *ElementSolution {
	nOfNodes := len(element.Nodes)

	return &ElementSolution{
		Element:      element,
		GlobalXDispl: make([]PointSolutionValue, nOfNodes),
		GlobalYDispl: make([]PointSolutionValue, nOfNodes),
		GlobalZRot:   make([]PointSolutionValue, nOfNodes),
		LocalXDispl:  make([]PointSolutionValue, nOfNodes),
		LocalYDispl:  make([]PointSolutionValue, nOfNodes),
		LocalZRot:    make([]PointSolutionValue, nOfNodes),

		AxialStress:   make([]PointSolutionValue, nOfNodes-1),
		ShearStress:   make([]PointSolutionValue, 2*nOfNodes-2),
		BendingMoment: make([]PointSolutionValue, 2*nOfNodes-2),
	}
}

/*
SolveUsingDisplacements sets the element's global and local displacements
given the structure's system of equations solution vector (the global node
displacements) and computes the stresses in each of the slices of the
preprocessed element.
*/
func (es *ElementSolution) SolveUsingDisplacements(globalDisp *vec.Vector) {
	es.setDisplacements(globalDisp)
	es.computeStresses()
}

/*
setDisplacements sets the global and local displacements given the structure's
system of equations solution vector (the global node displacements).
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
computeStresses use the displacements to compute the stress in each of the
slices of the preprocessed structure.

This method should be called after SetDisplacements, as it depends on the
displacements.
*/
func (es *ElementSolution) computeStresses() {
	var (
		trailNode, leadNode                               *preprocess.Node
		youngMod                                          = es.Element.Material().YoungMod
		iStrong                                           = es.Element.Section().IStrong
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
		axialStressValue := (leadDx - trailDx) * youngMod / length
		es.AxialStress[i-1] = PointSolutionValue{
			trailNode.T,
			axialStressValue,
		}

		/* Shear */
		var (
			shearDispTerm = 12.0 * eil3 * (trailDy - leadDy)
			shearRotTerm  = 6.0 * eil2 * (trailRz + leadRz)
			v             = shearDispTerm + shearRotTerm
		)
		es.ShearStress[j] = PointSolutionValue{trailNode.T, v - trailNode.LocalFy()}
		es.ShearStress[j+1] = PointSolutionValue{leadNode.T, v - leadNode.LocalFy()}

		/* Bending */
		var (
			bendStartDispTerm = 6.0 * eil2 * (leadDy - trailDy)
			bendStartRotTerm  = 2.0 * eil * (leadRz + 2.0*trailRz)
			bendEndDispTerm   = 6.0 * eil2 * (trailDy - leadDy)
			bendEndRotTerm    = 2.0 * eil * (trailRz + 2.0*leadRz)
		)
		es.BendingMoment[j] =
			PointSolutionValue{
				trailNode.T,
				bendStartDispTerm - bendStartRotTerm,
			}
		es.BendingMoment[j+1] =
			PointSolutionValue{
				leadNode.T,
				bendEndDispTerm + bendEndRotTerm,
			}
	}
}
