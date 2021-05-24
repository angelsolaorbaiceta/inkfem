package tests

import (
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"math"
	"testing"
)

func TestAxialMemberWithConcentratedLoad(t *testing.T) {
	var (
		l               = load.MakeConcentrated(load.FX, true, inkgeom.MaxT, 4000.0)
		str             = makeAxialElementStructure([]*load.ConcentratedLoad{l}, noDistLoads)
		sol             = solveStructure(str)
		solutionElement = sol.Elements[0]
	)

	t.Run("X displacements", func(t *testing.T) {
		var expectedXDispl = func(tParam inkgeom.TParam) float64 {
			var (
				ea = material.YoungMod * section.Area
				x  = tParam.Value() * length
				f  = l.Value
			)

			return f * x / ea
		}

		for _, disp := range solutionElement.LocalXDispl {
			if want := expectedXDispl(disp.T); !inkgeom.FloatsEqualEps(disp.Value, want, displError) {
				t.Errorf("Expected X displacement of %f, but got %f at %f", want, disp.Value, disp.T)
			}
		}
	})

	t.Run("Y displacements", func(t *testing.T) {
		for _, disp := range solutionElement.LocalYDispl {
			if !inkgeom.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("Z rotations", func(t *testing.T) {
		for _, disp := range solutionElement.LocalZRot {
			if !inkgeom.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})
}

func TestAxialMemberWithDistributedLoad(t *testing.T) {
	var (
		l               = load.MakeDistributed(load.FX, true, inkgeom.MinT, 400.0, inkgeom.MaxT, 0.0)
		str             = makeAxialElementStructure(noConcLoads, []*load.DistributedLoad{l})
		sol             = solveStructure(str)
		solutionElement = sol.Elements[0]
	)

	t.Run("X displacements", func(t *testing.T) {
		var expectedXDispl = func(tParam inkgeom.TParam) float64 {
			var (
				ea     = material.YoungMod * section.Area
				x      = tParam.Value() * length
				load_a = l.ValueAt(inkgeom.MinT)
				load_b = -load_a / length
			)

			return (load_a*x + 0.5*load_b*math.Pow(x, 2)) / ea
		}

		for _, disp := range solutionElement.LocalXDispl {
			if want := expectedXDispl(disp.T); !inkgeom.FloatsEqualEps(disp.Value, want, displError) {
				t.Errorf("Expected X displacement of %f, but got %f at %f", want, disp.Value, disp.T)
			}
		}
	})

	t.Run("Y displacements", func(t *testing.T) {
		for _, disp := range solutionElement.LocalYDispl {
			if !inkgeom.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("Z rotations", func(t *testing.T) {
		for _, disp := range solutionElement.LocalZRot {
			if !inkgeom.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})
}
