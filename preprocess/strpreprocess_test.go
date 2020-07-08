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
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

func TestStartNodesDofs(t *testing.T) {
	str := makeStructure()
	assignDof(str)

	if dofs := str.Elements[0].Nodes[0].DegreesOfFreedomNum(); dofs != [3]int{0, 1, 2} {
		t.Errorf("Structural node expected to have DOFs [0 1 2], but found %v", dofs)
	}
	if dofs := str.Elements[1].Nodes[0].DegreesOfFreedomNum(); dofs != [3]int{0, 1, 9} {
		t.Errorf("Structural node expected to have DOFs [0 1 9], but found %v", dofs)
	}
}

func TestMiddleNodesDofs(t *testing.T) {
	str := makeStructure()
	assignDof(str)

	if dofs := str.Elements[0].Nodes[1].DegreesOfFreedomNum(); dofs != [3]int{3, 4, 5} {
		t.Errorf("Structural node expected to have DOFs [3 4 5], but found %v", dofs)
	}
	if dofs := str.Elements[1].Nodes[1].DegreesOfFreedomNum(); dofs != [3]int{10, 11, 12} {
		t.Errorf("Structural node expected to have DOFs [10 11 12], but found %v", dofs)
	}
}

func TestEndNodesDofs(t *testing.T) {
	str := makeStructure()
	assignDof(str)

	if dofs := str.Elements[0].Nodes[2].DegreesOfFreedomNum(); dofs != [3]int{6, 7, 8} {
		t.Errorf("Structural node expected to have DOFs [6 7 8], but found %v", dofs)
	}
	if dofs := str.Elements[1].Nodes[2].DegreesOfFreedomNum(); dofs != [3]int{13, 14, 15} {
		t.Errorf("Structural node expected to have DOFs [13 14 15], but found %v", dofs)
	}
}

func TestDofsCount(t *testing.T) {
	str := makeStructure()
	assignDof(str)

	if count := str.DofsCount; count != 16 {
		t.Errorf("Sliced structure expected to have 16 degrees of freedom, but had %d", count)
	}
}

/* Utils */
func makeStructure() *Structure {
	var (
		nA = structure.MakeFreeNodeAtPosition(1, 0, 0)
		nB = structure.MakeFreeNodeAtPosition(2, 0, 100)
		nC = structure.MakeFreeNodeAtPosition(3, 100, 0)

		elemOrigA = structure.MakeElement(
			1, nA, nB, structure.FullConstraint, structure.FullConstraint,
			structure.MakeUnitMaterial(), structure.MakeUnitSection(), []load.Load{},
		)
		elemOrigB = structure.MakeElement(
			2, nA, nC, structure.DispConstraint, structure.FullConstraint,
			structure.MakeUnitMaterial(), structure.MakeUnitSection(), []load.Load{},
		)
	)

	return &Structure{
		Nodes: map[int]*structure.Node{nA.Id: nA, nB.Id: nB, nC.Id: nC},
		Elements: []*Element{
			MakeElement(elemOrigA, []*Node{
				MakeUnloadedNode(inkgeom.MinT, nA.Position),
				MakeUnloadedNode(inkgeom.MakeTParam(0.5), g2d.MakePoint(0, 50)),
				MakeUnloadedNode(inkgeom.MaxT, nB.Position)}),
			MakeElement(elemOrigB, []*Node{
				MakeUnloadedNode(inkgeom.MinT, nA.Position),
				MakeUnloadedNode(inkgeom.MakeTParam(0.5), g2d.MakePoint(50, 0)),
				MakeUnloadedNode(inkgeom.MaxT, nC.Position)}),
		},
	}
}
