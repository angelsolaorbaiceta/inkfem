package plot

import (
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

// structureRectBounds computes the rectangle containing the structure's geometry
// after applying the chosen drawing scale and adding, at least, the minimum margin
// specified in the options.
//
// To scale the box that contains the structure, the innerScale factor and the
// plot options scale are multiplied.
func structureRectBounds(
	st *structure.Structure,
	options *StructurePlotOps,
	unitsScale unitsScale,
) *g2d.Rect {
	nodePositions := make([]*g2d.Point, st.NodesCount())

	var failIfError = func(err error) {
		if err != nil {
			panic(err)
		}
	}

	for i, node := range st.GetAllNodes() {
		nodePositions[i] = node.Position
	}

	rect, err := g2d.MakeRectContaining(nodePositions)
	failIfError(err)

	rect, err = rect.WithScaledSize(options.Scale * unitsScale.value())
	failIfError(err)

	rect, err = rect.WithMargins(float64(options.MinMargin), float64(options.MinMargin))
	failIfError(err)

	return rect
}
