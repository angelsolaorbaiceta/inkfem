package pre

import (
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

func makeTestOriginalStructure() *structure.Structure {
	var (
		metadata = structure.StrMetadata{
			MajorVersion: 2,
			MinorVersion: 3,
		}
		nodeOne   = structure.MakeNodeAtPosition("n1", 0, 0, &structure.FullConstraint)
		nodeTwo   = structure.MakeFreeNodeAtPosition("n2", 200, 0)
		nodesById = map[contracts.StrID]*structure.Node{
			nodeOne.GetID(): nodeOne,
			nodeTwo.GetID(): nodeTwo,
		}
		element = structure.MakeElementBuilder("b1").
			WithStartNode(nodeOne, &structure.FullConstraint).
			WithEndNode(nodeTwo, &structure.FullConstraint).
			WithSection(structure.MakeUnitSection()).
			WithMaterial(structure.MakeUnitMaterial()).
			Build()
	)

	return structure.Make(metadata, nodesById, []*structure.Element{element})
}

func makeTestPreprocessedStructure() *preprocess.Structure {
	var (
		original        = makeTestOriginalStructure()
		originalElement = original.Elements[0]
		preNodes        = []*preprocess.Node{
			preprocess.MakeNode(nums.MinT, originalElement.StartPoint(), 10, 20, 30),
			preprocess.MakeNode(nums.HalfT, originalElement.PointAt(nums.HalfT), 11, 21, 31),
			preprocess.MakeNode(nums.MaxT, originalElement.EndPoint(), 12, 22, 32),
		}
		elements = []*preprocess.Element{
			preprocess.MakeElement(originalElement, preNodes),
		}
	)

	// Add left load to first node
	preNodes[0].AddLocalLeftLoad(5, 10, 15)

	// Add right load to last node
	preNodes[2].AddLocalRightLoad(-5, -10, -15)

	return preprocess.MakeStructure(original.Metadata, original.NodesById, elements)
}
