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

package preprocess

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkmath/mat"
)

/*
Element after slicing original structural element
*/
type Element struct {
	*structure.Element
	Nodes          []*Node
	globalStiffMat []mat.ReadOnlyMatrix
}

/* <-- Construction --> */

/*
MakeElement creates a new element given the original element and the nodes
of the sliced result.
*/
func MakeElement(originalElement *structure.Element, nodes []*Node) *Element {
	matrices := make([]mat.ReadOnlyMatrix, len(nodes)-1)
	return &Element{originalElement, nodes, matrices}
}

/* <-- Properties --> */

/*
NodesCount returns the number of nodes in the sliced element.
*/
func (e Element) NodesCount() int {
	return len(e.Nodes)
}

/*
OriginalElementString returns the string representation of the original
structural element.
*/
func (e Element) OriginalElementString() string {
	return e.String()
}

/* <-- Methods --> */

/*
ComputeStiffnessMatrices sets the global stiffness matrices for this element
in place (mutates the element).

Each element has a stiffness matrix between two contiguous nodes, so
in total that makes n - 1 matrices, where n is the number of nodes.

TODO: remove the channel thingy from here?
*/
func (e *Element) ComputeStiffnessMatrices(c chan<- Element) {
	var trail, lead *Node

	for i := 1; i < len(e.Nodes); i++ {
		trail = e.Nodes[i-1]
		lead = e.Nodes[i]
		e.globalStiffMat[i-1] = e.StiffnessGlobalMat(trail.T, lead.T)
	}

	c <- *e
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
	iStart := a[i].Geometry.Start
	jStart := a[j].Geometry.Start
	if pos := iStart.Compare(jStart); pos != 0 {
		return pos < 0
	}

	iEnd := a[i].Geometry.End
	jEnd := a[j].Geometry.End
	return iEnd.Compare(jEnd) < 0
}
