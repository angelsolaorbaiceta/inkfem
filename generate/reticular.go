package generate

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

type ReticStructureParams struct {
	Spans    int
	Levels   int
	Span     float64
	Height   float64
	Section  *structure.Section
	Material *structure.Material
}

func Reticular(params ReticStructureParams) *structure.Structure {
	var (
		nodes = generateNodes(params)
		bars  = generateBars(params, nodes)
	)

	return structure.Make(
		// TODO: read from version file
		structure.StrMetadata{
			MajorVersion: 1,
			MinorVersion: 0,
		},
		nodes,
		bars,
	)
}

func generateNodes(params ReticStructureParams) map[contracts.StrID]*structure.Node {
	var (
		nodes      = make(map[contracts.StrID]*structure.Node)
		nodeIndex  = 0
		nodeId     contracts.StrID
		constraint *structure.Constraint
	)

	for i := 0; i < params.Levels+1; i++ {
		if i == 0 {
			constraint = &structure.FullConstraint
		} else {
			constraint = &structure.NilConstraint
		}

		for j := 0; j < params.Spans+1; j++ {
			nodeId = fmt.Sprint(nodeIndex + 1)

			nodes[nodeId] = structure.MakeNodeAtPosition(
				nodeId,
				float64(j)*params.Span,
				float64(i)*params.Height,
				constraint,
			)

			nodeIndex += 1
		}
	}

	return nodes
}

func generateBars(params ReticStructureParams, nodes map[contracts.StrID]*structure.Node) []*structure.Element {
	var (
		rows      = params.Levels + 1
		cols      = params.Spans + 1
		barsCount = 2*cols*rows - 2*cols - rows + 1
		bars      = make([]*structure.Element, barsCount)
		barIndex  = 0
	)

	var isLowestNodesRow = func(id int) bool {
		return id <= cols
	}

	var isRowsLastNode = func(id int) bool {
		return id%cols == 0
	}

	var isUpperNodesRow = func(id int) bool {
		return id > cols*(rows-1)
	}

	for i := 1; i <= len(nodes); i++ {
		if !isRowsLastNode(i) && !isLowestNodesRow(i) {
			bars[barIndex] = structure.MakeElementBuilder(fmt.Sprint(barIndex+1)).
				WithStartNode(nodes[fmt.Sprint(i)], &structure.FullConstraint).
				WithEndNode(nodes[fmt.Sprint(i+1)], &structure.FullConstraint).
				WithMaterial(params.Material).
				WithSection(params.Section).
				Build()

			barIndex += 1
		}

		if !isUpperNodesRow(i) {
			bars[barIndex] = structure.MakeElementBuilder(fmt.Sprint(barIndex+1)).
				WithStartNode(nodes[fmt.Sprint(i)], &structure.FullConstraint).
				WithEndNode(nodes[fmt.Sprint(i+cols)], &structure.FullConstraint).
				WithMaterial(params.Material).
				WithSection(params.Section).
				Build()

			barIndex += 1
		}
	}

	return bars
}
