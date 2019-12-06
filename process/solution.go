package process

import (
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

/*
ElementSolution is the displacements and stresses for a given preprocessed element.

Displacements are stored in both local and global coordinates. Stresses are referred
only to the local reference frame.
*/
type ElementSolution struct {
	Element *preprocess.Element `json:"-"`

	GlobalXDispl []PointSolutionValue `json:"g_x_disp"`
	GlobalYDispl []PointSolutionValue `json:"g_y_disp"`
	GlobalZRot   []PointSolutionValue `json:"g_z_rot"`

	LocalXDispl []PointSolutionValue `json:"l_x_disp"`
	LocalYDispl []PointSolutionValue `json:"l_y_disp"`
	LocalZRot   []PointSolutionValue `json:"l_z_rot"`

	AxialStress   []PointSolutionValue `json:"axial"`
	ShearStress   []PointSolutionValue `json:"shear"`
	BendingMoment []PointSolutionValue `json:"bending"`
}

/*
MakeElementSolution creates an empty solution for the given element.
*/
func MakeElementSolution(element preprocess.Element) ElementSolution {
	nOfNodes := len(element.Nodes)

	return ElementSolution{
		Element:       &element,
		GlobalXDispl:  make([]PointSolutionValue, nOfNodes),
		GlobalYDispl:  make([]PointSolutionValue, nOfNodes),
		GlobalZRot:    make([]PointSolutionValue, nOfNodes),
		LocalXDispl:   make([]PointSolutionValue, nOfNodes),
		LocalYDispl:   make([]PointSolutionValue, nOfNodes),
		LocalZRot:     make([]PointSolutionValue, nOfNodes),
		AxialStress:   make([]PointSolutionValue, 2*nOfNodes-2),
		ShearStress:   make([]PointSolutionValue, 2*nOfNodes-2),
		BendingMoment: make([]PointSolutionValue, 3*nOfNodes-3),
	}
}

/*
Solution is the group of all element solutions with the structure metadata.
*/
type Solution struct {
	Metadata *structure.StrMetadata `json:"metadata"`
	Elements []ElementSolution      `json:"elements"`
}
