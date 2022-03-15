package generate

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

func TestGenerateReticularStructure(t *testing.T) {
	var (
		params = ReticStructureParams{
			Spans:    1,
			Levels:   2,
			Span:     300.0,
			Height:   200.0,
			Section:  structure.MakeUnitSection(),
			Material: structure.MakeUnitMaterial(),
		}
		str = Reticular(params)
	)

	t.Run("generates the nodes", func(t *testing.T) {
		n1 := str.GetNodeById("1")
		if !g2d.MakePoint(0, 0).Equals(n1.Position) {
			t.Errorf("Wrong node position: %v", n1.Position)
		}
		if n1.ExternalConstraint != &structure.FullConstraint {
			t.Error("Wrong node constraint")
		}

		n2 := str.GetNodeById("2")
		if !g2d.MakePoint(params.Span, 0).Equals(n2.Position) {
			t.Errorf("Wrong node position: %v", n2.Position)
		}
		if n2.ExternalConstraint != &structure.FullConstraint {
			t.Error("Wrong node constraint")
		}

		n3 := str.GetNodeById("3")
		if !g2d.MakePoint(0, params.Height).Equals(n3.Position) {
			t.Errorf("Wrong node position: %v", n3.Position)
		}
		if n3.ExternalConstraint != &structure.NilConstraint {
			t.Error("Wrong node constraint")
		}

		n4 := str.GetNodeById("4")
		if !g2d.MakePoint(params.Span, params.Height).Equals(n4.Position) {
			t.Errorf("Wrong node position: %v", n4.Position)
		}
		if n4.ExternalConstraint != &structure.NilConstraint {
			t.Error("Wrong node constraint")
		}
	})
}
