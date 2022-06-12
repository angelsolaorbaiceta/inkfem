package def

import (
	"testing"

	inkio "github.com/angelsolaorbaiceta/inkfem/io"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

func TestReadDefinition(t *testing.T) {
	var (
		wantStr = inkio.MakeTestOriginalStructure()
		reader  = inkio.MakeTestDefinitionReader()
		str     = Read(reader, inkio.ReaderOptions{ShouldIncludeOwnWeight: false})
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

	t.Run("parses the nodes", func(t *testing.T) {
		var (
			wantN1 = wantStr.GetNodeById("n1")
			wantN2 = wantStr.GetNodeById("n2")
		)

		if got := str.GetNodeById("n1"); !got.Equals(wantN1) {
			t.Errorf("Want %v, got %v", wantN1, got)
		}
		if got := str.GetNodeById("n2"); !got.Equals(wantN2) {
			t.Errorf("Want %v, got %v", wantN2, got)
		}
	})

	t.Run("parses the materials", func(t *testing.T) {
		wantMaterial := wantStr.GetMaterialsByName()["mat_yz"]

		if str.MaterialsCount() != 1 {
			t.Error("Want one material")
		}

		if got := str.GetMaterialsByName()[wantMaterial.Name]; !got.Equals(wantMaterial) {
			t.Errorf("Want %v, got %v", wantMaterial, got)
		}
	})

	t.Run("parses the sections", func(t *testing.T) {
		wantSection := wantStr.GetSectionsByName()["sec_xy"]

		if str.SectionsCount() != 1 {
			t.Error("Want one section")
		}

		if got := str.GetSectionsByName()[wantSection.Name]; !got.Equals(wantSection) {
			t.Errorf("Want %v, got %v", wantSection, got)
		}
	})

	t.Run("parses the loads", func(t *testing.T) {
		var (
			wantConcLoad = wantStr.GetElementById("b1").ConcentratedLoads[0]
			wantDistLoad = wantStr.GetElementById("b1").DistributedLoads[0]
		)

		if got := str.GetElementById("b1").ConcentratedLoads[0]; !got.Equals(wantConcLoad) {
			t.Errorf("Want %v, got %v", wantConcLoad, got)
		}
		if got := str.GetElementById("b1").DistributedLoads[0]; !got.Equals(wantDistLoad) {
			t.Errorf("Want %v, got %v", wantDistLoad, got)
		}
	})

	t.Run("parses the bars", func(t *testing.T) {
		wantBar := wantStr.GetElementById("b1")

		if got := str.GetElementById("b1"); !got.Equals(wantBar) {
			t.Errorf("Want %v, got %v", wantBar, got)
		}
	})
}
