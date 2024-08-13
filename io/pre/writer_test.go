package pre

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/stretchr/testify/assert"
)

func TestWritePreprocessedStructure(t *testing.T) {
	var (
		str             = inkio.MakeTestPreprocessedStructure()
		writer          bytes.Buffer
		nodesOffset     = 3
		materialsOffset = nodesOffset + 3
		sectionsOffset  = materialsOffset + 2
		barsOffset      = sectionsOffset + 2
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
		assert.Equal(t, want, gotLines[0])
	})

	t.Run("the second line is the degrees of freedom count", func(t *testing.T) {
		// 3 nodes x 3dof = 9 total dofs
		assert.Equal(t, "dof_count: 9", gotLines[1])
	})

	t.Run("the third line is the own weight inclusion", func(t *testing.T) {
		assert.Equal(t, "includes_own_weight: no", gotLines[2])
	})

	t.Run("then go the original nodes", func(t *testing.T) {
		var (
			wantHeader         = "|nodes|"
			wantNodeOnePattern = `n1 -> 0(\.[0]+)? 0(\.[0]+)? { } | \[6 7 8\]`
			wantNodeTwoPattern = `n2 -> 200(\.[0]+)? 0(\.[0]+)? { dx dy rz } | \[0 1 2\]`
		)

		assert.Equal(t, wantHeader, gotLines[nodesOffset])

		// Order in which the nodes appear isn't guaranteed
		nodeLines := gotLines[nodesOffset+1] + " " + gotLines[nodesOffset+2]

		assert.Regexp(t, wantNodeOnePattern, nodeLines)
		assert.Regexp(t, wantNodeTwoPattern, nodeLines)
	})

	t.Run("then go the materials", func(t *testing.T) {
		var (
			wantHeader          = "|materials|"
			wantMaterialPattern = `'mat_yz' -> 1(\.[0]*)? 2(\.[0]*)? 3(\.[0]*)? 4(\.[0]*)? 5(\.[0]*)? 6(\.[0]*)?`
		)

		assert.Equal(t, wantHeader, gotLines[materialsOffset])
		assert.Regexp(t, wantMaterialPattern, gotLines[materialsOffset+1])
	})

	t.Run("then go the sections", func(t *testing.T) {
		var (
			wantHeader         = "|sections|"
			wantSectionPattern = `'sec_xy' -> 1(\.[0]*)? 2(\.[0]*)? 3(\.[0]*)? 4(\.[0]*)? 5(\.[0]*)?`
		)

		assert.Equal(t, wantHeader, gotLines[sectionsOffset])
		assert.Regexp(t, wantSectionPattern, gotLines[sectionsOffset+1])
	})

	t.Run("lastly go the bars", func(t *testing.T) {
		var (
			wantHeader = "|bars|"
			wantBar    = "b1 -> n1 { dx dy rz } n2 { dx dy rz } 'mat_yz' 'sec_xy' >> 3"
		)

		assert.Equal(t, wantHeader, gotLines[barsOffset])
		assert.Equal(t, wantBar, gotLines[barsOffset+1])

		// first node
		var (
			wantFirstNodePattern      = `0(\.[0]+)? : 0(\.[0]+)? 0(\.[0]+)?`
			wantFirstNodeExtPattern   = `\s+ext\s+: {10(\.[0]+)? 20(\.[0]+)? 30(\.[0]+)?}`
			wantFirstNodeLeftPattern  = `\s+left\s+: {5(\.[0]+)? 10(\.[0]+)? 15(\.[0]+)?}`
			wantFirstNodeRightPattern = `\s+right\s+: {0(\.[0]+)? 0(\.[0]+)? 0(\.[0]+)?}`
			wantFirstNodeNetPattern   = `\s+net\s+: {15(\.[0]+)? 30(\.[0]+)? 45(\.[0]+)?}`
			wantFirstNodeDofPattern   = `\s+dof\s+: \[0 1 2\]`
		)
		assert.Regexp(t, wantFirstNodePattern, gotLines[barsOffset+2])
		assert.Regexp(t, wantFirstNodeExtPattern, gotLines[barsOffset+3])
		assert.Regexp(t, wantFirstNodeLeftPattern, gotLines[barsOffset+4])
		assert.Regexp(t, wantFirstNodeRightPattern, gotLines[barsOffset+5])
		assert.Regexp(t, wantFirstNodeNetPattern, gotLines[barsOffset+6])
		assert.Regexp(t, wantFirstNodeDofPattern, gotLines[barsOffset+7])

		// second node
		var (
			wantSecondNodePattern      = `0\.5[0]+ : 100(\.[0]+)? 0(\.[0]+)?`
			wantSecondNodeExtPattern   = `\s+ext\s+: {11(\.[0]+)? 21(\.[0]+)? 31(\.[0]+)?}`
			wantSecondNodeLeftPattern  = `\s+left\s+: {0(\.[0]+)? 0(\.[0]+)? 0(\.[0]+)?}`
			wantSecondNodeRightPattern = `\s+right\s+: {0(\.[0]+)? 0(\.[0]+)? 0(\.[0]+)?}`
			wantSecondNodeNetPattern   = `\s+net\s+: {11(\.[0]+)? 21(\.[0]+)? 31(\.[0]+)?}`
			wantSecondNodeDofPattern   = `\s+dof\s+: \[3 4 5\]`
		)
		assert.Regexp(t, wantSecondNodePattern, gotLines[barsOffset+8])
		assert.Regexp(t, wantSecondNodeExtPattern, gotLines[barsOffset+9])
		assert.Regexp(t, wantSecondNodeLeftPattern, gotLines[barsOffset+10])
		assert.Regexp(t, wantSecondNodeRightPattern, gotLines[barsOffset+11])
		assert.Regexp(t, wantSecondNodeNetPattern, gotLines[barsOffset+12])
		assert.Regexp(t, wantSecondNodeDofPattern, gotLines[barsOffset+13])

		// third node
		var (
			wantThirdNodePattern      = `1(\.[0]+)? : 200(\.[0]+)? 0(\.[0]+)?`
			wantThirdNodeExtPattern   = `\s+ext\s+: {12(\.[0]+)? 22(\.[0]+)? 32(\.[0]+)?}`
			wantThirdNodeLeftPattern  = `\s+left\s+: {0(\.[0]+)? 0(\.[0]+)? 0(\.[0]+)?}`
			wantThirdNodeRightPattern = `\s+right\s+: {-5(\.[0]+)? -10(\.[0]+)? -15(\.[0]+)?}`
			wantThirdNodeNetPattern   = `\s+net\s+: {7(\.[0]+)? 12(\.[0]+)? 17(\.[0]+)?}`
			wantThirdNodeDofPattern   = `\s+dof\s+: \[6 7 8\]`
		)
		assert.Regexp(t, wantThirdNodePattern, gotLines[barsOffset+14])
		assert.Regexp(t, wantThirdNodeExtPattern, gotLines[barsOffset+15])
		assert.Regexp(t, wantThirdNodeLeftPattern, gotLines[barsOffset+16])
		assert.Regexp(t, wantThirdNodeRightPattern, gotLines[barsOffset+17])
		assert.Regexp(t, wantThirdNodeNetPattern, gotLines[barsOffset+18])
		assert.Regexp(t, wantThirdNodeDofPattern, gotLines[barsOffset+19])
	})
}
