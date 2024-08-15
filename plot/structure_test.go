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

		plotOps = &StructurePlotOps{
			Scale:     5.0,
			MinMargin: 100,
		}
	)

	t.Run("computes the right image size given the scale and margin", func(t *testing.T) {
		var b bytes.Buffer

		StructureToSVG(strDefinition, plotOps, &b)
		// structure has a width of 200, times the scale plus the two lateral margins = 1200px
		// structure has a height of 300, times the scale plust the two vertical margins = 1700px
		var (
			got               = b.String()
			wantWidthPattern  = "width=\"1200\""
			wantHeightPattern = "height=\"1700\""
		)

		if match, err := regexp.MatchString(wantWidthPattern, got); !match || err != nil {
			t.Errorf("Want %s, but didn't find in:\n%s\n", wantWidthPattern, got)
		}

		if match, err := regexp.MatchString(wantHeightPattern, got); !match || err != nil {
			t.Errorf("Want %s, but didn't find in:\n%s\n", wantHeightPattern, got)
		}
	})

	t.Run("applies an affine transformation", func(t *testing.T) {
		var b bytes.Buffer

		StructureToSVG(strDefinition, plotOps, &b)
		// sx = 5, sy = -5, tx = 100, ty = 1700 - 100 = 1600
		var (
			got                  = b.String()
			wantTransformPattern = "matrix\\(5(\\.[0]+)?,0,0,-5(\\.[0]+)?,100(\\.[0]+)?,1600(\\.[0]+)?\\)"
		)

		if match, err := regexp.MatchString(wantTransformPattern, got); !match || err != nil {
			t.Errorf("Want %s, but didn't find in:\n%s\n", wantTransformPattern, got)
		}
	})

	t.Run("draws the bars", func(t *testing.T) {
		var b bytes.Buffer

		StructureToSVG(strDefinition, plotOps, &b)
		var (
			got               = b.String()
			wantBarOnePattern = "<line x1=\"0\" y1=\"0\" x2=\"0\" y2=\"300\""
			wantBarTwoPattern = "<line x1=\"0\" y1=\"300\" x2=\"200\" y2=\"300\""
		)

		if match, err := regexp.MatchString(wantBarOnePattern, got); !match || err != nil {
			t.Errorf("Want %s, but didn't find in:\n%s\n", wantBarOnePattern, got)
		}
		if match, err := regexp.MatchString(wantBarTwoPattern, got); !match || err != nil {
			t.Errorf("Want %s, but didn't find in:\n%s\n", wantBarTwoPattern, got)
		}
	})
}
