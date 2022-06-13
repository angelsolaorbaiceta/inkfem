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

		reactionsOffset  = 1
		barsOffset       = reactionsOffset + 2
		gdxOffset        = barsOffset + 2
		gdyOffset        = gdxOffset + 4
		grzOffset        = gdyOffset + 4
		ldxOffset        = grzOffset + 4
		ldyOffset        = ldxOffset + 4
		lrzOffset        = ldyOffset + 4
		axialOffset      = lrzOffset + 4
		shearOffset      = axialOffset + 5
		bendingOffset    = shearOffset + 5
		bendStressOffset = bendingOffset + 5
	)

	Write(sol, &writer)
	var gotLines []string
	for _, line := range strings.Split(writer.String(), "\n") {
		if line != "" {
			gotLines = append(gotLines, line)
		}
	}

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
			wantHeader       = "|reactions|"
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
			wantHeader = "|bars|"
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
		wantGdx := []string{
			"__gdx__",
			"0.000000 : 0.000000",
			"0.500000 : 1.000000",
			"1.000000 : 3.000000",
		}

		for i := 0; i < len(wantGdx); i++ {
			if got := gotLines[gdxOffset+i]; got != wantGdx[i] {
				t.Errorf("Want '%s', got '%s'", wantGdx[i], got)
			}
		}
	})

	t.Run("each bar has global Y displacements", func(t *testing.T) {
		wantGdy := []string{
			"__gdy__",
			"0.000000 : 0.000000",
			"0.500000 : 2.000000",
			"1.000000 : 4.000000",
		}

		for i := 0; i < len(wantGdy); i++ {
			if got := gotLines[gdyOffset+i]; got != wantGdy[i] {
				t.Errorf("Want '%s', got '%s'", wantGdy[i], got)
			}
		}
	})

	t.Run("each bar has global Z rotations", func(t *testing.T) {
		wantGrz := []string{
			"__grz__",
			"0.000000 : 0.000000",
			"0.500000 : 0.500000",
			"1.000000 : 0.700000",
		}

		for i := 0; i < len(wantGrz); i++ {
			if got := gotLines[grzOffset+i]; got != wantGrz[i] {
				t.Errorf("Want '%s', got '%s'", wantGrz[i], got)
			}
		}
	})

	t.Run("each bar has local X displacements", func(t *testing.T) {
		wantLdx := []string{
			"__ldx__",
			"0.000000 : 0.000000",
			"0.500000 : 1.000000",
			"1.000000 : 3.000000",
		}

		for i := 0; i < len(wantLdx); i++ {
			if got := gotLines[ldxOffset+i]; got != wantLdx[i] {
				t.Errorf("Want '%s', got '%s'", wantLdx[i], got)
			}
		}
	})

	t.Run("each bar has local Y displacements", func(t *testing.T) {
		wantLdy := []string{
			"__ldy__",
			"0.000000 : 0.000000",
			"0.500000 : 2.000000",
			"1.000000 : 4.000000",
		}

		for i := 0; i < len(wantLdy); i++ {
			if got := gotLines[ldyOffset+i]; got != wantLdy[i] {
				t.Errorf("Want '%s', got '%s'", wantLdy[i], got)
			}
		}
	})

	t.Run("each bar has local Z rotations", func(t *testing.T) {
		wantLrz := []string{
			"__lrz__",
			"0.000000 : 0.000000",
			"0.500000 : 0.500000",
			"1.000000 : 0.700000",
		}

		for i := 0; i < len(wantLrz); i++ {
			if got := gotLines[lrzOffset+i]; got != wantLrz[i] {
				t.Errorf("Want '%s', got '%s'", wantLrz[i], got)
			}
		}
	})

	t.Run("each bar has local axial forces", func(t *testing.T) {
		wantAxial := []string{
			"__axial__",
			"0.000000 : 5.020000",
			"0.500000 : 0.020000",
			"0.500000 : 0.040000",
			"1.000000 : 5.040000",
		}

		for i := 0; i < len(wantAxial); i++ {
			if got := gotLines[axialOffset+i]; got != wantAxial[i] {
				t.Errorf("Want '%s', got '%s'", wantAxial[i], got)
			}
		}
	})

	t.Run("each bar has local shear forces", func(t *testing.T) {
		wantShear := []string{
			"__shear__",
			"0.000000 : -9.998896",
			"0.500000 : 0.001104",
			"0.500000 : 0.002784",
			"1.000000 : -9.997216",
		}

		for i := 0; i < 5; i++ {
			if got := gotLines[shearOffset+i]; got != wantShear[i] {
				t.Errorf("Want '%s', got '%s'", wantShear[i], got)
			}
		}
	})

	t.Run("each bar has local bending moments", func(t *testing.T) {
		wantBending := []string{
			"__bend__",
			"0.000000 : 14.964800",
			"0.500000 : 0.075200",
			"0.500000 : -0.131200",
			"1.000000 : 15.147200",
		}

		for i := 0; i < len(wantBending); i++ {
			if got := gotLines[bendingOffset+i]; got != wantBending[i] {
				t.Errorf("Want '%s', got '%s'", wantBending[i], got)
			}
		}
	})

	t.Run("each bar has axial stresses associated to the bending moments", func(t *testing.T) {
		wantBendStress := []string{
			"__bend_axial_stress__",
			"0.000000 : 3.741200",
			"0.500000 : 0.018800",
			"0.500000 : -0.032800",
			"1.000000 : 3.786800",
		}

		for i := 0; i < len(wantBendStress); i++ {
			if got := gotLines[bendStressOffset+i]; got != wantBendStress[i] {
				t.Errorf("Want '%s', got '%s'", wantBendStress[i], got)
			}
		}
	})
}
