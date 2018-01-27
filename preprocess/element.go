/*
Package preprocess defines the 'preprocessed' or 'sliced' structure model which
is used for the Finite Element Analysis.

This package also provides the means for slicing or preprocessing the structure
as it is defined in the 'structure' package.
*/
package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom"
)

// Element after slicing original structural element
type Element struct {
	originalElement *structure.Element
	Nodes           []Node
}

/* ::::::::::::::: Construction ::::::::::::::: */

/*
MakeElement creates a new element given the original element and the nodes
of the sliced result.
*/
func MakeElement(originalElement *structure.Element, nodes []Node) Element {
	return Element{originalElement, nodes}
}

/* ::::::::::::::: Properties ::::::::::::::: */

// ID returns the id of the original structural element.
func (e Element) ID() int {
	return e.originalElement.Id
}

// Geometry returns a pointer to the geometry of the original structural element.
func (e Element) Geometry() *inkgeom.Segment {
	return &e.originalElement.Geometry
}

// StartNodeID returns the id of the start node in the original structural element.
func (e Element) StartNodeID() int {
	return e.originalElement.StartNodeId
}

// EndNodeID returns the id of the end node in the original structural element.
func (e Element) EndNodeID() int {
	return e.originalElement.EndNodeId
}

// StartLink returns the link of the original structural element with the start node.
func (e Element) StartLink() *structure.Constraint {
	return e.originalElement.StartLink
}

// EndLink returns the link of the original structural element with the end node.
func (e Element) EndLink() *structure.Constraint {
	return e.originalElement.EndLink
}

// OriginalElementString returns the string representation of the original structural element.
func (e Element) OriginalElementString() string {
	return e.originalElement.String()
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
	iStart := a[i].Geometry().Start
	jStart := a[j].Geometry().Start
	if pos := iStart.Compare(jStart); pos != 0 {
		return pos < 0
	}

	iEnd := a[i].Geometry().End
	jEnd := a[j].Geometry().End
	return iEnd.Compare(jEnd) < 0
}
