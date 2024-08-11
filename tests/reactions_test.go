package tests

import (
	"math"
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/build"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/stretchr/testify/assert"
)

func TestTwoElementsCantileverReactions(t *testing.T) {
	build.ReadBuildInfo()

	var (
		qyValue           = -200.0
		fyValue           = -4000.0
		elTwoAngle        = math.Atan(0.5)
		elTwoLength       = math.Sqrt(math.Pow(length, 2) + math.Pow(0.5*length, 2))
		str               = makeTwoElementsCantileverReactionsStructure(qyValue, fyValue)
		sol               = solveStructure(str)
		reactions         = sol.NodeReactions()["n2"]
		reactionsErrorEps = 4.0
	)

	t.Run("Fx reaction", func(t *testing.T) {
		want := fyValue * math.Sin(elTwoAngle)
		got := reactions.Fx()

		assert.True(t, nums.FloatsEqualEps(got, want, reactionsErrorEps))
	})

	t.Run("Fy reaction", func(t *testing.T) {
		want := -fyValue*math.Cos(elTwoAngle) - qyValue*length
		got := reactions.Fy()

		assert.True(t, nums.FloatsEqualEps(got, want, reactionsErrorEps))
	})

	t.Run("Mz reaction", func(t *testing.T) {
		want := -fyValue*elTwoLength + qyValue*length*length*0.5
		got := reactions.Mz()

		assert.True(t, nums.FloatsEqualEps(got, want, reactionsErrorEps))
	})
}
