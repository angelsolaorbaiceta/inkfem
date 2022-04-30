package generate

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/build"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

func TestGenerateReticularStructure(t *testing.T) {
	build.Info = &build.BuildInfo{
		MajorVersion: 3,
		MinorVersion: 2,
	}

	var (
		params = ReticStructureParams{
			Spans:         1,
			Levels:        2,
			Span:          300.0,
			Height:        200.0,
			LoadDistValue: 50.0,
			Section:       structure.MakeUnitSection(),
			Material:      structure.MakeUnitMaterial(),
		}
		str = Reticular(params)

		n1 = structure.MakeNodeAtPosition("1", 0, 0, &structure.FullConstraint)
		n2 = structure.MakeNodeAtPosition("2", 300, 0, &structure.FullConstraint)
		n3 = structure.MakeNodeAtPosition("3", 0, 200, &structure.NilConstraint)
		n4 = structure.MakeNodeAtPosition("4", 300, 200, &structure.NilConstraint)
		n5 = structure.MakeNodeAtPosition("5", 0, 400, &structure.NilConstraint)
		n6 = structure.MakeNodeAtPosition("6", 300, 400, &structure.NilConstraint)

		e1 = structure.MakeElementBuilder("1").
			WithStartNode(n1, &structure.FullConstraint).
			WithEndNode(n3, &structure.FullConstraint).
			WithMaterial(structure.MakeUnitMaterial()).
			WithSection(structure.MakeUnitSection()).
			Build()
		e2 = structure.MakeElementBuilder("2").
			WithStartNode(n2, &structure.FullConstraint).
			WithEndNode(n4, &structure.FullConstraint).
			WithMaterial(structure.MakeUnitMaterial()).
			WithSection(structure.MakeUnitSection()).
			Build()
		e3 = structure.MakeElementBuilder("3").
			WithStartNode(n3, &structure.FullConstraint).
			WithEndNode(n4, &structure.FullConstraint).
			WithMaterial(structure.MakeUnitMaterial()).
			WithSection(structure.MakeUnitSection()).
			Build()
		e4 = structure.MakeElementBuilder("4").
			WithStartNode(n3, &structure.FullConstraint).
			WithEndNode(n5, &structure.FullConstraint).
			WithMaterial(structure.MakeUnitMaterial()).
			WithSection(structure.MakeUnitSection()).
			Build()
		e5 = structure.MakeElementBuilder("5").
			WithStartNode(n4, &structure.FullConstraint).
			WithEndNode(n6, &structure.FullConstraint).
			WithMaterial(structure.MakeUnitMaterial()).
			WithSection(structure.MakeUnitSection()).
			Build()
		e6 = structure.MakeElementBuilder("6").
			WithStartNode(n5, &structure.FullConstraint).
			WithEndNode(n6, &structure.FullConstraint).
			WithMaterial(structure.MakeUnitMaterial()).
			WithSection(structure.MakeUnitSection()).
			Build()

		load = load.MakeDistributed(
			load.FY,
			true,
			nums.MinT, -params.LoadDistValue,
			nums.MaxT, -params.LoadDistValue,
		)
	)

	t.Run("uses the binary's version", func(t *testing.T) {
		if got := str.Metadata.MajorVersion; got != 3 {
			t.Errorf("got %v, want %v", got, 3)
		}
		if got := str.Metadata.MinorVersion; got != 2 {
			t.Errorf("got %v, want %v", got, 2)
		}
	})

	t.Run("generates the nodes", func(t *testing.T) {
		if got := str.NodesCount(); got != 6 {
			t.Errorf("Want %d nodes, got %d", 6, got)
		}

		if got := str.GetNodeById("1"); !n1.Equals(got) {
			t.Errorf("Want %v, got %v", n1, got)
		}

		if got := str.GetNodeById("2"); !n2.Equals(got) {
			t.Errorf("Want %v, got %v", n2, got)
		}

		if got := str.GetNodeById("3"); !n3.Equals(got) {
			t.Errorf("Want %v, got %v", n3, got)
		}

		if got := str.GetNodeById("4"); !n4.Equals(got) {
			t.Errorf("Want %v, got %v", n4, got)
		}

		if got := str.GetNodeById("5"); !n5.Equals(got) {
			t.Errorf("Want %v, got %v", n5, got)
		}

		if got := str.GetNodeById("6"); !n6.Equals(got) {
			t.Errorf("Want %v, got %v", n6, got)
		}
	})

	t.Run("generates the elements", func(t *testing.T) {
		if got := str.ElementsCount(); got != 6 {
			t.Errorf("Want %d elements, got %d", 6, got)
		}

		if got := str.GetElementById("1"); !got.Equals(e1) {
			t.Errorf("Want %v, got %v", e1, got)
		}

		if got := str.GetElementById("2"); !got.Equals(e2) {
			t.Errorf("Want %v, got %v", e2, got)
		}

		if got := str.GetElementById("3"); !got.Equals(e3) {
			t.Errorf("Want %v, got %v", e3, got)
		}

		if got := str.GetElementById("4"); !got.Equals(e4) {
			t.Errorf("Want %v, got %v", e4, got)
		}

		if got := str.GetElementById("5"); !got.Equals(e5) {
			t.Errorf("Want %v, got %v", e5, got)
		}

		if got := str.GetElementById("6"); !got.Equals(e6) {
			t.Errorf("Want %v, got %v", e6, got)
		}
	})

	t.Run("doesn't add loads to the columns", func(t *testing.T) {
		b1 := str.GetElementById("1")
		if len(b1.ConcentratedLoads) > 0 || len(b1.DistributedLoads) > 0 {
			t.Error("Expected bar to not have any load applied")
		}

		b2 := str.GetElementById("2")
		if len(b2.ConcentratedLoads) > 0 || len(b2.DistributedLoads) > 0 {
			t.Error("Expected bar to not have any load applied")
		}

		b4 := str.GetElementById("4")
		if len(b4.ConcentratedLoads) > 0 || len(b4.DistributedLoads) > 0 {
			t.Error("Expected bar to not have any load applied")
		}

		b5 := str.GetElementById("5")
		if len(b5.ConcentratedLoads) > 0 || len(b5.DistributedLoads) > 0 {
			t.Error("Expected bar to not have any load applied")
		}
	})

	t.Run("adds horizontal distributed loads to the beams", func(t *testing.T) {
		b3 := str.GetElementById("3")
		if len(b3.ConcentratedLoads) > 0 {
			t.Error("Expected bar to not have any concentrated load applied")
		}
		if len(b3.DistributedLoads) != 1 || !b3.DistributedLoads[0].Equals(load) {
			t.Error("Wrong distributed load")
		}

		b6 := str.GetElementById("6")
		if len(b6.ConcentratedLoads) > 0 {
			t.Error("Expected bar to not have any concentrated load applied")
		}
		if len(b6.DistributedLoads) != 1 || !b6.DistributedLoads[0].Equals(load) {
			t.Error("Wrong distributed load")
		}
	})
}
