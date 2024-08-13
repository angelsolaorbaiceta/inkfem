package plot

import (
	"fmt"

	svg "github.com/ajstarks/svgo"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

func drawExternalConstraints(
	canvas *svg.SVG,
	st *structure.Structure,
	config *plotConfig,
) {
	// TODO: Derive the length from the average element length.
	l := 50

	canvas.Gstyle(
		fmt.Sprintf("stroke:%s;stroke-width:%d", config.ExternalConstColor, config.ExternalConstWidth),
	)

	for _, node := range st.GetAllNodes() {
		if node.IsExternallyConstrained() {
			// A group that sets the coordinate origin at the node's position and the
			// y-axis pointing downwards.
			canvas.Gtransform(
				fmt.Sprintf("translate(%d,%d) scale(1,-1)",
					int(node.Position.X()), int(node.Position.Y()),
				),
			)

			if node.ExternalConstraint == &structure.FullConstraint {
				drawGround(canvas, l)
			}

			canvas.Gend()
		}
	}

	canvas.Gend()
}

func drawGround(canvas *svg.SVG, l int) {
	halfL := l / 2

	canvas.Rect(-halfL, 0, l, halfL, "stroke=\"none\"", fmt.Sprintf("fill=\"url(#%s)\"", diagonalLinesPatternId))
	canvas.Line(-halfL, 0, halfL, 0, "vector-effect=\"non-scaling-stroke\"")
}
