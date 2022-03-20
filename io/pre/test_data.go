package pre

import (
	"io"
	"strings"

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
		originalElement = original.Elements()[0]
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

	return preprocess.
		MakeStructure(original.Metadata, original.NodesById, elements).
		AssignDof()
}

func makePreprocessedReader() io.Reader {
	return strings.NewReader(`inkfem v2.3
	
	dof_count: 9
	
	|nodes| 2
	n1 -> 0.000000 0.000000 { dx dy rz } | [0 1 2]
	n2 -> 200.000000 0.000000 { } | [6 7 8]

	|materials| 1
	'unit_material' -> 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000

	|sections| 1
	'unit_section' -> 1.000000 1.000000 1.000000 1.000000 1.000000
	
	|bars| 1
	b1 -> n1 { dx dy rz } n2 { dx dy rz } 'unit_material' 'unit_section' >> 3
	0.000000 : 0.000000 0.000000
					ext   : {10.000000 20.000000 30.000000}
					left  : {5.000000 10.000000 15.000000}
					right : {0.000000 0.000000 0.000000}
					net   : {15.000000 30.000000 45.000000}
					dof   : [0 1 2]
	0.500000 : 100.000000 0.000000
					ext   : {11.000000 21.000000 31.000000}
					left  : {0.000000 0.000000 0.000000}
					right : {0.000000 0.000000 0.000000}
					net   : {11.000000 21.000000 31.000000}
					dof   : [3 4 5]
	1.000000 : 200.000000 0.000000
					ext   : {12.000000 22.000000 32.000000}
					left  : {0.000000 0.000000 0.000000}
					right : {-5.000000 -10.000000 -15.000000}
					net   : {7.000000 12.000000 17.000000}
					dof   : [6 7 8]
	`)
}
