package plot

import (
	"bytes"
	"regexp"
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

		plotOps = StructurePlotOps{
			Scale:     5.0,
			MinMargin: 100,
		}
	)

	t.Run("computes the right image size given the scale and margin", func(t *testing.T) {
		var b bytes.Buffer

		StructureToSVG(strDefinition, plotOps, &b)
		var (
			got = b.String()
			// structure has a width of 200, times the scale plus the two lateral margins = 1200px
			wantWidthPattern = "width=\"1200\""
			// structure has a height of 300, times the scale plust the two vertical margins = 1700px
			wantHeightPattern = "height=\"1700\""
		)

		if match, err := regexp.MatchString(wantWidthPattern, got); !match || err != nil {
			t.Errorf("Want %s, but didn't find in:\n%s\n", wantWidthPattern, got)
		}

		if match, err := regexp.MatchString(wantHeightPattern, got); !match || err != nil {
			t.Errorf("Want %s, but didn't find in:\n%s\n", wantHeightPattern, got)
		}
	})
}
