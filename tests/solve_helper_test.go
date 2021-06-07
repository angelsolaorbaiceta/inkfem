package tests

import (
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/process"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
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
	pre := preprocess.DoStructure(str)
	return process.Solve(pre, solveOptions)
}

func makeCantileverBeamStructure(
	concentratedLoads []*load.ConcentratedLoad,
	distributedLoads []*load.DistributedLoad,
) *structure.Structure {
	var (
		nodeOne = structure.MakeNode("fixed-node", g2d.MakePoint(0, 0), &structure.FullConstraint)
		nodeTwo = structure.MakeNode("free-node", g2d.MakePoint(length, 0), &structure.NilConstraint)
		beam    = structure.MakeElement(
			"beam",
			nodeOne,
			nodeTwo,
			&structure.FullConstraint,
			&structure.FullConstraint,
			material,
			section,
			concentratedLoads,
			distributedLoads,
		)
	)

	return &structure.Structure{
		Metadata: structure.StrMetadata{MajorVersion: 1, MinorVersion: 0},
		Nodes: map[contracts.StrID]*structure.Node{
			nodeOne.Id: nodeOne,
			nodeTwo.Id: nodeTwo,
		},
		Elements: []*structure.Element{beam},
	}
}

func makeAxialElementStructure(
	concentratedLoads []*load.ConcentratedLoad,
	distributedLoads []*load.DistributedLoad,
) *structure.Structure {
	var (
		nodeOne      = structure.MakeNode("fixed-node", g2d.MakePoint(0, 0), &structure.FullConstraint)
		nodeTwo      = structure.MakeNode("free-node", g2d.MakePoint(length, 0), &structure.NilConstraint)
		axialElement = structure.MakeElement(
			"axial-element",
			nodeOne,
			nodeTwo,
			&structure.FullConstraint,
			&structure.FullConstraint,
			material,
			section,
			concentratedLoads,
			distributedLoads,
		)
	)

	return &structure.Structure{
		Metadata: structure.StrMetadata{MajorVersion: 1, MinorVersion: 0},
		Nodes: map[contracts.StrID]*structure.Node{
			nodeOne.Id: nodeOne,
			nodeTwo.Id: nodeTwo,
		},
		Elements: []*structure.Element{axialElement},
	}
}

func makeTwoElementsCantileverReactionsStructure(distLoadVal, concLoadValue float64) *structure.Structure {
	var (
		nodeOne    = structure.MakeNode("n1", g2d.MakePoint(0, 0), &structure.NilConstraint)
		nodeTwo    = structure.MakeNode("n2", g2d.MakePoint(length, 0), &structure.FullConstraint)
		nodeThree  = structure.MakeNode("n3", g2d.MakePoint(2*length, 0.5*length), &structure.NilConstraint)
		elementOne = structure.MakeElement(
			"el-1",
			nodeOne,
			nodeTwo,
			&structure.FullConstraint,
			&structure.FullConstraint,
			material,
			section,
			noConcLoads,
			[]*load.DistributedLoad{
				load.MakeDistributed(load.FY, true, inkgeom.MinT, distLoadVal, inkgeom.MaxT, distLoadVal),
			},
		)
		elementTwo = structure.MakeElement(
			"el-2",
			nodeTwo,
			nodeThree,
			&structure.FullConstraint,
			&structure.FullConstraint,
			material,
			section,
			[]*load.ConcentratedLoad{
				load.MakeConcentrated(load.FY, true, inkgeom.MaxT, concLoadValue),
			},
			noDistLoads,
		)
	)

	return &structure.Structure{
		Metadata: structure.StrMetadata{MajorVersion: 1, MinorVersion: 0},
		Nodes: map[contracts.StrID]*structure.Node{
			nodeOne.Id:   nodeOne,
			nodeTwo.Id:   nodeTwo,
			nodeThree.Id: nodeThree,
		},
		Elements: []*structure.Element{elementOne, elementTwo},
	}
}
