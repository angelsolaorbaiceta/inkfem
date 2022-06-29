package process

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/preprocess"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/angelsolaorbaiceta/inkmath/vec"
)

// This test uses an horizontal cantilever bar with nodes at (0, 0) and (400, 0)
// with three preprocess nodes, each of them subject to horizontal forces.
// The left preprocessed node doesn't move, but the two others move 5 and 10 units, horizontally.
//
// This yields a uniform axial stress on the bar: σ = E * 𝛅 / L = 100 * 5 / 200 = 2.5
func TestAxialBarSolution(t *testing.T) {
	var (
		bar    = makeElementSolutionTestOriginalBar()
		preBar = preprocess.MakeElement(
			bar,
			[]*preprocess.Node{
				preprocess.MakeNodeWithDofs(nums.MinT, g2d.MakePoint(0, 0), [3]int{0, 1, 2}),
				preprocess.MakeNodeWithDofs(nums.HalfT, g2d.MakePoint(200, 0), [3]int{3, 4, 5}),
				preprocess.MakeNodeWithDofs(nums.MaxT, g2d.MakePoint(400, 0), [3]int{6, 7, 8}),
			},
		)
		barSolution = MakeElementSolution(
			preBar,
			vec.MakeWithValues([]float64{0, 0, 0, 5, 0, 0, 10, 0, 0}),
		)
	)

	t.Run("The axial stress", func(t *testing.T) {
		var (
			axial          = barSolution.AxialStress
			wantAxialLeft  = PointSolutionValue{nums.MinT, 2.5}
			wantAxialMid   = PointSolutionValue{nums.HalfT, 2.5}
			wantAxialRight = PointSolutionValue{nums.MaxT, 2.5}
		)

		// Since the middle node should have the same axial stress both at its left and right,
		// there should be only three axial stress values: the left, the middle and the right.
		if len(axial) != 3 {
			t.Errorf("Expected 3 axial stress values, got %d", len(axial))
		}

		if !axial[0].Equals(wantAxialLeft) {
			t.Errorf("Expected axial stress %f, got %f", wantAxialLeft, axial[0])
		}
		if !axial[1].Equals(wantAxialMid) {
			t.Errorf("Expected axial stress %f, got %f", wantAxialMid, axial[1])
		}
		if !axial[2].Equals(wantAxialRight) {
			t.Errorf("Expected axial stress %f, got %f", wantAxialRight, axial[2])
		}
	})
}

// This test uses an horizontal cantilever bar with nodes at (0, 0) and (400, 0)
// with three preprocess nodes, each of them subject to vertical forces.
// The left preprocessed node doesn't move, but the two others move 5 and 10 units, vertically.
// The rotations are chosen to that the shear force due to rotation is uniform throughout the bar.
//
// This yields a uniform shear force on the bar generated by the displacements:
// (trail_dy - lead_dy) * 12 * E * I / L^3
// = (0 + 5) * 12 * 100 * 8 / 200^3
// = 0.006
//
// And a uniform shear force on the bar generated by the rotations:
// (trailRz + leadRz) * 6.0 * E * I / L^2
// = (0 - 0.1) * 6.0 * 100 * 8 / 200^2
// = -0.012
//
// The total shear force is therefore:
// 0.006 + (-0.012) = -0.006
func TestShearBarSolution(t *testing.T) {
	var (
		bar    = makeElementSolutionTestOriginalBar()
		preBar = preprocess.MakeElement(
			bar,
			[]*preprocess.Node{
				preprocess.MakeNodeWithDofs(nums.MinT, g2d.MakePoint(0, 0), [3]int{0, 1, 2}),
				preprocess.MakeNodeWithDofs(nums.HalfT, g2d.MakePoint(200, 0), [3]int{3, 4, 5}),
				preprocess.MakeNodeWithDofs(nums.MaxT, g2d.MakePoint(400, 0), [3]int{6, 7, 8}),
			},
		)
		barSolution = MakeElementSolution(
			preBar,
			vec.MakeWithValues([]float64{0, 0, 0, 0, -5, -0.1, 0, -10, 0.0}),
		)
	)

	t.Run("The shear force", func(t *testing.T) {
		var (
			shear          = barSolution.ShearForce
			wantShearLeft  = PointSolutionValue{nums.MinT, -0.006}
			wantShearMid   = PointSolutionValue{nums.HalfT, -0.006}
			wantShearRight = PointSolutionValue{nums.MaxT, -0.006}
		)

		// Since the middle node should have the same shear force both at its left and right,
		// there should be only three shear force values: the left, the middle and the right.
		if len(shear) != 3 {
			t.Errorf("Expected 3 shear force values, got %d", len(shear))
		}

		if !shear[0].Equals(wantShearLeft) {
			t.Errorf("Expected shear force %f, got %f", wantShearLeft, shear[0])
		}
		if !shear[1].Equals(wantShearMid) {
			t.Errorf("Expected shear force %f, got %f", wantShearMid, shear[1])
		}
		if !shear[2].Equals(wantShearRight) {
			t.Errorf("Expected shear force %f, got %f", wantShearRight, shear[2])
		}
	})
}

// Makes a bar with nodes at (0, 0) and (400, 0), where the left node is fixed (cantilever beam).
// The section has a cross section of 5 and both moments of inertia are 8.
// The material has a Young modulus of 100.
func makeElementSolutionTestOriginalBar() *structure.Element {
	var (
		nodeOne  = structure.MakeNode("n1", g2d.MakePoint(0, 0), &structure.FullConstraint)
		nodeTwo  = structure.MakeNode("n2", g2d.MakePoint(400, 0), &structure.NilConstraint)
		section  = structure.MakeSection("section", 5.0, 8.0, 8.0, 0.0, 0.0)
		material = structure.MakeMaterial("material", 0.0, 100.0, 0.0, 0.0, 0.0, 0.0)
		bar      = structure.MakeElementBuilder("b1").
				WithStartNode(nodeOne, &structure.FullConstraint).
				WithEndNode(nodeTwo, &structure.FullConstraint).
				WithSection(section).
				WithMaterial(material).
				Build()
	)

	return bar
}
