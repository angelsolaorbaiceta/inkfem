package def

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

func TestDeserializeNode(t *testing.T) {
	t.Run("deserializes the node", func(t *testing.T) {
		var (
			got  = DeserializeNode("1 -> 10.1 20.2 { dx dy rz }")
			want = structure.MakeNode("1", g2d.MakePoint(10.1, 20.2), &structure.FullConstraint)
		)

		if !got.Equals(want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})

	t.Run("deserializes the node with scientific notation coordinates", func(t *testing.T) {
		var (
			got  = DeserializeNode("1 -> 1e+2 2.0e-2 { dx dy rz }")
			want = structure.MakeNode("1", g2d.MakePoint(100.0, 0.02), &structure.FullConstraint)
		)

		if !got.Equals(want) {
			t.Errorf("Expected %v, got %v", want, got)
		}
	})
}
