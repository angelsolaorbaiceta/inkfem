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
Element after slicing original structural element.
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

/* <-- Methods --> */

/*
GlobalStiffMatrixAt returns the global stiffness matrix at position i, that is, between
nodes i and i + 1.
*/
func (e Element) GlobalStiffMatrixAt(i int) mat.ReadOnlyMatrix {
	var (
		trail = e.Nodes[i]
		lead  = e.Nodes[i+1]
	)

	return e.StiffnessGlobalMat(trail.T, lead.T)
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
