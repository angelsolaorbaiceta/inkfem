package pre

import (
	"testing"

	"github.com/angelsolaorbaiceta/inkfem/structure"
)

func TestRead(t *testing.T) {
	var (
		originalStr        = makeTestOriginalStructure()
		preprocessedReader = makePreprocessedReader()
		str                = Read(*originalStr, preprocessedReader)
	)

	t.Run("parses the metadata", func(t *testing.T) {
		var (
			want = structure.StrMetadata{MajorVersion: 2, MinorVersion: 3}
			got  = str.Metadata
		)

		if got.MajorVersion != want.MajorVersion || got.MinorVersion != want.MinorVersion {
			t.Errorf("Want %v, got %v", want, got)
		}
	})

	t.Run("parses the degrees of freedom count", func(t *testing.T) {
		if got := str.DofsCount(); got != 9 {
			t.Errorf("Want 9 DOFs, got %d", got)
		}
	})
}
