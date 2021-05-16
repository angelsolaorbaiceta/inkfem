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
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

func applyDistributedLoadsToNodes(nodes []*Node, element *structure.Element) {
	var (
		trailNode, leadNode *Node
		length, halfLength  float64
		avgLoadValVect      [3]float64
		elemRefFrame        = element.Geometry.RefFrame()
	)

	for i, j := 0, 1; j < len(nodes); i, j = i+1, j+1 {
		trailNode, leadNode = nodes[i], nodes[j]
		length = element.Geometry.LengthBetween(trailNode.T, leadNode.T)
		halfLength = 0.5 * length

		for _, load := range element.Loads {
			avgLoadValVect = load.AvgValueVectorBetween(trailNode.T, leadNode.T)

			if !load.IsInLocalCoords {
				localForces := elemRefFrame.ProjectVector(
					g2d.MakeVector(avgLoadValVect[0], avgLoadValVect[1]),
				)
				avgLoadValVect[0] = localForces.X
				avgLoadValVect[1] = localForces.Y
			}

			trailNode.AddLoad([3]float64{
				avgLoadValVect[0] * halfLength,
				avgLoadValVect[1] * halfLength,
				(avgLoadValVect[2] * halfLength) + (avgLoadValVect[1] * length * length / 12.0),
			})
			leadNode.AddLoad([3]float64{
				avgLoadValVect[0] * halfLength,
				avgLoadValVect[1] * halfLength,
				(avgLoadValVect[2] * halfLength) - (avgLoadValVect[1] * length * length / 12.0),
			})
		}
	}
}
