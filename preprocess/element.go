/*
Package preprocess defines the 'preprocessed' or 'sliced' structure model which
is used for the Finite Element Analysis.

This package also provides the means for slicing or preprocessing the structure
as it is defined in the 'structure' package.
*/
package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// Element after slicing original structural element
type Element struct {
	OriginalElement structure.Element
	Nodes           []Node
}

/* ::::::::::::::: Construction ::::::::::::::: */

// MakeElement creates a new element given the original element and the nodes
// of the sliced result.
func MakeElement(originalElement structure.Element, nodes []Node) Element {
	return Element{originalElement, nodes}
}

/* ::::::::::::::: sort.Interface ::::::::::::::: */

// ByGeometryPos implements sort.Interface for []Element based on the position of the original geometry.
type ByGeometryPos []Element

func (a ByGeometryPos) Len() int {
	return len(a)
}

func (a ByGeometryPos) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByGeometryPos) Less(i, j int) bool {
	iStart := a[i].OriginalElement.Geometry.Start
	jStart := a[j].OriginalElement.Geometry.Start
	if pos := iStart.Compare(jStart); pos != 0 {
		return pos < 0
	}

	iEnd := a[i].OriginalElement.Geometry.End
	jEnd := a[j].OriginalElement.Geometry.End
	return iEnd.Compare(jEnd) < 0
}
