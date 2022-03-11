package preprocess

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

func TestStartNodesDofs(t *testing.T) {
	str := makeStructure()
	str.assignDof()

	if dofs := str.Elements[0].NodeAt(0).DegreesOfFreedomNum(); dofs != [3]int{0, 1, 2} {
		t.Errorf("Structural node expected to have DOFs [0 1 2], but found %v", dofs)
	}
	if dofs := str.Elements[1].NodeAt(0).DegreesOfFreedomNum(); dofs != [3]int{0, 1, 9} {
		t.Errorf("Structural node expected to have DOFs [0 1 9], but found %v", dofs)
	}
}

func TestMiddleNodesDofs(t *testing.T) {
	str := makeStructure()
	str.assignDof()

	if dofs := str.Elements[0].NodeAt(1).DegreesOfFreedomNum(); dofs != [3]int{3, 4, 5} {
		t.Errorf("Structural node expected to have DOFs [3 4 5], but found %v", dofs)
	}
	if dofs := str.Elements[1].NodeAt(1).DegreesOfFreedomNum(); dofs != [3]int{10, 11, 12} {
		t.Errorf("Structural node expected to have DOFs [10 11 12], but found %v", dofs)
	}
}

func TestEndNodesDofs(t *testing.T) {
	str := makeStructure()
	str.assignDof()

	if dofs := str.Elements[0].NodeAt(2).DegreesOfFreedomNum(); dofs != [3]int{6, 7, 8} {
		t.Errorf("Structural node expected to have DOFs [6 7 8], but found %v", dofs)
	}
	if dofs := str.Elements[1].NodeAt(2).DegreesOfFreedomNum(); dofs != [3]int{13, 14, 15} {
		t.Errorf("Structural node expected to have DOFs [13 14 15], but found %v", dofs)
	}
}

func TestDofsCount(t *testing.T) {
	str := makeStructure()
	str.assignDof()

	if count := str.DofsCount(); count != 16 {
		t.Errorf("Sliced structure expected to have 16 degrees of freedom, but had %d", count)
	}
}

/* Utils */
func makeStructure() *Structure {
	var (
		nA = structure.MakeFreeNodeAtPosition("1", 0, 0)
		nB = structure.MakeFreeNodeAtPosition("2", 0, 100)
		nC = structure.MakeFreeNodeAtPosition("3", 100, 0)

		elemOrigA = structure.MakeElementBuilder(
			"1",
		).WithStartNode(
			nA, &structure.FullConstraint,
		).WithEndNode(
			nB, &structure.FullConstraint,
		).WithMaterial(
			structure.MakeUnitMaterial(),
		).WithSection(
			structure.MakeUnitSection(),
		).Build()

		elemOrigB = structure.MakeElementBuilder(
			"2",
		).WithStartNode(
			nA, &structure.DispConstraint,
		).WithEndNode(
			nC, &structure.FullConstraint,
		).WithMaterial(
			structure.MakeUnitMaterial(),
		).WithSection(
			structure.MakeUnitSection(),
		).Build()
	)

	return &Structure{
		NodesById: structure.MakeNodesById(
			map[contracts.StrID]*structure.Node{
				nA.GetID(): nA,
				nB.GetID(): nB,
				nC.GetID(): nC,
			},
		),
		Elements: []*Element{
			MakeElement(elemOrigA, []*Node{
				MakeUnloadedNode(nums.MinT, nA.Position),
				MakeUnloadedNode(nums.MakeTParam(0.5), g2d.MakePoint(0, 50)),
				MakeUnloadedNode(nums.MaxT, nB.Position)}),
			MakeElement(elemOrigB, []*Node{
				MakeUnloadedNode(nums.MinT, nA.Position),
				MakeUnloadedNode(nums.MakeTParam(0.5), g2d.MakePoint(50, 0)),
				MakeUnloadedNode(nums.MaxT, nC.Position)}),
		},
	}
}
