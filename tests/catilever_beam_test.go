package tests

import (
	"math"
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/build"
	strmath "github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

func TestCantileverBeamWithConcentratedVerticalLoadAtEnd(t *testing.T) {
	build.Info = &build.BuildInfo{MajorVersion: 3, MinorVersion: 2}

	var (
		l               = load.MakeConcentrated(load.FY, true, nums.MaxT, -2000)
		str             = makeCantileverBeamStructure([]*load.ConcentratedLoad{l}, noDistLoads)
		sol             = solveStructure(str)
		solutionElement = sol.Elements[0]
		maxYDispl       = -200.0 / 1908.0 // PL³ / 3EI
		maxZRot         = -1.0 / 636.0    // -PL² / 2EI
	)

	t.Run("global X displacements", func(t *testing.T) {
		for _, disp := range solutionElement.GlobalXDispl {
			if !nums.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("local X displacements", func(t *testing.T) {
		for _, disp := range solutionElement.LocalXDispl {
			if !nums.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("global Y displacements", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalYDispl)

		if got := solutionElement.GlobalYDispl[0].Value; !nums.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Y displacement in the constrained end")
		}
		if got := solutionElement.GlobalYDispl[nOfValues-1].Value; !nums.FloatsEqualEps(got, maxYDispl, displError) {
			t.Errorf("expected max Y displacement of %f, but got %f", maxYDispl, got)
		}
	})

	t.Run("local Y displacements", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalYDispl)

		if got := solutionElement.LocalYDispl[0].Value; !nums.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Y displacement in the constrained end")
		}
		if got := solutionElement.LocalYDispl[nOfValues-1].Value; !nums.FloatsEqualEps(got, maxYDispl, displError) {
			t.Errorf("expected max Y displacement of %f, but got %f", maxYDispl, got)
		}
	})

	t.Run("global Z rotations", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalZRot)

		if got := solutionElement.GlobalZRot[0].Value; !nums.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Z rotation in the constrained end")
		}
		if got := solutionElement.GlobalZRot[nOfValues-1].Value; !nums.FloatsEqualEps(got, maxZRot, displError) {
			t.Errorf("expected max Z rotation of %f, but got %f", maxZRot, got)
		}
	})

	t.Run("local Z rotations", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalZRot)

		if got := solutionElement.GlobalZRot[0].Value; !nums.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Z rotation in the constrained end")
		}
		if got := solutionElement.GlobalZRot[nOfValues-1].Value; !nums.FloatsEqualEps(got, maxZRot, displError) {
			t.Errorf("expected max Z rotation of %f, but got %f", maxZRot, got)
		}
	})

	t.Run("Axial stress", func(t *testing.T) {
		for _, axial := range solutionElement.AxialStress {
			if !nums.FloatsEqualEps(axial.Value, 0.0, displError) {
				t.Errorf("Expected no axial stress, but got %f", axial.Value)
			}
		}
	})

	t.Run("Shear force", func(t *testing.T) {
		expectedShear := -l.Value

		for _, shear := range solutionElement.ShearForce {
			if !nums.FloatsEqualEps(shear.Value, expectedShear, displError) {
				t.Errorf("Expected a Shear force of %f, but got %f at t = %f", expectedShear, shear.Value, shear.T)
			}
		}
	})

	t.Run("Bending moment", func(t *testing.T) {
		var expectedBending = func(tParam nums.TParam) float64 {
			return l.Value * length * (nums.MaxT.Value() - tParam.Value())
		}

		for _, bending := range solutionElement.BendingMoment {
			var (
				got  = bending.Value
				want = expectedBending(bending.T)
			)

			if !nums.FloatsEqualEps(got, want, displError) {
				t.Errorf("Expected a bending moment of %f, but got %f at t = %f", want, got, bending.T)
			}
		}
	})

	t.Run("Reaction Torsor", func(t *testing.T) {
		var (
			want = strmath.MakeTorsor(0.0, -l.Value, -l.Value*length)
			got  = sol.NodeReactions()["fixed-node"]
		)

		if !nums.FloatsEqualEps(got.Fx(), want.Fx(), displError) {
			t.Errorf("Expected Fx reaction %f, but got %f", want.Fx(), got.Fx())
		}
		if !nums.FloatsEqualEps(got.Fy(), want.Fy(), displError) {
			t.Errorf("Expected Fy reaction %f, but got %f", want.Fy(), got.Fy())
		}
		if !nums.FloatsEqualEps(got.Mz(), want.Mz(), displError) {
			t.Errorf("Expected Mz reaction %f, but got %f", want.Mz(), got.Mz())
		}
	})
}

func TestCantileverBeamWithDistributedVerticalLoad(t *testing.T) {
	build.Info = &build.BuildInfo{MajorVersion: 3, MinorVersion: 2}

	var (
		l               = load.MakeDistributed(load.FY, true, nums.MinT, -200.0, nums.MaxT, 0.0)
		str             = makeCantileverBeamStructure(noConcLoads, []*load.DistributedLoad{l})
		sol             = solveStructure(str)
		solutionElement = sol.Elements[0]
		maxYDispl       = -200.0 / 1908.0 // WL⁴ / 30EI
		maxZRot         = -20.0 / 15264.0 // WL³ / 24EI
	)

	t.Run("global X displacements", func(t *testing.T) {
		for _, disp := range solutionElement.GlobalXDispl {
			if !nums.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("local X displacements", func(t *testing.T) {
		for _, disp := range solutionElement.LocalXDispl {
			if !nums.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("global Y displacements", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalYDispl)

		if got := solutionElement.GlobalYDispl[0].Value; !nums.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Y displacement in the constrained end")
		}
		if got := solutionElement.GlobalYDispl[nOfValues-1].Value; !nums.FloatsEqualEps(got, maxYDispl, displError) {
			t.Errorf("expected max Y displacement of %f, but got %f", maxYDispl, got)
		}
	})

	t.Run("local Y displacements", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalYDispl)

		if got := solutionElement.LocalYDispl[0].Value; !nums.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Y displacement in the constrained end")
		}
		if got := solutionElement.LocalYDispl[nOfValues-1].Value; !nums.FloatsEqualEps(got, maxYDispl, displError) {
			t.Errorf("expected max Y displacement of %f, but got %f", maxYDispl, got)
		}
	})

	t.Run("global Z rotations", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalZRot)

		if got := solutionElement.GlobalZRot[0].Value; !nums.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Z rotation in the constrained end")
		}
		if got := solutionElement.GlobalZRot[nOfValues-1].Value; !nums.FloatsEqualEps(got, maxZRot, displError) {
			t.Errorf("expected max Z rotation of %f, but got %f", maxZRot, got)
		}
	})

	t.Run("local Z rotations", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalZRot)

		if got := solutionElement.GlobalZRot[0].Value; !nums.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Z rotation in the constrained end")
		}
		if got := solutionElement.GlobalZRot[nOfValues-1].Value; !nums.FloatsEqualEps(got, maxZRot, displError) {
			t.Errorf("expected max Z rotation of %f, but got %f", maxZRot, got)
		}
	})

	t.Run("Axial stress", func(t *testing.T) {
		for _, axial := range solutionElement.AxialStress {
			if !nums.FloatsEqualEps(axial.Value, 0.0, displError) {
				t.Errorf("Expected no axial stress, but got %f", axial.Value)
			}
		}
	})

	t.Run("Shear force", func(t *testing.T) {
		var expectedShear = func(tParam nums.TParam) float64 {
			var (
				qStart = l.ValueAt(nums.MinT)
				x      = length * tParam.Value()
			)

			return qStart * (-0.5*length + x - 0.5*x*x/length)
		}

		for _, shear := range solutionElement.ShearForce {
			var (
				got  = shear.Value
				want = expectedShear(shear.T)
			)

			if !nums.FloatsEqualEps(got, want, displError) {
				t.Errorf("Expected a Shear force of %f, but got %f at t = %f", want, got, shear.T)
			}
		}
	})

	t.Run("Bending moment", func(t *testing.T) {
		var expectedBending = func(tParam nums.TParam) float64 {
			var (
				qStart = l.ValueAt(nums.MinT)
				mStart = qStart * math.Pow(length, 2) / 6.0
				x      = length * tParam.Value()
			)

			return mStart + 0.5*qStart*(-length*x+math.Pow(x, 2)-math.Pow(x, 3)/(3.0*length))
		}

		for _, bending := range solutionElement.BendingMoment {
			var (
				got  = bending.Value
				want = expectedBending(bending.T)
			)

			if !nums.FloatsEqualEps(got, want, displError) {
				t.Errorf("Expected a bending moment of %f, but got %f at t = %f", want, got, bending.T)
			}
		}
	})
}
