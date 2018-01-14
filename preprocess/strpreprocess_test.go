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

	expectedDofs := [3]int{0, 1, 2}
	actualDofs := [2][3]int{
		str.Nodes[1].DegreesOfFreedomNum(),
		str.Elements[0].Nodes[0].DegreesOfFreedomNum(),
	}

	for _, dofs := range actualDofs {
		if dofs != expectedDofs {
			t.Errorf("Structural node expected to have DOFs %v, but found %v", expectedDofs, dofs)
		}
	}
}

/* Utils */
func makeStructure() Structure {
	var (
		nA = structure.MakeFreeNodeFromProjs(1, 0, 0)
		nB = structure.MakeFreeNodeFromProjs(2, 0, 100)
		nC = structure.MakeFreeNodeFromProjs(3, 100, 0)

		elemOrigA = structure.MakeElement(1, nA, nB, structure.MakeFullConstraint(), structure.MakeFullConstraint(), structure.MakeUnitMaterial(), structure.MakeUnitSection(), []load.Load{})
		elemOrigB = structure.MakeElement(2, nA, nC, structure.MakeDispConstraint(), structure.MakeFullConstraint(), structure.MakeUnitMaterial(), structure.MakeUnitSection(), []load.Load{})
	)

	return Structure{
		Nodes: map[int]structure.Node{nA.Id: nA, nB.Id: nB, nC.Id: nC},
		Elements: []Element{
			MakeElement(elemOrigA, []Node{
				MakeUnloadedNode(inkgeom.MIN_T, nA.Position),
				MakeUnloadedNode(inkgeom.MakeTParam(0.5), inkgeom.MakePoint(0, 50)),
				MakeUnloadedNode(inkgeom.MAX_T, nB.Position)}),
			MakeElement(elemOrigB, []Node{
				MakeUnloadedNode(inkgeom.MIN_T, nA.Position),
				MakeUnloadedNode(inkgeom.MakeTParam(0.5), inkgeom.MakePoint(50, 0)),
				MakeUnloadedNode(inkgeom.MAX_T, nC.Position)}),
		},
	}
}
