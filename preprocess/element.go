package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkmath/mat"
)

// Element after slicing original structural element.
type Element struct {
	*structure.Element
	Nodes          []*Node
	globalStiffMat []mat.ReadOnlyMatrix
}

// MakeElement creates a new element given the original element and the nodes of the sliced result.
func MakeElement(originalElement *structure.Element, nodes []*Node) *Element {
	matrices := make([]mat.ReadOnlyMatrix, len(nodes)-1)
	return &Element{originalElement, nodes, matrices}
}

// NodesCount returns the number of nodes in the sliced element.
func (e Element) NodesCount() int {
	return len(e.Nodes)
}

/*
ComputeStiffnessMatrices sets the global stiffness matrices for this element in place
(mutates the element).

Each element has a stiffness matrix between two contiguous nodes, so in total that makes n - 1
matrices, where n is the number of nodes.
*/
func (e *Element) ComputeStiffnessMatrices() {
	var trail, lead *Node

	for i := 1; i < len(e.Nodes); i++ {
		trail = e.Nodes[i-1]
		lead = e.Nodes[i]
		e.globalStiffMat[i-1] = e.StiffnessGlobalMat(trail.T, lead.T)
	}
}

/*
GlobalStiffMatrixAt returns the global stiffness matrix at position i, that is, between
nodes i and i + 1.
*/
func (e Element) GlobalStiffMatrixAt(i int) mat.ReadOnlyMatrix {
	return e.globalStiffMat[i]
}

/* <-- sort.Interface --> */

/*
ByGeometryPos implements sort.Interface for []Element based on the position of the
original geometry.
*/
type ByGeometryPos []*Element

func (a ByGeometryPos) Len() int {
	return len(a)
}

func (a ByGeometryPos) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByGeometryPos) Less(i, j int) bool {
	iStart := a[i].StartPoint()
	jStart := a[j].StartPoint()
	if pos := iStart.Compare(jStart); pos != 0 {
		return pos < 0
	}

	iEnd := a[i].EndPoint()
	jEnd := a[j].EndPoint()
	return iEnd.Compare(jEnd) < 0
}
