package def

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

func TestDeserializeSection(t *testing.T) {
	t.Run("deserializes the section", func(t *testing.T) {
		var (
			got      = DeserializeSection("'IPE 100' -> 1.1 2.2 3.3 4.4 5.5")
			wantName = "IPE 100"
			want     = structure.MakeSection(wantName, 1.1, 2.2, 3.3, 4.4, 5.5)
		)

		if got.Name != wantName {
			t.Errorf("Expected name '%s', got '%s'", wantName, got.Name)
		}
		if !got.Equals(want) {
			t.Errorf("Expected section %v, got %v", want, got)
		}
	})

	t.Run("deserializes the section using scientific notation numbers", func(t *testing.T) {
		var (
			got  = DeserializeSection("'IPE 100' -> 1.1e2 2.2e-2 3e3 4.4 5.5")
			want = structure.MakeSection("IPE 100", 110.0, 0.022, 3000, 4.4, 5.5)
		)

		if !got.Equals(want) {
			t.Errorf("Expected section %v, got %v", want, got)
		}
	})
}
