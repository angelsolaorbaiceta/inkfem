package sol

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
)

func TestWriteSolution(t *testing.T) {
	var (
		sol    = inkio.MakeTestSolution()
		writer bytes.Buffer

		reactionsOffset = 1
		barsOffset      = reactionsOffset + 2
		gdxOffset       = barsOffset + 2
		gdyOffset       = gdxOffset + 4
		grzOffset       = gdyOffset + 4
	)

	Write(sol, &writer)
	var gotLines []string
	for _, line := range strings.Split(writer.String(), "\n") {
		if line != "" {
			gotLines = append(gotLines, line)
		}
	}

	fmt.Println(gotLines)

	t.Run("first line is always the header with the version", func(t *testing.T) {
		var (
			want = fmt.Sprintf("inkfem v%d.%d", sol.Metadata.MajorVersion, sol.Metadata.MinorVersion)
			got  = gotLines[0]
		)

		if got != want {
			t.Errorf("Want '%s', got '%s'", want, got)
		}
	})

	t.Run("then go the node reactions", func(t *testing.T) {
		var (
			wantHeader       = "|reactions| 1"
			wantReactPattern = `n1 -> -?[\d\.]+ -?[\d\.]+ -?[\d\.]+`
			gotReaction      = gotLines[reactionsOffset+1]
		)

		if got := gotLines[reactionsOffset]; got != wantHeader {
			t.Errorf("Want '%s', got '%s'", wantHeader, got)
		}
		if match, _ := regexp.MatchString(wantReactPattern, gotReaction); !match {
			t.Errorf("Want '%s', got '%s'", wantReactPattern, gotReaction)
		}
	})

	t.Run("then goes the bars", func(t *testing.T) {
		var (
			wantHeader = "|bars| 1"
			wantBar    = "b1 -> n1 { dx dy rz } n2 { dx dy rz } 'mat_yz' 'sec_xy'"
		)

		if got := gotLines[barsOffset]; got != wantHeader {
			t.Errorf("Want '%s', got '%s'", wantHeader, got)
		}
		if got := gotLines[barsOffset+1]; got != wantBar {
			t.Errorf("Want '%s', got '%s'", wantBar, got)
		}
	})

	t.Run("each bar has global X displacements", func(t *testing.T) {
		var (
			wantHeader = "__gdx__"
			wantGdx1   = "0.000000 : 0.000000"
			wantGdx2   = "0.500000 : 1.000000"
			wantGdx3   = "1.000000 : 3.000000"
		)

		if got := gotLines[gdxOffset]; got != wantHeader {
			t.Errorf("Want '%s', got '%s'", wantHeader, got)
		}
		if got := gotLines[gdxOffset+1]; got != wantGdx1 {
			t.Errorf("Want '%s', got '%s'", wantGdx1, got)
		}
		if got := gotLines[gdxOffset+2]; got != wantGdx2 {
			t.Errorf("Want '%s', got '%s'", wantGdx2, got)
		}
		if got := gotLines[gdxOffset+3]; got != wantGdx3 {
			t.Errorf("Want '%s', got '%s'", wantGdx3, got)
		}
	})

	t.Run("each bar has global Y displacements", func(t *testing.T) {
		var (
			wantHeader = "__gdy__"
			wantGdy1   = "0.000000 : 0.000000"
			wantGdy2   = "0.500000 : 2.000000"
			wantGdy3   = "1.000000 : 4.000000"
		)

		if got := gotLines[gdyOffset]; got != wantHeader {
			t.Errorf("Want '%s', got '%s'", wantHeader, got)
		}
		if got := gotLines[gdyOffset+1]; got != wantGdy1 {
			t.Errorf("Want '%s', got '%s'", wantGdy1, got)
		}
		if got := gotLines[gdyOffset+2]; got != wantGdy2 {
			t.Errorf("Want '%s', got '%s'", wantGdy2, got)
		}
		if got := gotLines[gdyOffset+3]; got != wantGdy3 {
			t.Errorf("Want '%s', got '%s'", wantGdy3, got)
		}
	})

	t.Run("each bar has global Z rotations", func(t *testing.T) {
		var (
			wantHeader = "__grz__"
			wantGrz1   = "0.000000 : 0.000000"
			wantGrz2   = "0.500000 : 0.500000"
			wantGrz3   = "1.000000 : 0.700000"
		)

		if got := gotLines[grzOffset]; got != wantHeader {
			t.Errorf("Want '%s', got '%s'", wantHeader, got)
		}
		if got := gotLines[grzOffset+1]; got != wantGrz1 {
			t.Errorf("Want '%s', got '%s'", wantGrz1, got)
		}
		if got := gotLines[grzOffset+2]; got != wantGrz2 {
			t.Errorf("Want '%s', got '%s'", wantGrz2, got)
		}
		if got := gotLines[grzOffset+3]; got != wantGrz3 {
			t.Errorf("Want '%s', got '%s'", wantGrz3, got)
		}
	})
}
