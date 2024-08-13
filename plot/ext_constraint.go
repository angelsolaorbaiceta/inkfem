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
	l := 75

	canvas.Gstyle(
		fmt.Sprintf("stroke:%s;stroke-width:%d;fill:none", config.ExternalConstColor, config.ExternalConstWidth),
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

			if node.ExternalConstraint.Equals(&structure.FullConstraint) {
				drawGround(canvas, l, 0)
			} else if node.ExternalConstraint.Equals(&structure.DispConstraint) {
				drawTriangle(canvas, l)
				drawGround(canvas, l, int(l/2))
			} else if node.ExternalConstraint.Equals(&structure.DispYConstraint) {
				drawTriangle(canvas, l)
				drawWheels(canvas, l)
				drawGround(canvas, l, int(3*l/4))
			}

			canvas.Gend()
		}
	}

	canvas.Gend()
}

func drawTriangle(canvas *svg.SVG, l int) {
	var (
		halfL   = int(l / 2)
		fourthL = int(l / 4)
	)

	canvas.Polygon(
		[]int{0, -fourthL, fourthL},
		[]int{0, halfL, halfL},
	)
}

func drawWheels(canvas *svg.SVG, l int) {
	var (
		r      = int(l / 8)
		y      = int(l/2) + r
		leftX  = -int(l/4) + r
		rightX = int(l/4) - r
	)

	canvas.Circle(leftX, y, r)
	canvas.Circle(rightX, y, r)
}

func drawGround(canvas *svg.SVG, l int, deltaY int) {
	var (
		halfL = l / 2
		y     = deltaY
	)

	canvas.Rect(
		-halfL, y, l, halfL,
		"stroke=\"none\"",
		fmt.Sprintf("fill=\"url(#%s)\"", diagonalLinesPatternId),
	)
	canvas.Line(-halfL, y, halfL, y)
}
