/*
Preprocess package defines the 'preprocessed' or 'sliced' structure model which
is used for the Finite Element Analysis.

This package also provides the means for slicing or preprocessing the structure
as it is defined in the 'structure' package.
*/
package preprocess

import (
    "github.com/angelsolaorbaiceta/inkfem/structure"
)

type Element struct {
    originalElement structure.Element
    Nodes []Node
}

/* Construction */
func MakeElement(originalElement structure.Element, nodes []Node) Element {
    return Element{originalElement, nodes}
}
