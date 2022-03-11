package tests

import (
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/process"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

var (
	noDistLoads = []*load.DistributedLoad{}
	noConcLoads = []*load.ConcentratedLoad{}
	material    = &structure.Material{
		Name:             "steel",
		Density:          0,
		YoungMod:         20e6,
		ShearMod:         0,
		PoissonRatio:     1,
		YieldStrength:    0,
		UltimateStrength: 0,
	}
	section = &structure.Section{
		Name:    "IPE 120",
		Area:    14,
		IStrong: 318,
		IWeak:   28,
		SStrong: 53,
		SWeak:   9,
	}
	length     = 100.0
	displError = 1e-5
)

func solveStructure(str *structure.Structure) *process.Solution {
	solveOptions := process.SolveOptions{
		SaveSysMatrixImage:    false,
		OutputPath:            "",
		SafeChecks:            true,
		MaxDisplacementsError: displError,
	}
	pre := preprocess.StructureModel(str)
	return process.Solve(pre, solveOptions)
}

func makeCantileverBeamStructure(
	concentratedLoads []*load.ConcentratedLoad,
	distributedLoads []*load.DistributedLoad,
) *structure.Structure {
	var (
		nodeOne = structure.MakeNode("fixed-node", g2d.MakePoint(0, 0), &structure.FullConstraint)
		nodeTwo = structure.MakeNode("free-node", g2d.MakePoint(length, 0), &structure.NilConstraint)
		beam    = structure.MakeElementBuilder(
			"beam",
		).WithStartNode(
			nodeOne, &structure.FullConstraint,
		).WithEndNode(
			nodeTwo, &structure.FullConstraint,
		).WithMaterial(
			material,
		).WithSection(
			section,
		).AddConcentratedLoads(
			concentratedLoads,
		).AddDistributedLoads(
			distributedLoads,
		).Build()
	)

	return structure.Make(
		structure.StrMetadata{MajorVersion: 1, MinorVersion: 0},
		map[contracts.StrID]*structure.Node{
			nodeOne.GetID(): nodeOne,
			nodeTwo.GetID(): nodeTwo,
		},
		[]*structure.Element{beam},
	)
}

func makeAxialElementStructure(
	concentratedLoads []*load.ConcentratedLoad,
	distributedLoads []*load.DistributedLoad,
) *structure.Structure {
	var (
		nodeOne      = structure.MakeNode("fixed-node", g2d.MakePoint(0, 0), &structure.FullConstraint)
		nodeTwo      = structure.MakeNode("free-node", g2d.MakePoint(length, 0), &structure.NilConstraint)
		axialElement = structure.MakeElementBuilder(
			"axial-element",
		).WithStartNode(
			nodeOne, &structure.FullConstraint,
		).WithEndNode(
			nodeTwo, &structure.FullConstraint,
		).WithMaterial(
			material,
		).WithSection(
			section,
		).AddConcentratedLoads(
			concentratedLoads,
		).AddDistributedLoads(
			distributedLoads,
		).Build()
	)

	return structure.Make(
		structure.StrMetadata{MajorVersion: 1, MinorVersion: 0},
		map[contracts.StrID]*structure.Node{
			nodeOne.GetID(): nodeOne,
			nodeTwo.GetID(): nodeTwo,
		},
		[]*structure.Element{axialElement},
	)
}

func makeTwoElementsCantileverReactionsStructure(distLoadVal, concLoadValue float64) *structure.Structure {
	var (
		nodeOne   = structure.MakeNode("n1", g2d.MakePoint(0, 0), &structure.NilConstraint)
		nodeTwo   = structure.MakeNode("n2", g2d.MakePoint(length, 0), &structure.FullConstraint)
		nodeThree = structure.MakeNode("n3", g2d.MakePoint(2*length, 0.5*length), &structure.NilConstraint)

		elementOne = structure.MakeElementBuilder(
			"el-1",
		).WithStartNode(
			nodeOne, &structure.FullConstraint,
		).WithEndNode(
			nodeTwo, &structure.FullConstraint,
		).WithMaterial(
			material,
		).WithSection(
			section,
		).AddDistributedLoads(
			[]*load.DistributedLoad{
				load.MakeDistributed(load.FY, true, nums.MinT, distLoadVal, nums.MaxT, distLoadVal),
			},
		).Build()

		elementTwo = structure.MakeElementBuilder(
			"el-2",
		).WithStartNode(
			nodeTwo,
			&structure.FullConstraint,
		).WithEndNode(
			nodeThree, &structure.FullConstraint,
		).WithMaterial(
			material,
		).WithSection(
			section,
		).AddConcentratedLoads(
			[]*load.ConcentratedLoad{
				load.MakeConcentrated(load.FY, true, nums.MaxT, concLoadValue),
			},
		).Build()
	)

	return structure.Make(
		structure.StrMetadata{MajorVersion: 1, MinorVersion: 0},
		map[contracts.StrID]*structure.Node{
			nodeOne.GetID():   nodeOne,
			nodeTwo.GetID():   nodeTwo,
			nodeThree.GetID(): nodeThree,
		},
		[]*structure.Element{elementOne, elementTwo},
	)
}
