package def

import (
	"testing"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/stretchr/testify/assert"
)

func TestReadDefinition(t *testing.T) {
	var (
		wantStr = inkio.MakeTestOriginalStructure()
		reader  = inkio.MakeTestDefinitionReader()
		str     = Read(reader)
	)

	t.Run("parses the metadata", func(t *testing.T) {
		var (
			want = structure.StrMetadata{MajorVersion: 2, MinorVersion: 3}
			got  = str.Metadata
		)

		assert.Equal(t, want.MajorVersion, got.MajorVersion)
		assert.Equal(t, want.MinorVersion, got.MinorVersion)
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

	t.Run("parses the loads", func(t *testing.T) {
		var (
			wantConcLoad = wantStr.GetElementById("b1").ConcentratedLoads[0]
			wantDistLoad = wantStr.GetElementById("b1").DistributedLoads[0]
		)

		assert.Equal(t, wantConcLoad, str.GetElementById("b1").ConcentratedLoads[0])
		assert.Equal(t, wantDistLoad, str.GetElementById("b1").DistributedLoads[0])
	})

	t.Run("parses the bars", func(t *testing.T) {
		wantBar := wantStr.GetElementById("b1")

		assert.Equal(t, wantBar, str.GetElementById("b1"))
	})
}

func TestReadDefinitionInverseOrder(t *testing.T) {
	var (
		wantStr = inkio.MakeTestOriginalStructure()
		reader  = inkio.MakeTestDefinitionReaderInverseOrder()
		str     = Read(reader)
	)

	t.Run("parses the nodes", func(t *testing.T) {
		var (
			wantN1 = wantStr.GetNodeById("n1")
			wantN2 = wantStr.GetNodeById("n2")
		)

		assert.Equal(t, wantN1, str.GetNodeById("n1"))
		assert.Equal(t, wantN2, str.GetNodeById("n2"))
	})
}
