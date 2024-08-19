package load

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/stretchr/testify/assert"
)

func TestDistLoadEquation(t *testing.T) {
	xAxisLength := 100.0

	t.Run("can't get the equation if start T and end T are equal", func(t *testing.T) {
		load := MakeDistributed(FX, true, nums.MakeTParam(0.2), 20.0, nums.MakeTParam(0.2), 40.0)
		_, err := load.AsEquation(xAxisLength, 2.0)

		assert.Error(t, err)
	})

	t.Run("get equation for a distributed load", func(t *testing.T) {
		var (
			load    = MakeDistributed(FX, true, nums.MakeTParam(0.2), 20.0, nums.MakeTParam(0.8), 50.0)
			eq, err = load.AsEquation(xAxisLength, 2.0)
			want    = math.MakeSVLE(1.0, 20.0)
		)

		assert.NoError(t, err)
		assert.True(t, want.Equals(*eq))
	})
}
