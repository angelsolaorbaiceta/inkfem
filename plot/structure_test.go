package plot

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

func TestStructureToSVG(t *testing.T) {
	var (
		nodeOne   = structure.MakeNodeAtPosition("n1", 0, 0, &structure.FullConstraint)
		nodeTwo   = structure.MakeFreeNodeAtPosition("n2", 0, 300)
		nodeThree = structure.MakeFreeNodeAtPosition("n3", 200, 300)

		barOne = structure.
			MakeElementBuilder("b1").
			WithStartNode(nodeOne, &structure.FullConstraint).
			WithEndNode(nodeTwo, &structure.FullConstraint).
			WithMaterial(structure.MakeUnitMaterial()).
			WithSection(structure.MakeUnitSection()).
			Build()
		barTwo = structure.
			MakeElementBuilder("b2").
			WithStartNode(nodeTwo, &structure.FullConstraint).
			WithEndNode(nodeThree, &structure.FullConstraint).
			WithMaterial(structure.MakeUnitMaterial()).
			WithSection(structure.MakeUnitSection()).
			Build()
	)

	strDefinition := &structure.Structure{
		Metadata: structure.StrMetadata{
			MajorVersion: 1,
			MinorVersion: 0,
		},
		Nodes: map[contracts.StrID]*structure.Node{
			"n1": nodeOne,
			"n2": nodeTwo,
			"n3": nodeThree,
		},
		Elements: []*structure.Element{barOne, barTwo},
	}

	t.Run("computes the right image size given the scale", func(t *testing.T) {
		var b bytes.Buffer

		StructureToSVG(strDefinition, StructurePlotOps{Scale: 5.0}, &b)
		got := b.String()

		fmt.Printf("---> %s\n", got)
	})
}
