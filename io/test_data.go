package io

import (
	"io"
	"strings"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/process"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

func MakeTestOriginalStructure() *structure.Structure {
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
		section  = structure.MakeSection("sec_xy", 1, 2, 3, 4, 5)
		material = structure.MakeMaterial("mat_yz", 1, 2, 3, 4, 5, 6)
		concLoad = load.MakeConcentrated(load.FX, true, nums.HalfT, -50.6)
		distLoad = load.MakeDistributed(load.FY, false, nums.MinT, 20.4, nums.MaxT, 40.5)
		element  = structure.MakeElementBuilder("b1").
				WithStartNode(nodeOne, &structure.FullConstraint).
				WithEndNode(nodeTwo, &structure.FullConstraint).
				WithSection(section).
				WithMaterial(material).
				AddConcentratedLoad(concLoad).
				AddDistributedLoad(distLoad).
				Build()
	)

	return structure.Make(metadata, nodesById, []*structure.Element{element})
}

func MakeTestDefinitionReader() io.Reader {
	return strings.NewReader(`inkfem v2.3
	 |nodes|
	 n1 -> 0 0 {dx dy rz}
	 n2 -> 200 0 {}

	 |sections|
	 'sec_xy' -> 1 2 3 4 5

	 |materials|
	 'mat_yz' -> 1 2 3 4 5 6

	 |loads|
	 fx lc b1 0.5 -50.6
	 fy gd b1 0 20.4 1 40.5

	 |bars|
	 b1 -> n1 {dx dy rz} n2{dx dy rz} 'mat_yz' 'sec_xy'
	 `)
}

func MakeTestDefinitionReaderInverseOrder() io.Reader {
	return strings.NewReader(`inkfem v2.3
	 |bars|
	 b1 -> n1 {dx dy rz} n2{dx dy rz} 'mat_yz' 'sec_xy'

	 |loads|
	 fy gd b1 0 20.4 1 40.5
	 fx lc b1 0.5 -50.6

	 |materials|
	 'mat_yz' -> 1 2 3 4 5 6

	 |sections|
	 'sec_xy' -> 1 2 3 4 5

	 |nodes|
	 n2 -> 200 0 {}
	 n1 -> 0 0 {dx dy rz}
	 `)
}

func MakeTestPreprocessedStructure() *preprocess.Structure {
	var (
		original        = MakeTestOriginalStructure()
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

func MakeTestSolution() *process.Solution {
	var (
		preStructure  = MakeTestPreprocessedStructure()
		preElement    = preStructure.GetElementById("b1")
		displacements = vec.MakeWithValues([]float64{
			0.0, 0.0, 0.0, // First node is fixed
			1.0, 2.0, 0.5, // Second node
			3.0, 4.0, 0.7, // Third node
		})
		solElement = process.MakeElementSolution(preElement, displacements, 1e-5)
	)

	return process.MakeSolution(
		structure.StrMetadata{
			MajorVersion: 2,
			MinorVersion: 3,
		},
		preStructure.NodesById,
		[]*process.ElementSolution{solElement},
	)
}

func MakeTestPreprocessedReader() io.Reader {
	return strings.NewReader(`inkfem v2.3
	
	dof_count: 9
	
	|nodes|
	n1 -> 0.000000 0.000000 { dx dy rz } | [0 1 2]
	n2 -> 200.000000 0.000000 { } | [6 7 8]

	|materials|
	'mat_yz' -> 1.000000 2.000000 3.000000 4.000000 5.000000 6.000000

	|sections|
	'sec_xy' -> 1.000000 2.000000 3.000000 4.000000 5.000000
	
	|bars|
	b1 -> n1 { dx dy rz } n2 { dx dy rz } 'mat_yz' 'sec_xy' >> 3
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
