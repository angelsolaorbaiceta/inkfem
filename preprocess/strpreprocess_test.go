package preprocess

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
)

func TestStartNodesDofs(t *testing.T) {
	str := makeStructure()
	assignDof(&str)

	if dofs := str.Elements[0].Nodes[0].DegreesOfFreedomNum(); dofs != [3]int{0, 1, 2} {
		t.Errorf("Structural node expected to have DOFs [0 1 2], but found %v", dofs)
	}
	if dofs := str.Elements[1].Nodes[0].DegreesOfFreedomNum(); dofs != [3]int{0, 1, 9} {
		t.Errorf("Structural node expected to have DOFs [0 1 9], but found %v", dofs)
	}
}

func TestMiddleNodesDofs(t *testing.T) {
	str := makeStructure()
	assignDof(&str)

	if dofs := str.Elements[0].Nodes[1].DegreesOfFreedomNum(); dofs != [3]int{3, 4, 5} {
		t.Errorf("Structural node expected to have DOFs [3 4 5], but found %v", dofs)
	}
	if dofs := str.Elements[1].Nodes[1].DegreesOfFreedomNum(); dofs != [3]int{10, 11, 12} {
		t.Errorf("Structural node expected to have DOFs [10 11 12], but found %v", dofs)
	}
}

func TestEndNodesDofs(t *testing.T) {
	str := makeStructure()
	assignDof(&str)

	if dofs := str.Elements[0].Nodes[2].DegreesOfFreedomNum(); dofs != [3]int{6, 7, 8} {
		t.Errorf("Structural node expected to have DOFs [6 7 8], but found %v", dofs)
	}
	if dofs := str.Elements[1].Nodes[2].DegreesOfFreedomNum(); dofs != [3]int{13, 14, 15} {
		t.Errorf("Structural node expected to have DOFs [13 14 15], but found %v", dofs)
	}
}

func TestDofsCount(t *testing.T) {
	str := makeStructure()
	dofsCount := assignDof(&str)

	if dofsCount != 16 {
		t.Errorf("Sliced structure expected to have 16 degrees of freedom, but had %d", dofsCount)
	}
}

/* Utils */
func makeStructure() Structure {
	var (
		nA = structure.MakeFreeNodeFromProjs(1, 0, 0)
		nB = structure.MakeFreeNodeFromProjs(2, 0, 100)
		nC = structure.MakeFreeNodeFromProjs(3, 100, 0)

		elemOrigA = structure.MakeElement(
			1, nA, nB, structure.MakeFullConstraint(), structure.MakeFullConstraint(),
			structure.MakeUnitMaterial(), structure.MakeUnitSection(), []load.Load{},
		)
		elemOrigB = structure.MakeElement(
			2, nA, nC, structure.MakeDispConstraint(), structure.MakeFullConstraint(),
			structure.MakeUnitMaterial(), structure.MakeUnitSection(), []load.Load{},
		)
	)

	return Structure{
		Nodes: map[int]*structure.Node{nA.Id: nA, nB.Id: nB, nC.Id: nC},
		Elements: []*Element{
			MakeElement(elemOrigA, []*Node{
				MakeUnloadedNode(inkgeom.MinT, nA.Position),
				MakeUnloadedNode(inkgeom.MakeTParam(0.5), inkgeom.MakePoint(0, 50)),
				MakeUnloadedNode(inkgeom.MaxT, nB.Position)}),
			MakeElement(elemOrigB, []*Node{
				MakeUnloadedNode(inkgeom.MinT, nA.Position),
				MakeUnloadedNode(inkgeom.MakeTParam(0.5), inkgeom.MakePoint(50, 0)),
				MakeUnloadedNode(inkgeom.MaxT, nC.Position)}),
		},
	}
}
