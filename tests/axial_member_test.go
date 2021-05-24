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

	t.Run("Axial stress", func(t *testing.T) {
		expectedAxial := l.Value / section.Area
		
		for _, axial := range solutionElement.AxialStress {
			if !inkgeom.FloatsEqualEps(axial.Value, expectedAxial, displError) {
				t.Errorf("Expected axial stress of %f, but got %f at t = %f", expectedAxial, axial.Value, axial.T)
			}
		}
	})

	t.Run("Shear stress", func(t *testing.T) {
		for _, shear := range solutionElement.ShearStress {
			if !inkgeom.FloatsEqualEps(shear.Value, 0.0, displError) {
				t.Errorf("Expected no shear stress but got %f at t = %f", shear.Value, shear.T)
			}
		}
	})

	t.Run("Bending moment", func(t *testing.T) {
		for _, bending := range solutionElement.BendingMoment {
			if !inkgeom.FloatsEqualEps(bending.Value, 0.0, displError) {
				t.Errorf("Expected no bending moment but got %f at t = %f", bending.Value, bending.T)
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
				x2     = math.Pow(x, 2)
				x3     = math.Pow(x, 3)
				load_a = l.ValueAt(inkgeom.MinT)
				load_b = -load_a / length
			)

			return (load_a*x2/2.0 + load_b*x3/3.0) / ea
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
