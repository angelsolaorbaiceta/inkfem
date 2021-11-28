package load

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

func TestLoadIsNodal(t *testing.T) {
	t.Run("concentrated in the start position is nodal", func(t *testing.T) {
		load := MakeConcentrated(FX, true, nums.MinT, 45.0)

		if !load.IsNodal() {
			t.Error("Expected load to be nodal (t = 0.0)")
		}
	})

	t.Run("concentrated in the end position is nodal", func(t *testing.T) {
		load := MakeConcentrated(FX, true, nums.MaxT, 45.0)

		if !load.IsNodal() {
			t.Error("Expected load to be nodal (t = 1.0)")
		}
	})
}
