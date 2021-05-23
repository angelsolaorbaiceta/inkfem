package tests

import (
	"github.com/angelsolaorbaiceta/inkfem/contracts"
	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/process"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"math"
	"testing"
)

var (
	noDistLoads = []*load.DistributedLoad{}
	noConcLoads = []*load.ConcentratedLoad{}
	material    = &structure.Material{
		Name:             "steel",
		Density:          0,
		YoungMod:         20e6,
		ShearMod:         0,
		PoissonRatio:     1,
		YieldStrength:    0,
		UltimateStrength: 0,
	}
	section = &structure.Section{
		Name:    "IPE 120",
		Area:    14,
		IStrong: 318,
		IWeak:   28,
		SStrong: 53,
		SWeak:   9,
	}
	length     = 100.0
	displError = 1e-5
)

func TestCantileverBeamWithConcentratedVerticalLoadAtEnd(t *testing.T) {
	var (
		l               = load.MakeConcentrated(load.FY, true, inkgeom.MaxT, -2000)
		str             = makeBeamStructure([]*load.ConcentratedLoad{l}, noDistLoads)
		sol             = solveStructure(str)
		solutionElement = sol.Elements[0]
		maxYDispl       = -200.0 / 1908.0 // PL³ / 3EI
		maxZRot         = -1.0 / 636.0    // -PL² / 2EI
	)

	t.Run("global X displacements", func(t *testing.T) {
		for _, disp := range solutionElement.GlobalXDispl {
			if !inkgeom.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("local X displacements", func(t *testing.T) {
		for _, disp := range solutionElement.LocalXDispl {
			if !inkgeom.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("global Y displacements", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalYDispl)

		if got := solutionElement.GlobalYDispl[0].Value; !inkgeom.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Y displacement in the constrained end")
		}
		if got := solutionElement.GlobalYDispl[nOfValues-1].Value; !inkgeom.FloatsEqualEps(got, maxYDispl, displError) {
			t.Errorf("expected max Y displacement of %f, but got %f", maxYDispl, got)
		}
	})

	t.Run("local Y displacements", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalYDispl)

		if got := solutionElement.LocalYDispl[0].Value; !inkgeom.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Y displacement in the constrained end")
		}
		if got := solutionElement.LocalYDispl[nOfValues-1].Value; !inkgeom.FloatsEqualEps(got, maxYDispl, displError) {
			t.Errorf("expected max Y displacement of %f, but got %f", maxYDispl, got)
		}
	})

	t.Run("global Z rotations", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalZRot)

		if got := solutionElement.GlobalZRot[0].Value; !inkgeom.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Z rotation in the constrained end")
		}
		if got := solutionElement.GlobalZRot[nOfValues-1].Value; !inkgeom.FloatsEqualEps(got, maxZRot, displError) {
			t.Errorf("expected max Z rotation of %f, but got %f", maxZRot, got)
		}
	})

	t.Run("local Z rotations", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalZRot)

		if got := solutionElement.GlobalZRot[0].Value; !inkgeom.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Z rotation in the constrained end")
		}
		if got := solutionElement.GlobalZRot[nOfValues-1].Value; !inkgeom.FloatsEqualEps(got, maxZRot, displError) {
			t.Errorf("expected max Z rotation of %f, but got %f", maxZRot, got)
		}
	})

	t.Run("Axial stress", func(t *testing.T) {
		for _, axial := range solutionElement.AxialStress {
			if !inkgeom.FloatsEqualEps(axial.Value, 0.0, displError) {
				t.Errorf("Expected no axial stress, but got %f", axial.Value)
			}
		}
	})

	t.Run("Shear stress", func(t *testing.T) {
		expectedShear := -l.Value

		for _, shear := range solutionElement.ShearStress {
			if !inkgeom.FloatsEqualEps(shear.Value, expectedShear, displError) {
				t.Errorf("Expected a shear stress of %f, but got %f at t = %f", expectedShear, shear.Value, shear.T)
			}
		}
	})

	t.Run("Bending moment", func(t *testing.T) {
		var expectedBending = func(tParam inkgeom.TParam) float64 {
			return l.Value * length * (inkgeom.MaxT.Value() - tParam.Value())
		}

		for _, bending := range solutionElement.BendingMoment {
			var (
				got  = bending.Value
				want = expectedBending(bending.T)
			)

			if !inkgeom.FloatsEqualEps(got, want, displError) {
				t.Errorf("Expected a bending moment of %f, but got %f at t = %f", want, got, bending.T)
			}
		}
	})
}

func TestCantileverBeamWithDistributedVerticalLoad(t *testing.T) {
	var (
		l               = load.MakeDistributed(load.FY, true, inkgeom.MinT, -200.0, inkgeom.MaxT, 0.0)
		str             = makeBeamStructure(noConcLoads, []*load.DistributedLoad{l})
		sol             = solveStructure(str)
		solutionElement = sol.Elements[0]
		maxYDispl       = -200.0 / 1908.0 // WL⁴ / 30EI
		maxZRot         = -20.0 / 15264.0 // WL³ / 24EI
	)

	t.Run("global X displacements", func(t *testing.T) {
		for _, disp := range solutionElement.GlobalXDispl {
			if !inkgeom.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("local X displacements", func(t *testing.T) {
		for _, disp := range solutionElement.LocalXDispl {
			if !inkgeom.FloatsEqualEps(disp.Value, 0.0, displError) {
				t.Errorf("Expected no X displacement, but got %f", disp.Value)
			}
		}
	})

	t.Run("global Y displacements", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalYDispl)

		if got := solutionElement.GlobalYDispl[0].Value; !inkgeom.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Y displacement in the constrained end")
		}
		if got := solutionElement.GlobalYDispl[nOfValues-1].Value; !inkgeom.FloatsEqualEps(got, maxYDispl, displError) {
			t.Errorf("expected max Y displacement of %f, but got %f", maxYDispl, got)
		}
	})

	t.Run("local Y displacements", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalYDispl)

		if got := solutionElement.LocalYDispl[0].Value; !inkgeom.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Y displacement in the constrained end")
		}
		if got := solutionElement.LocalYDispl[nOfValues-1].Value; !inkgeom.FloatsEqualEps(got, maxYDispl, displError) {
			t.Errorf("expected max Y displacement of %f, but got %f", maxYDispl, got)
		}
	})

	t.Run("global Z rotations", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalZRot)

		if got := solutionElement.GlobalZRot[0].Value; !inkgeom.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Z rotation in the constrained end")
		}
		if got := solutionElement.GlobalZRot[nOfValues-1].Value; !inkgeom.FloatsEqualEps(got, maxZRot, displError) {
			t.Errorf("expected max Z rotation of %f, but got %f", maxZRot, got)
		}
	})

	t.Run("local Z rotations", func(t *testing.T) {
		nOfValues := len(solutionElement.GlobalZRot)

		if got := solutionElement.GlobalZRot[0].Value; !inkgeom.FloatsEqualEps(got, 0.0, displError) {
			t.Error("expected no Z rotation in the constrained end")
		}
		if got := solutionElement.GlobalZRot[nOfValues-1].Value; !inkgeom.FloatsEqualEps(got, maxZRot, displError) {
			t.Errorf("expected max Z rotation of %f, but got %f", maxZRot, got)
		}
	})

	t.Run("Axial stress", func(t *testing.T) {
		for _, axial := range solutionElement.AxialStress {
			if !inkgeom.FloatsEqualEps(axial.Value, 0.0, displError) {
				t.Errorf("Expected no axial stress, but got %f", axial.Value)
			}
		}
	})

	t.Run("Shear stress", func(t *testing.T) {
		var expectedShear = func(tParam inkgeom.TParam) float64 {
			var (
				qStart = l.ValueAt(inkgeom.MinT)
				x      = length * tParam.Value()
			)

			return qStart * (-0.5*length + x - 0.5*x*x/length)
		}

		for _, shear := range solutionElement.ShearStress {
			var (
				got  = shear.Value
				want = expectedShear(shear.T)
			)

			if !inkgeom.FloatsEqualEps(got, want, displError) {
				t.Errorf("Expected a shear stress of %f, but got %f at t = %f", want, got, shear.T)
			}
		}
	})

	t.Run("Bending moment", func(t *testing.T) {
		var expectedBending = func(tParam inkgeom.TParam) float64 {
			var (
				qStart = l.ValueAt(inkgeom.MinT)
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

			if !inkgeom.FloatsEqualEps(got, want, displError) {
				t.Errorf("Expected a bending moment of %f, but got %f at t = %f", want, got, bending.T)
			}
		}
	})
}

func makeBeamStructure(
	concentratedLoads []*load.ConcentratedLoad,
	distributedLoads []*load.DistributedLoad,
) *structure.Structure {
	var (
		nodeOne = structure.MakeNode("fixed-node", g2d.MakePoint(0, 0), structure.FullConstraint)
		nodeTwo = structure.MakeNode("free-node", g2d.MakePoint(length, 0), structure.NilConstraint)
		beam    = structure.MakeElement(
			"beam",
			nodeOne,
			nodeTwo,
			structure.FullConstraint,
			structure.FullConstraint,
			material,
			section,
			concentratedLoads,
			distributedLoads,
		)
	)

	return &structure.Structure{
		structure.StrMetadata{1, 0},
		map[contracts.StrID]*structure.Node{
			nodeOne.Id: nodeOne,
			nodeTwo.Id: nodeTwo,
		},
		[]*structure.Element{beam},
	}
}

func solveStructure(str *structure.Structure) *process.Solution {
	solveOptions := process.SolveOptions{false, "", true, displError}
	pre := preprocess.DoStructure(str)
	return process.Solve(pre, solveOptions)
}
