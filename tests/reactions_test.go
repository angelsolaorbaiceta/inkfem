package tests

import (
	"math"
	"testing"

	"github.com/angelsolaorbaiceta/inkgeom"
)

func TestTwoElementsCantileverReactions(t *testing.T) {
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

		if got := reactions.Fx(); !inkgeom.FloatsEqualEps(got, want, reactionsErrorEps) {
			t.Errorf("Expected Fx = %f, but got %f", want, got)
		}
	})

	t.Run("Fy reaction", func(t *testing.T) {
		want := -fyValue*math.Cos(elTwoAngle) - qyValue*length

		if got := reactions.Fy(); !inkgeom.FloatsEqualEps(got, want, reactionsErrorEps) {
			t.Errorf("Expected Fy = %f, but got %f", want, got)
		}
	})

	t.Run("Mz reaction", func(t *testing.T) {
		want := -fyValue*elTwoLength + qyValue*length*length*0.5

		if got := reactions.Mz(); !inkgeom.FloatsEqualEps(got, want, reactionsErrorEps) {
			t.Errorf("Expected Mz = %f, but got %f", want, got)
		}
	})
}
