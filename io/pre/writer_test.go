package pre

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestWritePreprocessedStructure(t *testing.T) {
	var (
		str            = makeTestPreprocessedStructure()
		writer         bytes.Buffer
		nodesOffset    = 2
		materiasOffset = nodesOffset + 3
		sectionsOffset = materiasOffset
		barsOffset     = sectionsOffset + 1
	)

	Write(str, &writer)
	fmt.Println(writer.String())
	var gotLines []string
	for _, line := range strings.Split(writer.String(), "\n") {
		if line != "" {
			gotLines = append(gotLines, line)
		}
	}

	fmt.Println(writer.String())

	t.Run("first line is always the header with the version", func(t *testing.T) {
		want := fmt.Sprintf("inkfem v%d.%d", str.Metadata.MajorVersion, str.Metadata.MinorVersion)
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

		if got := gotLines[nodesOffset]; got != wantHeader {
			t.Errorf("Want '%s', got '%s'", wantHeader, got)
		}

		// Order in which the nodes appear isn't guaranteed
		nodeLines := gotLines[nodesOffset+1] + " " + gotLines[nodesOffset+2]

		if match, _ := regexp.MatchString(wantNodeOnePattern, nodeLines); !match {
			t.Error("Want node one")
		}
		if match, _ := regexp.MatchString(wantNodeTwoPattern, nodeLines); !match {
			t.Error("Want node two")
		}
	})

	t.Run("then go the materials", func(t *testing.T) {
		var (
			wantHeader = "|materials| 1"
			// wantMaterialPattern = "b1 -> n1 { dx dy rz } n2 { dx dy rz } 'unit_material' 'unit_section' >> 3"
		)

		if got := gotLines[materiasOffset]; got != wantHeader {
			t.Errorf("want '%s', got '%s'", wantHeader, got)
		}
	})

	t.Run("then go the sections", func(t *testing.T) {})

	t.Run("lastly go the bars", func(t *testing.T) {
		var (
			wantHeader = "|bars| 1"
			wantBar    = "b1 -> n1 { dx dy rz } n2 { dx dy rz } 'unit_material' 'unit_section' >> 3"
		)

		if got := gotLines[barsOffset]; got != wantHeader {
			t.Errorf("want '%s', got '%s'", wantHeader, got)
		}
		if got := gotLines[barsOffset+1]; got != wantBar {
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
		if matches, _ := regexp.MatchString(wantFirstNodePattern, gotLines[barsOffset+2]); !matches {
			t.Errorf("Want first node position: %s", gotLines[barsOffset+2])
		}
		if matches, _ := regexp.MatchString(wantFirstNodeLeftPattern, gotLines[barsOffset+3]); !matches {
			t.Errorf("Want first node left load: %s", gotLines[barsOffset+3])
		}
		if matches, _ := regexp.MatchString(wantFirstNodeRightPattern, gotLines[barsOffset+4]); !matches {
			t.Errorf("Want first node right load: %s", gotLines[barsOffset+4])
		}
		if matches, _ := regexp.MatchString(wantFirstNodeNetPattern, gotLines[barsOffset+5]); !matches {
			t.Errorf("Want first node net load: %s", gotLines[barsOffset+5])
		}
		if matches, _ := regexp.MatchString(wantFirstNodeDofPattern, gotLines[barsOffset+6]); !matches {
			t.Errorf("Want first node dofs: %s", gotLines[barsOffset+6])
		}

		// second node
		var (
			wantSecondNodePattern      = "0\\.5[0]+ : 100(\\.[0]+)? 0(\\.[0]+)?"
			wantSecondNodeLeftPattern  = "\\s+left\\s+: {0(\\.[0]+)? 0(\\.[0]+)? 0(\\.[0]+)?}"
			wantSecondNodeRightPattern = "\\s+right\\s+: {0(\\.[0]+)? 0(\\.[0]+)? 0(\\.[0]+)?}"
			wantSecondNodeNetPattern   = "\\s+net\\s+: {11(\\.[0]+)? 21(\\.[0]+)? 31(\\.[0]+)?}"
			wantSecondNodeDofPattern   = "\\s+dof\\s+: \\[3 4 5\\]"
		)
		if matches, _ := regexp.MatchString(wantSecondNodePattern, gotLines[barsOffset+7]); !matches {
			t.Errorf("Want second node position: %s", gotLines[barsOffset+7])
		}
		if matches, _ := regexp.MatchString(wantSecondNodeLeftPattern, gotLines[barsOffset+8]); !matches {
			t.Errorf("Want second node left load: %s", gotLines[barsOffset+8])
		}
		if matches, _ := regexp.MatchString(wantSecondNodeRightPattern, gotLines[barsOffset+9]); !matches {
			t.Errorf("Want second node right load: %s", gotLines[barsOffset+9])
		}
		if matches, _ := regexp.MatchString(wantSecondNodeNetPattern, gotLines[barsOffset+10]); !matches {
			t.Errorf("Want second node net load: %s", gotLines[barsOffset+10])
		}
		if matches, _ := regexp.MatchString(wantSecondNodeDofPattern, gotLines[barsOffset+11]); !matches {
			t.Errorf("Want second node dofs: %s", gotLines[barsOffset+11])
		}

		// third node
		var (
			wantThirdNodePattern      = "1(\\.[0]+)? : 200(\\.[0]+)? 0(\\.[0]+)?"
			wantThirdNodeLeftPattern  = "\\s+left\\s+: {0(\\.[0]+)? 0(\\.[0]+)? 0(\\.[0]+)?}"
			wantThirdNodeRightPattern = "\\s+right\\s+: {-5(\\.[0]+)? -10(\\.[0]+)? -15(\\.[0]+)?}"
			wantThirdNodeNetPattern   = "\\s+net\\s+: {7(\\.[0]+)? 12(\\.[0]+)? 17(\\.[0]+)?}"
			wantThirdNodeDofPattern   = "\\s+dof\\s+: \\[6 7 8\\]"
		)
		if matches, _ := regexp.MatchString(wantThirdNodePattern, gotLines[barsOffset+12]); !matches {
			t.Errorf("Want second node position: %s", gotLines[barsOffset+12])
		}
		if matches, _ := regexp.MatchString(wantThirdNodeLeftPattern, gotLines[barsOffset+13]); !matches {
			t.Errorf("Want second node left load: %s", gotLines[barsOffset+13])
		}
		if matches, _ := regexp.MatchString(wantThirdNodeRightPattern, gotLines[barsOffset+14]); !matches {
			t.Errorf("Want second node right load: %s", gotLines[barsOffset+14])
		}
		if matches, _ := regexp.MatchString(wantThirdNodeNetPattern, gotLines[barsOffset+15]); !matches {
			t.Errorf("Want second node net load: %s", gotLines[barsOffset+15])
		}
		if matches, _ := regexp.MatchString(wantThirdNodeDofPattern, gotLines[barsOffset+16]); !matches {
			t.Errorf("Want second node dofs: %s", gotLines[barsOffset+16])
		}
	})
}