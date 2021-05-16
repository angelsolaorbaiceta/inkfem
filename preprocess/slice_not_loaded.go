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
	"github.com/angelsolaorbaiceta/inkgeom"
)

/*
Non axial elements which have no loads applied are sliced just by subdividing their
geometry into a given number of slices, so that the slices have the same length.
*/
func sliceElementWithoutLoads(e *structure.Element, slices int) *Element {
	tPos := inkgeom.SubTParamCompleteRangeTimes(slices)
	nodes := make([]*Node, len(tPos))

	for i := 0; i < len(tPos); i++ {
		nodes[i] = MakeUnloadedNode(tPos[i], e.PointAt(tPos[i]))
	}

	return MakeElement(e, nodes)
}
