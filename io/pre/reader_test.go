package pre

import (
	"testing"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/stretchr/testify/assert"
)

func TestReadPreprocessModel(t *testing.T) {
	var (
		wantStr            = inkio.MakeTestPreprocessedStructure()
		preprocessedReader = inkio.MakeTestPreprocessedReader()
		str                = Read(preprocessedReader)
	)

	t.Run("parses the metadata", func(t *testing.T) {
		var (
			want = structure.StrMetadata{MajorVersion: 2, MinorVersion: 3}
			got  = str.Metadata
		)

		assert.Equal(t, want, got)
	})

	t.Run("parses the degrees of freedom count", func(t *testing.T) {
		assert.Equal(t, 9, str.DofsCount())
	})

	t.Run("parses the nodes", func(t *testing.T) {
		var (
			wantN1 = wantStr.GetNodeById("n1")
			wantN2 = wantStr.GetNodeById("n2")
		)

		assert.Equal(t, wantN1, str.GetNodeById("n1"))
		assert.Equal(t, wantN2, str.GetNodeById("n2"))
	})

	t.Run("parses the materials", func(t *testing.T) {
		wantMaterial := wantStr.GetMaterialsByName()["mat_yz"]

		assert.Equal(t, 1, str.MaterialsCount())
		assert.Equal(t, wantMaterial, str.GetMaterialsByName()[wantMaterial.Name])
	})

	t.Run("parses the sections", func(t *testing.T) {
		wantSection := wantStr.GetSectionsByName()["sec_xy"]

		assert.Equal(t, 1, str.SectionsCount())
		assert.Equal(t, wantSection, str.GetSectionsByName()[wantSection.Name])
	})

	t.Run("parses the bars", func(t *testing.T) {
		wantBar := wantStr.GetElementById("b1")

		assert.True(t, wantBar.Equals(str.GetElementById("b1")))
	})
}
