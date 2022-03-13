package plot

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

// structureRectBounds computes the rectangle containing the structure's geometry after applying
// the chosen drawing scale and adding, at least, the minimum margin specified in the options.
func structureRectBounds(st *structure.Structure, options StructurePlotOps) *g2d.Rect {
	nodePositions := make([]*g2d.Point, 0, st.NodesCount())

	var failIfError = func(err error) {
		if err != nil {
			panic("Can't compute the structure's rectangular bounds")
		}
	}

	for _, node := range st.GetAllNodes() {
		nodePositions = append(nodePositions, node.Position)
	}

	rect, err := g2d.MakeRectContaining(nodePositions)
	failIfError(err)

	rect, err = rect.WithScaledSize(options.Scale)
	failIfError(err)

	rect, err = rect.WithMargins(float64(options.MinMargin), float64(options.MinMargin))
	failIfError(err)

	return rect
}
