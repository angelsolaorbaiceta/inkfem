package io

import (
	"bytes"
	"fmt"
	"regexp"
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
		preNodes = []*preprocess.Node{
			preprocess.MakeNode(nums.MinT, originalElement.StartPoint(), 10, 20, 30),
			preprocess.MakeNode(nums.HalfT, originalElement.PointAt(nums.HalfT), 11, 21, 31),
			preprocess.MakeNode(nums.MaxT, originalElement.EndPoint(), 12, 22, 32),
		}
		elements = []*preprocess.Element{
			preprocess.MakeElement(originalElement, preNodes),
		}
		str    = preprocess.MakeStructure(metadata, nodesById, elements)
		writer bytes.Buffer
	)

	// Add left load to first node
	preNodes[0].AddLocalLeftLoad(5, 10, 15)

	// Add right load to last node
	preNodes[2].AddLocalRightLoad(-5, -10, -15)

	WritePreprocessedStructure(str, &writer)

	var gotLines []string
	for _, line := range strings.Split(writer.String(), "\n") {
		if line != "" {
			gotLines = append(gotLines, line)
		}
	}

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
		if got := gotLines[1]; got != want {
			t.Errorf("Want '%s', got '%s'", want, got)
		}
	})

	t.Run("then go the original nodes", func(t *testing.T) {
		var (
			wantHeader         = "|nodes| 2"
			wantNodeOnePattern = "n1 -> 0(\\.[0]+)? 0(\\.[0]+)? { } | DOF: \\[6 7 8\\]"
			wantNodeTwoPattern = "n2 -> 200(\\.[0]+)? 0(\\.[0]+)? { dx dy rz } | DOF: \\[0 1 2\\]"
		)

		if got := gotLines[2]; got != wantHeader {
			t.Errorf("Want '%s', got '%s'", wantHeader, got)
		}

		// Order in which the nodes appear isn't guaranteed
		nodeLines := gotLines[3] + " " + gotLines[4]

		if match, _ := regexp.MatchString(wantNodeOnePattern, nodeLines); !match {
			t.Error("Want node one")
		}
		if match, _ := regexp.MatchString(wantNodeTwoPattern, nodeLines); !match {
			t.Error("Want node two")
		}
	})

	t.Run("lastly go the bars", func(t *testing.T) {
		var (
			wantHeader = "|bars| 1"
			wantBar    = "b1 -> n1 { dx dy rz } n2 { dx dy rz } 'unit_material' 'unit_section' >> 3"
		)

		if got := gotLines[5]; got != wantHeader {
			t.Errorf("want '%s', got '%s'", wantHeader, got)
		}
		if got := gotLines[6]; got != wantBar {
			t.Errorf("want '%s', got '%s'", wantBar, got)
		}

		// first node
		var (
			wantFirstNodePattern      = "0(\\.[0]+)? : 0(\\.[0]+)? 0(\\.[0]+)?"
			wantFirstNodeLeftPattern  = "\\s+left\\s+: {5(\\.[0]+)? 10(\\.[0]+)? 15(\\.[0]+)?}"
			wantFirstNodeRightPattern = "\\s+right\\s+: {0(\\.[0]+)? 0(\\.[0]+)? 0(\\.[0]+)?}"
			wantFirstNodeNetPattern   = "\\s+net\\s+: {15(\\.[0]+)? 30(\\.[0]+)? 45(\\.[0]+)?}"
			wantFirstNodeDofPattern   = "\\s+dof\\s+: \\[0 1 2\\]"
		)
		if matches, _ := regexp.MatchString(wantFirstNodePattern, gotLines[7]); !matches {
			t.Errorf("Want first node position: %s", gotLines[7])
		}
		if matches, _ := regexp.MatchString(wantFirstNodeLeftPattern, gotLines[8]); !matches {
			t.Errorf("Want first node left load: %s", gotLines[8])
		}
		if matches, _ := regexp.MatchString(wantFirstNodeRightPattern, gotLines[9]); !matches {
			t.Errorf("Want first node right load: %s", gotLines[9])
		}
		if matches, _ := regexp.MatchString(wantFirstNodeNetPattern, gotLines[10]); !matches {
			t.Errorf("Want first node net load: %s", gotLines[10])
		}
		if matches, _ := regexp.MatchString(wantFirstNodeDofPattern, gotLines[11]); !matches {
			t.Errorf("Want first node dofs: %s", gotLines[11])
		}

		// second node
		var (
			wantSecondNodePattern      = "0\\.5[0]+ : 100(\\.[0]+)? 0(\\.[0]+)?"
			wantSecondNodeLeftPattern  = "\\s+left\\s+: {0(\\.[0]+)? 0(\\.[0]+)? 0(\\.[0]+)?}"
			wantSecondNodeRightPattern = "\\s+right\\s+: {0(\\.[0]+)? 0(\\.[0]+)? 0(\\.[0]+)?}"
			wantSecondNodeNetPattern   = "\\s+net\\s+: {11(\\.[0]+)? 21(\\.[0]+)? 31(\\.[0]+)?}"
			wantSecondNodeDofPattern   = "\\s+dof\\s+: \\[3 4 5\\]"
		)
		if matches, _ := regexp.MatchString(wantSecondNodePattern, gotLines[12]); !matches {
			t.Errorf("Want second node position: %s", gotLines[12])
		}
		if matches, _ := regexp.MatchString(wantSecondNodeLeftPattern, gotLines[13]); !matches {
			t.Errorf("Want second node left load: %s", gotLines[13])
		}
		if matches, _ := regexp.MatchString(wantSecondNodeRightPattern, gotLines[14]); !matches {
			t.Errorf("Want second node right load: %s", gotLines[14])
		}
		if matches, _ := regexp.MatchString(wantSecondNodeNetPattern, gotLines[15]); !matches {
			t.Errorf("Want second node net load: %s", gotLines[15])
		}
		if matches, _ := regexp.MatchString(wantSecondNodeDofPattern, gotLines[16]); !matches {
			t.Errorf("Want second node dofs: %s", gotLines[16])
		}
	})
}
