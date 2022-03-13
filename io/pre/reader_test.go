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
		want := structure.StrMetadata{MajorVersion: 2, MinorVersion: 3}
		if got := str.Metadata; got.MajorVersion != want.MajorVersion || got.MinorVersion != want.MinorVersion {
			t.Errorf("Want %v, got %v", want, got)
		}
	})
}
