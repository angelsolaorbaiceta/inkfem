package load

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/stretchr/testify/assert"
)

func TestLoadIsNodal(t *testing.T) {
	t.Run("concentrated in the start position is nodal", func(t *testing.T) {
		load := MakeConcentrated(FX, true, nums.MinT, 45.0)
		assert.True(t, load.IsNodal())
	})

	t.Run("concentrated in the end position is nodal", func(t *testing.T) {
		load := MakeConcentrated(FX, true, nums.MaxT, 45.0)
		assert.True(t, load.IsNodal())
	})
}
