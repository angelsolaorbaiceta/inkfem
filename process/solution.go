package process

import (
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom"
)

/*
ElementSolution is the displacements and stresses for a given preprocessed element.

Displacements are stored in both local and global coordinates. Stresses are referred
only to the local reference frame.
*/
type ElementSolution struct {
	Element       *preprocess.Element
	GlobalDispl   map[inkgeom.TParam][3]float64
	LocalDispl    map[inkgeom.TParam][3]float64
	Points        map[inkgeom.TParam]inkgeom.Projectable
	AxialStress   map[inkgeom.TParam]float64
	ShearStress   map[inkgeom.TParam]float64
	BendingMoment map[inkgeom.TParam]float64
}

/*
MakeElementSolution creates an empty solution for the given element.
*/
func MakeElementSolution(element preprocess.Element) ElementSolution {
	nOfNodes := len(element.Nodes)

	return ElementSolution{
		Element:       &element,
		GlobalDispl:   make(map[inkgeom.TParam][3]float64, nOfNodes),
		LocalDispl:    make(map[inkgeom.TParam][3]float64, nOfNodes),
		Points:        make(map[inkgeom.TParam]inkgeom.Projectable, nOfNodes),
		AxialStress:   make(map[inkgeom.TParam]float64, nOfNodes),
		ShearStress:   make(map[inkgeom.TParam]float64, nOfNodes),
		BendingMoment: make(map[inkgeom.TParam]float64, 2*nOfNodes-1),
	}
}

/*
Solution is the group of all element solutions with the structure metadata.
*/
type Solution struct {
	Metadata *structure.StrMetadata
	Elements []ElementSolution
}
