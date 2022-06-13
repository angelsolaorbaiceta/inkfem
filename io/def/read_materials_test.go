package def

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

func TestDeserializeMaterial(t *testing.T) {
	t.Run("deserializes the material", func(t *testing.T) {
		var (
			got      = DeserializeMaterial("'mat steel' -> 1.1 2.2 3.3 4.4 5.5 6.6")
			wantName = "mat steel"
			want     = structure.MakeMaterial(wantName, 1.1, 2.2, 3.3, 4.4, 5.5, 6.6)
		)

		if got.Name != wantName {
			t.Errorf("Expected name %s, got '%s'", wantName, got.Name)
		}
		if !got.Equals(want) {
			t.Errorf("Wrong material. Want %v, got %v", want, got)
		}
	})

	t.Run("deserializes the material using scientific notation numbers", func(t *testing.T) {
		var (
			got  = DeserializeMaterial("'steel' -> 1.1e2 2.2e-2 3e3 4.4 5.5 6.6")
			want = structure.MakeMaterial("steel", 110.0, 0.022, 3000, 4.4, 5.5, 6.6)
		)

		if !got.Equals(want) {
			t.Errorf("Wrong material. Want %v, got %v", want, got)
		}
	})
}
