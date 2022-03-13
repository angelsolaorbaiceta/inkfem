package io

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

func TestWritePreprocessedStructure(t *testing.T) {
	var (
		metadata = structure.StrMetadata{
			MajorVersion: 2,
			MinorVersion: 3,
		}
		nodeOne   = structure.MakeNodeAtPosition("n1", 0, 0, &structure.FullConstraint)
		nodeTwo   = structure.MakeFreeNodeAtPosition("n2", 200, 0)
		nodesById = structure.MakeNodesById(map[contracts.StrID]*structure.Node{
			nodeOne.GetID(): nodeOne,
			nodeTwo.GetID(): nodeTwo,
		})
		originalElement = structure.MakeElementBuilder("b1").
				WithStartNode(nodeOne, &structure.FullConstraint).
				WithEndNode(nodeTwo, &structure.FullConstraint).
				WithSection(structure.MakeUnitSection()).
				WithMaterial(structure.MakeUnitMaterial()).
				Build()
		elements = []*preprocess.Element{
			preprocess.MakeElement(
				originalElement,
				[]*preprocess.Node{
					preprocess.MakeNode(nums.MinT, originalElement.StartPoint(), 10, 20, 30),
					preprocess.MakeNode(nums.HalfT, originalElement.PointAt(nums.HalfT), 11, 21, 31),
					preprocess.MakeNode(nums.MaxT, originalElement.EndPoint(), 12, 22, 32),
				},
			),
		}
		str    = preprocess.MakeStructure(metadata, nodesById, elements)
		writer bytes.Buffer
	)

	WritePreprocessedStructure(str, &writer)
	gotLines := strings.Split(writer.String(), "\n")

	fmt.Println(writer.String())

	t.Run("first line is always the header with the version", func(t *testing.T) {
		want := fmt.Sprintf("inkfem v%d.%d", metadata.MajorVersion, metadata.MinorVersion)
		if got := gotLines[0]; got != want {
			t.Errorf("Want '%s', got '%s'", want, got)
		}
	})

	t.Run("then goes the degrees of freedom count", func(t *testing.T) {
		// 3 nodes x 3dof = 9 total dofs
		want := "dof_count: 9"
		if got := gotLines[2]; got != want {
			t.Errorf("Want '%s', got '%s'", want, got)
		}
	})
}
