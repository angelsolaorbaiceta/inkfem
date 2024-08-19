package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeSVLE(t *testing.T) {
	t.Run("create a new single value linear equation from the origin and another point", func(t *testing.T) {
		svle, err := MakeSVLEFromPoints(0, 0, 1, 1)

		assert.Nil(t, err)
		assert.Equal(t, 1.0, svle.Slope())
		assert.Equal(t, 0.0, svle.YIntercept())
	})

	t.Run("create a new single value linear equation from two points", func(t *testing.T) {
		svle, err := MakeSVLEFromPoints(1, 1, 2, 3)

		assert.Nil(t, err)
		assert.Equal(t, 2.0, svle.Slope())
		assert.Equal(t, -1.0, svle.YIntercept())
	})

	t.Run("return an error if the two points have the same x value", func(t *testing.T) {
		_, err := MakeSVLEFromPoints(1, 1, 1, 3)

		assert.NotNil(t, err)
	})

	t.Run("is not a horizontal line", func(t *testing.T) {
		svle := MakeSVLE(2, 3)

		assert.False(t, svle.IsHorizontal())
	})

	t.Run("is a horizontal line", func(t *testing.T) {
		svle := MakeSVLE(0, 3)

		assert.True(t, svle.IsHorizontal())
	})

	t.Run("find Y value for a given X value", func(t *testing.T) {
		svle := MakeSVLE(2, 3)

		assert.Equal(t, 5.0, svle.YAt(1))
	})

	t.Run("find X value for a given Y value", func(t *testing.T) {
		svle := MakeSVLE(2, 3)
		x, err := svle.XAt(5)

		assert.Nil(t, err)
		assert.Equal(t, 1.0, x)
	})

	t.Run("return an error when trying to find X value for a horizontal line", func(t *testing.T) {
		svle := MakeSVLE(0, 3)
		_, err := svle.XAt(5)

		assert.NotNil(t, err)
	})
}
