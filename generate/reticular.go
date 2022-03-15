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
	)

	return structure.Make(
		// TODO: read from version file
		structure.StrMetadata{
			MajorVersion: 1,
			MinorVersion: 0,
		},
		nodes,
		[]*structure.Element{},
	)
}

func generateNodes(params ReticStructureParams) map[contracts.StrID]*structure.Node {
	var (
		nodes     = make(map[contracts.StrID]*structure.Node)
		nodeIndex = 0
		nodeId    string
	)

	for i := 0; i < params.Spans+1; i++ {
		for j := 0; j < params.Levels+1; j++ {
			nodeId = fmt.Sprint(nodeIndex + 1)

			nodes[nodeId] = structure.MakeNodeAtPosition(
				nodeId,
				0.0,
				0.0,
				&structure.FullConstraint,
			)
		}
	}

	return nodes
}
