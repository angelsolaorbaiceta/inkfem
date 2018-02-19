package process

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom"
)

// PointSolutionValue is a tuple of T and Value.
type PointSolutionValue struct {
	T     inkgeom.TParam
	Value float64
}

func (psv PointSolutionValue) String() string {
	return fmt.Sprintf("T = %f : %f", psv.T.Value(), psv.Value)
}

/*
ElementSolution is the displacements and stresses for a given preprocessed element.

Displacements are stored in both local and global coordinates. Stresses are referred
only to the local reference frame.
*/
type ElementSolution struct {
	Element *preprocess.Element

	GlobalXDispl []PointSolutionValue
	GlobalYDispl []PointSolutionValue
	GlobalZRot   []PointSolutionValue

	LocalXDispl []PointSolutionValue
	LocalYDispl []PointSolutionValue
	LocalZRot   []PointSolutionValue

	// Points map[inkgeom.TParam]inkgeom.Projectable

	AxialStress   []PointSolutionValue
	ShearStress   []PointSolutionValue
	BendingMoment []PointSolutionValue
}

/*
MakeElementSolution creates an empty solution for the given element.
*/
func MakeElementSolution(element preprocess.Element) ElementSolution {
	nOfNodes := len(element.Nodes)

	return ElementSolution{
		Element:      &element,
		GlobalXDispl: make([]PointSolutionValue, nOfNodes),
		GlobalYDispl: make([]PointSolutionValue, nOfNodes),
		GlobalZRot:   make([]PointSolutionValue, nOfNodes),
		LocalXDispl:  make([]PointSolutionValue, nOfNodes),
		LocalYDispl:  make([]PointSolutionValue, nOfNodes),
		LocalZRot:    make([]PointSolutionValue, nOfNodes),
		// Points:        make([]PointSolutionValue, nOfNodes),
		AxialStress:   make([]PointSolutionValue, nOfNodes),
		ShearStress:   make([]PointSolutionValue, nOfNodes),
		BendingMoment: make([]PointSolutionValue, 2*nOfNodes-1),
	}
}

/*
Solution is the group of all element solutions with the structure metadata.
*/
type Solution struct {
	Metadata *structure.StrMetadata
	Elements []ElementSolution
}
