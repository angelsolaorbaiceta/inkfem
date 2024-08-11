package def

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/stretchr/testify/assert"
)

// TODO: test preprocessed nodes
func TestDeserializeBars(t *testing.T) {
	var (
		lineOne = "1 -> 1{ dx dy rz } 2{ dx dy } 'mat' 'sec'"
		lineTwo = "2 -> 1{ dx dy } 3{ dx dy rz } 'mat' 'sec'"

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
		distributedLoads = structure.DistLoadsById{}
		data             = &structure.StructureData{
			Nodes:             nodes,
			Materials:         materials,
			Sections:          sections,
			ConcentratedLoads: concentratedLoads,
			DistributedLoads:  distributedLoads,
		}

		wantBarOneDTO = &DeserializedBarDTO{
			Id:           "1",
			StartNodeId:  nodes["1"].GetID(),
			StartLink:    &structure.FullConstraint,
			EndNodeId:    nodes["2"].GetID(),
			EndLink:      &structure.DispConstraint,
			MaterialName: "mat",
			SectionName:  "sec",
		}
		wantBarTwoDTO = &DeserializedBarDTO{
			Id:           "2",
			StartNodeId:  nodes["1"].GetID(),
			StartLink:    &structure.DispConstraint,
			EndNodeId:    nodes["3"].GetID(),
			EndLink:      &structure.FullConstraint,
			MaterialName: "mat",
			SectionName:  "sec",
		}

		wantBarOne = structure.MakeElementBuilder(
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

		wantBarTwo = structure.MakeElementBuilder(
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
		barOneDTO, _ = DeserializeBar(lineOne)
		barTwoDTO, _ = DeserializeBar(lineTwo)
		bars         = BarsFromDeserialization(
			[]*DeserializedBarDTO{barOneDTO, barTwoDTO},
			data,
		)
	)

	t.Run("Deserialize bars into a DTO", func(t *testing.T) {
		assert.Equal(t, wantBarOneDTO, barOneDTO)
		assert.Equal(t, wantBarTwoDTO, barTwoDTO)
	})

	t.Run("Bars concentrated loads", func(t *testing.T) {
		barOne, barTwo := bars[0], bars[1]

		assert.Equal(t, wantBarOne.ConcentratedLoads, barOne.ConcentratedLoads)
		assert.Equal(t, wantBarTwo.ConcentratedLoads, barTwo.ConcentratedLoads)
	})
}
