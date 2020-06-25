package process

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
)

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

	AxialStress   []PointSolutionValue
	ShearStress   []PointSolutionValue
	BendingMoment []PointSolutionValue
}

func (es *ElementSolution) OriginalElementString() string {
	return es.Element.OriginalElementString()
}

/*
MakeElementSolution creates an empty solution for the given element.
*/
func MakeElementSolution(element *preprocess.Element) *ElementSolution {
	nOfNodes := len(element.Nodes)

	fmt.Printf("-> make element %s\n", element.OriginalElementString())

	return &ElementSolution{
		Element:       element,
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
