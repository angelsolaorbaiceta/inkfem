package plot

import (
	"bytes"
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/stretchr/testify/assert"
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

		strDefinition = structure.Make(
			structure.StrMetadata{
				MajorVersion: 1,
				MinorVersion: 0,
			},
			map[contracts.StrID]*structure.Node{
				"n1": nodeOne,
				"n2": nodeTwo,
				"n3": nodeThree,
			},
			[]*structure.Element{barOne, barTwo},
		)

		plotConfig = DefaultPlotConfig()
		plotOps    = &StructurePlotOps{
			Scale:     5.0,
			MinMargin: 100,
		}

		b   bytes.Buffer
		got string
	)

	StructureToSVG(strDefinition, plotOps, plotConfig, &b)
	got = b.String()

	t.Run("computes the right image size given the scale and margin", func(t *testing.T) {
		// structure has a width of 200, times the scale plus the two lateral margins = 1200px
		// structure has a height of 300, times the scale plus the two vertical margins = 1700px
		var (
			wantWidthPattern  = "width=\"1200\""
			wantHeightPattern = "height=\"1700\""
		)

		assert.Regexp(t, wantWidthPattern, got)
		assert.Regexp(t, wantHeightPattern, got)
	})

	t.Run("applies an affine transformation", func(t *testing.T) {
		// sx = 5, sy = -5, tx = 100, ty = 1700 - 100 = 1600
		wantTransformPattern := "matrix\\(5(\\.[0]+)?,0,0,-5(\\.[0]+)?,100(\\.[0]+)?,1600(\\.[0]+)?\\)"
		assert.Regexp(t, wantTransformPattern, got)
	})

	t.Run("draws the bars", func(t *testing.T) {
		var (
			wantBarOnePattern = "<line x1=\"0\" y1=\"0\" x2=\"0\" y2=\"300\" id=\"bar__b1\""
			wantBarTwoPattern = "<line x1=\"0\" y1=\"300\" x2=\"200\" y2=\"300\" id=\"bar__b2\""
		)

		assert.Regexp(t, wantBarOnePattern, got)
		assert.Regexp(t, wantBarTwoPattern, got)
	})

	t.Run("draws the nodes", func(t *testing.T) {
		var (
			wantNodeOnePattern   = "<circle cx=\"0\" cy=\"0\" r=\"10\" id=\"node__n1\""
			wantNodeTwoPattern   = "<circle cx=\"0\" cy=\"300\" r=\"10\" id=\"node__n2\""
			wantNodeThreePattern = "<circle cx=\"200\" cy=\"300\" r=\"10\" id=\"node__n3\""
		)

		assert.Regexp(t, wantNodeOnePattern, got)
		assert.Regexp(t, wantNodeTwoPattern, got)
		assert.Regexp(t, wantNodeThreePattern, got)
	})
}
