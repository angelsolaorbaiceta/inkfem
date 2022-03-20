package io

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

// TODO: test preprocessed ndoes
func TestDeserializeBars(t *testing.T) {
	var (
		lines = []string{
			"1 -> 1{ dx dy rz } 2{ dx dy } 'mat' 'sec'",
			"2 -> 1{ dx dy } 3{ dx dy rz } 'mat' 'sec'",
		}
		nodes = map[contracts.StrID]*structure.Node{
			"1": structure.MakeFreeNodeAtPosition("1", 100.0, 200.0),
			"2": structure.MakeFreeNodeAtPosition("2", 300.0, 400.0),
			"3": structure.MakeFreeNodeAtPosition("3", 500.0, 600.0),
		}
		materials = map[string]*structure.Material{
			"mat": structure.MakeMaterial("mat", 1.1, 2.2, 3.3, 4.4, 5.5, 6.6),
		}
		sections = map[string]*structure.Section{
			"sec": structure.MakeSection("sec", 10.1, 20.2, 30.3, 40.4, 50.5),
		}
		concentratedLoads = structure.ConcLoadsById{
			"1": {load.MakeConcentrated(load.FY, true, nums.MinT, -50)},
			"2": {load.MakeConcentrated(load.MZ, true, nums.MaxT, -30)},
		}
		ownWeightLoadVal = -materials["mat"].Density * sections["sec"].Area
		ownWeightLoad    = []*load.DistributedLoad{
			load.MakeDistributed(load.FY, false, nums.MinT, ownWeightLoadVal, nums.MaxT, ownWeightLoadVal),
		}
		distributedLoads = structure.DistLoadsById{}
		data             = &structure.StructureData{
			Nodes:             nodes,
			Materials:         materials,
			Sections:          sections,
			ConcentratedLoads: concentratedLoads,
			DistributedLoads:  distributedLoads,
		}

		wantElOne = structure.MakeElementBuilder(
			"1",
		).WithStartNode(
			nodes["1"], &structure.FullConstraint,
		).WithEndNode(
			nodes["2"], &structure.DispConstraint,
		).WithMaterial(
			materials["mat"],
		).WithSection(
			sections["sec"],
		).AddDistributedLoads(
			distributedLoads["1"],
		).AddConcentratedLoads(
			concentratedLoads["1"],
		).Build()

		wantElTwo = structure.MakeElementBuilder(
			"2",
		).WithStartNode(
			nodes["1"], &structure.DispConstraint,
		).WithEndNode(
			nodes["3"], &structure.FullConstraint,
		).WithMaterial(
			materials["mat"],
		).WithSection(
			sections["sec"],
		).AddDistributedLoads(
			distributedLoads["2"],
		).AddConcentratedLoads(
			concentratedLoads["2"],
		).Build()
	)

	var (
		elOne, _ = DeserializeBar(lines[0], data, ReaderOptions{ShouldIncludeOwnWeight: true})
		elTwo, _ = DeserializeBar(lines[1], data, ReaderOptions{ShouldIncludeOwnWeight: true})
	)

	t.Run("Elements read", func(t *testing.T) {
		if !elOne.Equals(wantElOne) {
			t.Errorf("Expected element %v, got %v", wantElOne, elOne)
		}
		if !elTwo.Equals(wantElTwo) {
			t.Errorf("Expected element %v, got %v", wantElTwo, elTwo)
		}
	})

	t.Run("Elements concentrated loads", func(t *testing.T) {
		if !load.ConcentratedLoadsEqual(elOne.ConcentratedLoads, wantElOne.ConcentratedLoads) {
			t.Errorf(
				"Expected element concentrated loads %v, but got %v",
				wantElOne.ConcentratedLoads,
				elOne.ConcentratedLoads,
			)
		}
		if !load.ConcentratedLoadsEqual(elTwo.ConcentratedLoads, wantElTwo.ConcentratedLoads) {
			t.Errorf(
				"Expected element concentrated loads %v, but got %v",
				wantElTwo.ConcentratedLoads,
				elTwo.ConcentratedLoads,
			)
		}
	})

	t.Run("Elements distributed loads", func(t *testing.T) {
		if !load.DistributedLoadsEqual(elOne.DistributedLoads, ownWeightLoad) {
			t.Errorf(
				"Expected element distributed loads %v, but got %v",
				ownWeightLoad,
				elOne.DistributedLoads,
			)
		}
		if !load.DistributedLoadsEqual(elTwo.DistributedLoads, ownWeightLoad) {
			t.Errorf(
				"Expected element distributed loads %v, but got %v",
				ownWeightLoad,
				elTwo.DistributedLoads,
			)
		}
	})

}
