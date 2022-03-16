package generate

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
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
	)

	t.Run("generates the nodes", func(t *testing.T) {
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
		if got := str.GetElementById("1"); !got.Equals(e1) {
			t.Errorf("Want %v, got %v", e1, got)
		}

		if got := str.GetElementById("2"); !got.Equals(e2) {
			t.Errorf("Want %v, got %v", e2, got)
		}
	})
}
