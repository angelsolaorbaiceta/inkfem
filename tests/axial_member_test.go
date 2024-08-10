package tests

import (
	"math"
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/build"
	strmath "github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

func TestAxialMemberWithConcentratedLoad(t *testing.T) {
	build.Info = &build.BuildInfo{MajorVersion: 3, MinorVersion: 2}

	var (
		l               = load.MakeConcentrated(load.FX, true, nums.MaxT, 4000.0)
		str             = makeAxialElementStructure([]*load.ConcentratedLoad{l}, noDistLoads)
		sol             = solveStructure(str)
		solutionElement = sol.Elements[0]
	)

	t.Run("X displacements", func(t *testing.T) {
		var (
			ea = material.YoungMod * section.Area
			f  = l.Value
		)

		var expectedXDispl = func(tParam nums.TParam) float64 {
			x := tParam.Value() * length
			return f * x / ea
		}

		for _, disp := range solutionElement.LocalXDispl {
			if want := expectedXDispl(disp.T); !nums.FloatsEqualEps(disp.Value, want, displError) {
				t.Errorf("Expected X displacement of %f, but got %f at %f", want, disp.Value, disp.T)
			}
		}
	})

	t.Run("Y displacements", func(t *testing.T) {
		for _, disp := range solutionElement.LocalYDispl {
			if !nums.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("Z rotations", func(t *testing.T) {
		for _, disp := range solutionElement.LocalZRot {
			if !nums.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("Axial stress", func(t *testing.T) {
		expectedAxial := l.Value / section.Area

		for _, axial := range solutionElement.AxialStress {
			if !nums.FloatsEqualEps(axial.Value, expectedAxial, displError) {
				t.Errorf("Expected axial stress of %f, but got %f at t = %f", expectedAxial, axial.Value, axial.T)
			}
		}
	})

	t.Run("Shear force", func(t *testing.T) {
		for _, shear := range solutionElement.ShearForce {
			if !nums.FloatsEqualEps(shear.Value, 0.0, displError) {
				t.Errorf("Expected no Shear force but got %f at t = %f", shear.Value, shear.T)
			}
		}
	})

	t.Run("Bending moment", func(t *testing.T) {
		for _, bending := range solutionElement.BendingMoment {
			if !nums.FloatsEqualEps(bending.Value, 0.0, displError) {
				t.Errorf("Expected no bending moment but got %f at t = %f", bending.Value, bending.T)
			}
		}
	})

	t.Run("Reaction Torsor", func(t *testing.T) {
		want := strmath.MakeTorsor(-l.Value, 0.0, 0.0)

		if got := sol.NodeReactions()["fixed-node"]; !got.Equals(want) {
			t.Errorf("Expected reaction torsor %v, but got %v", want, got)
		}
	})
}

func TestAxialMemberWithConstantDistributedLoad(t *testing.T) {
	build.Info = &build.BuildInfo{MajorVersion: 3, MinorVersion: 2}

	var (
		l               = load.MakeDistributed(load.FX, true, nums.MinT, 400.0, nums.MaxT, 400.0)
		str             = makeAxialElementStructure(noConcLoads, []*load.DistributedLoad{l})
		sol             = solveStructure(str)
		solutionElement = sol.Elements[0]
	)

	t.Run("Reaction", func(t *testing.T) {
		var (
			got  = sol.NodeReactions()["fixed-node"]
			want = strmath.MakeTorsor(-l.ValueAt(nums.MinT)*length, 0.0, 0.0)
		)

		if !got.Equals(want) {
			t.Errorf("Expected reaction %v, but got %v", want, got)
		}
	})

	t.Run("X displacements", func(t *testing.T) {
		var (
			ea                   = material.YoungMod * section.Area
			load_a               = l.ValueAt(nums.MinT)
			normalReactionForce  = l.ValueAt(nums.MinT) * length
			normalReactionStress = normalReactionForce / section.Area
		)

		var expectedXDispl = func(tParam nums.TParam) float64 {
			var (
				x  = tParam.Value() * length
				x2 = math.Pow(x, 2)
			)

			return (normalReactionStress * x / material.YoungMod) - (load_a*x2/2.0)/ea
		}

		for _, disp := range solutionElement.LocalXDispl {
			if want := expectedXDispl(disp.T); !nums.FloatsEqualEps(disp.Value, want, displError) {
				t.Errorf("Expected X displacement of %f, but got %f at %f", want, disp.Value, disp.T)
			}
		}
	})

	t.Run("Y displacements", func(t *testing.T) {
		for _, disp := range solutionElement.LocalYDispl {
			if !nums.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("Z rotations", func(t *testing.T) {
		for _, disp := range solutionElement.LocalZRot {
			if !nums.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("Axial stress", func(t *testing.T) {
		var (
			normalReactionForce  = l.ValueAt(nums.MinT) * length
			normalReactionStress = normalReactionForce / section.Area
			loadStress           = l.ValueAt(nums.MinT) / section.Area
		)

		var expectedAxial = func(tParam nums.TParam) float64 {
			x := tParam.Value() * length
			return normalReactionStress - loadStress*x
		}

		for _, axial := range solutionElement.AxialStress {
			if want := expectedAxial(axial.T); !nums.FloatsEqualEps(axial.Value, want, displError) {
				t.Errorf("Expected axial stress of %f, but got %f at t = %f", want, axial.Value, axial.T)
			}
		}
	})

	t.Run("Shear force", func(t *testing.T) {
		for _, shear := range solutionElement.ShearForce {
			if !nums.FloatsEqualEps(shear.Value, 0.0, displError) {
				t.Errorf("Expected no Shear force but got %f at t = %f", shear.Value, shear.T)
			}
		}
	})

	t.Run("Bending moment", func(t *testing.T) {
		for _, bending := range solutionElement.BendingMoment {
			if !nums.FloatsEqualEps(bending.Value, 0.0, displError) {
				t.Errorf("Expected no bending moment but got %f at t = %f", bending.Value, bending.T)
			}
		}
	})
}

func TestAxialMemberWithDistributedLoad(t *testing.T) {
	var (
		l                    = load.MakeDistributed(load.FX, true, nums.MinT, 400.0, nums.MaxT, 0.0)
		str                  = makeAxialElementStructure(noConcLoads, []*load.DistributedLoad{l})
		sol                  = solveStructure(str)
		solutionElement      = sol.Elements[0]
		ea                   = material.YoungMod * section.Area
		normalReactionForce  = 0.5 * l.ValueAt(nums.MinT) * length
		normalReactionStress = normalReactionForce / section.Area
		load_a               = l.ValueAt(nums.MinT)
		load_b               = -load_a / length
	)

	t.Run("X displacements", func(t *testing.T) {
		var expectedXDispl = func(tParam nums.TParam) float64 {
			var (
				x  = tParam.Value() * length
				x2 = math.Pow(x, 2)
				x3 = math.Pow(x, 3)
			)

			return normalReactionStress*x/material.YoungMod - (load_a*x2/2.0+load_b*x3/6.0)/ea
		}

		for _, disp := range solutionElement.LocalXDispl {
			if want := expectedXDispl(disp.T); !nums.FloatsEqualEps(disp.Value, want, displError) {
				t.Errorf("Expected X displacement of %f, but got %f at %f", want, disp.Value, disp.T)
			}
		}
	})

	t.Run("Y displacements", func(t *testing.T) {
		for _, disp := range solutionElement.LocalYDispl {
			if !nums.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("Z rotations", func(t *testing.T) {
		for _, disp := range solutionElement.LocalZRot {
			if !nums.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("Axial stress", func(t *testing.T) {
		var expectedAxial = func(tParam nums.TParam) float64 {
			var (
				x  = tParam.Value() * length
				x2 = math.Pow(x, 2)
			)

			return normalReactionStress - (load_a*x+load_b*x2/2.0)/section.Area
		}

		for _, axial := range solutionElement.AxialStress {
			if want := expectedAxial(axial.T); !nums.FloatsEqualEps(axial.Value, want, displError) {
				t.Errorf("Expected axial stress of %f, but got %f at t = %f", want, axial.Value, axial.T)
			}
		}
	})

	t.Run("Shear force", func(t *testing.T) {
		for _, shear := range solutionElement.ShearForce {
			if !nums.FloatsEqualEps(shear.Value, 0.0, displError) {
				t.Errorf("Expected no Shear force but got %f at t = %f", shear.Value, shear.T)
			}
		}
	})

	t.Run("Bending moment", func(t *testing.T) {
		for _, bending := range solutionElement.BendingMoment {
			if !nums.FloatsEqualEps(bending.Value, 0.0, displError) {
				t.Errorf("Expected no bending moment but got %f at t = %f", bending.Value, bending.T)
			}
		}
	})
}
