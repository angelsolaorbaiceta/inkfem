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

/* Construction */

// MakeElement creates a new element given the original element and the nodes
// of the sliced result.
func MakeElement(originalElement structure.Element, nodes []Node) Element {
	return Element{originalElement, nodes}
}

/* Sorting */
func (e Element) Id() int {
	return e.OriginalElement.Id
}
