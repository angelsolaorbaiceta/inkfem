package def

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
)

func TestWriteDefinition(t *testing.T) {
	var (
		str            = inkio.MakeTestOriginalStructure()
		writer         bytes.Buffer
		nodesOffset    = 1
		materiasOffset = nodesOffset + 3
		sectionsOffset = materiasOffset + 2
		loadsOffset    = sectionsOffset + 2
		// barsOffset     = sectionsOffset + 2
	)

	Write(str, &writer)
	fmt.Println(writer.String())
	var gotLines []string
	for _, line := range strings.Split(writer.String(), "\n") {
		if line != "" {
			gotLines = append(gotLines, line)
		}
	}

	t.Run("first line is always the header with the version", func(t *testing.T) {
		want := fmt.Sprintf("inkfem v%d.%d", str.Metadata.MajorVersion, str.Metadata.MinorVersion)
		if got := gotLines[0]; got != want {
			t.Errorf("Want '%s', got '%s'", want, got)
		}
	})

	t.Run("then go the nodes", func(t *testing.T) {
		var (
			wantHeader         = "|nodes| 2"
			wantNodeOnePattern = `n1 -> 0(\.[0]+)? 0(\.[0]+)? { dx dy rz }`
			wantNodeTwoPattern = `n2 -> 200(\.[0]+)? 0(\.[0]+)? { }`
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
			wantHeader          = "|materials| 1"
			wantMaterialPattern = "'unit_material' -> 1\\.[0]+ 1\\.[0]+ 1\\.[0]+ 1\\.[0]+ 1\\.[0]+ 1\\.[0]+"
		)

		if got := gotLines[materiasOffset]; got != wantHeader {
			t.Errorf("want '%s', got '%s'", wantHeader, got)
		}
		if matches, _ := regexp.MatchString(wantMaterialPattern, gotLines[materiasOffset+1]); !matches {
			t.Errorf("Want material, got: %s", gotLines[materiasOffset+1])
		}
	})

	t.Run("then go the sections", func(t *testing.T) {
		var (
			wantHeader         = "|sections| 1"
			wantSectionPattern = `'unit_section' -> 1\.[0]+ 1\.[0]+ 1\.[0]+ 1\.[0]+ 1\.[0]+`
		)

		if got := gotLines[sectionsOffset]; got != wantHeader {
			t.Errorf("want '%s', got '%s'", wantHeader, got)
		}
		if matches, _ := regexp.MatchString(wantSectionPattern, gotLines[sectionsOffset+1]); !matches {
			t.Errorf("Want section, got: %s", gotLines[sectionsOffset+1])
		}
	})

	t.Run("then go the loads", func(t *testing.T) {
		var (
			wantHeader          = "|loads| 2"
			wantConcLoadPattern = `fx lc b1 0.5[0]* -50.6[0]*`
		)

		if got := gotLines[loadsOffset]; got != wantHeader {
			t.Errorf("want '%s', got '%s'", wantHeader, got)
		}
		if matches, _ := regexp.MatchString(wantConcLoadPattern, gotLines[loadsOffset+1]); !matches {
			t.Errorf("Want load, got: %s", gotLines[loadsOffset+1])
		}
	})
}
