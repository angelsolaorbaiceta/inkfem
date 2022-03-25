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
		barsOffset     = loadsOffset + 3
	)

	Write(str, &writer)
	var gotLines []string
	for _, line := range strings.Split(writer.String(), "\n") {
		if line != "" {
			gotLines = append(gotLines, line)
		}
	}

	t.Run("first line is always the header with the version", func(t *testing.T) {
		want := fmt.Sprintf("inkfem v%d.%d", str.Metadata.MajorVersion, str.Metadata.MinorVersion)
		if got := gotLines[0]; got != want {
			t.Errorf("want %s, got %s", want, got)
		}
	})

	t.Run("then go the nodes", func(t *testing.T) {
		var (
			wantHeader         = "|nodes| 2"
			wantNodeOnePattern = `n1 -> 0(\.[0]+)? 0(\.[0]+)? { dx dy rz }`
			wantNodeTwoPattern = `n2 -> 200(\.[0]+)? 0(\.[0]+)? { }`
		)

		if got := gotLines[nodesOffset]; got != wantHeader {
			t.Errorf("want %s, got %s", wantHeader, got)
		}

		// Order in which the nodes appear isn't guaranteed
		nodeLines := gotLines[nodesOffset+1] + " " + gotLines[nodesOffset+2]

		if match, _ := regexp.MatchString(wantNodeOnePattern, nodeLines); !match {
			t.Error("want node one")
		}
		if match, _ := regexp.MatchString(wantNodeTwoPattern, nodeLines); !match {
			t.Error("want node two")
		}
	})

	t.Run("then go the materials", func(t *testing.T) {
		var (
			wantHeader          = "|materials| 1"
			wantMaterialPattern = `'mat_yz' -> 1(\.[0]*)? 2(\.[0]*)? 3(\.[0]*)? 4(\.[0]*)? 5(\.[0]*)? 6(\.[0]*)?`
		)

		if got := gotLines[materiasOffset]; got != wantHeader {
			t.Errorf("want %s, got %s", wantHeader, got)
		}
		if matches, _ := regexp.MatchString(wantMaterialPattern, gotLines[materiasOffset+1]); !matches {
			t.Errorf("want material, got: %s", gotLines[materiasOffset+1])
		}
	})

	t.Run("then go the sections", func(t *testing.T) {
		var (
			wantHeader         = "|sections| 1"
			wantSectionPattern = `'sec_xy' -> 1(\.[0]*)? 2(\.[0]*)? 3(\.[0]*)? 4(\.[0]*)? 5(\.[0]*)?`
		)

		if got := gotLines[sectionsOffset]; got != wantHeader {
			t.Errorf("want %s, got %s", wantHeader, got)
		}
		if matches, _ := regexp.MatchString(wantSectionPattern, gotLines[sectionsOffset+1]); !matches {
			t.Errorf("want section, got: %s", gotLines[sectionsOffset+1])
		}
	})

	t.Run("then go the loads", func(t *testing.T) {
		var (
			wantHeader          = "|loads| 2"
			wantConcLoadPattern = `fx lc b1 0.5[0]* -50.6[0]*`
			wantDistLoadPattern = `fy gd b1 0(\.[0]+)? 20.4[0]* 1(\.[0]+)? 40.5[0]*`
		)

		if got := gotLines[loadsOffset]; got != wantHeader {
			t.Errorf("want %s, got %s", wantHeader, got)
		}
		if matches, _ := regexp.MatchString(wantConcLoadPattern, gotLines[loadsOffset+1]); !matches {
			t.Errorf("want concentrated load, got: %s", gotLines[loadsOffset+1])
		}
		if matches, _ := regexp.MatchString(wantDistLoadPattern, gotLines[loadsOffset+2]); !matches {
			t.Errorf("want distributed load, got %s", gotLines[loadsOffset+2])
		}
	})

	t.Run("lastly go the bars", func(t *testing.T) {
		var (
			wantHeader = "|bars| 1"
			wantBar    = "b1 -> n1 { dx dy rz } n2 { dx dy rz } 'mat_yz' 'sec_xy'"
		)

		if got := gotLines[barsOffset]; got != wantHeader {
			t.Errorf("want %s, got %s", wantHeader, got)
		}
		if got := gotLines[barsOffset+1]; got != wantBar {
			t.Errorf("want bar, got %s", got)
		}
	})
}
