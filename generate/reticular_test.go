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
		if got := n1.Position; !g2d.MakePoint(0, 0).Equals(got) {
			t.Error("Wrong node position")
		}
	})
}
