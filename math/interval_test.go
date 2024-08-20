package math

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/stretchr/testify/assert"
)

func TestIntCloseInterval(t *testing.T) {
	numbers := []int{1, 2, 3, 4, 5}

	t.Run("creates a new IntCloseInterval from a list of integers", func(t *testing.T) {
		interval := MakeIntCloseInterval(numbers)

		assert.Equal(t, 1, interval.Min())
		assert.Equal(t, 5, interval.Max())
	})

	t.Run("returns the value at a given parameter t in the interval", func(t *testing.T) {
		interval := MakeIntCloseInterval(numbers)

		assert.Equal(t, 1, interval.ValueAt(nums.MakeTParam(0.0)))
		assert.Equal(t, 1, interval.ValueAt(nums.MakeTParam(0.1)))
		assert.Equal(t, 2, interval.ValueAt(nums.MakeTParam(0.2)))
		assert.Equal(t, 2, interval.ValueAt(nums.MakeTParam(0.3)))
		assert.Equal(t, 3, interval.ValueAt(nums.MakeTParam(0.4)))
		assert.Equal(t, 3, interval.ValueAt(nums.MakeTParam(0.5)))
		assert.Equal(t, 3, interval.ValueAt(nums.MakeTParam(0.6)))
		assert.Equal(t, 4, interval.ValueAt(nums.MakeTParam(0.7)))
		assert.Equal(t, 4, interval.ValueAt(nums.MakeTParam(0.8)))
		assert.Equal(t, 5, interval.ValueAt(nums.MakeTParam(0.9)))
		assert.Equal(t, 5, interval.ValueAt(nums.MakeTParam(1.0)))
	})
}
