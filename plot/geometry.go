package plot

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

// drawGeometry draws the geometry of the structural elements in the given canvas, including
// them inside a group with the adequate affine transformation that results in the y-axis
// pointing upwards and the chosen drawing scale.
func drawGeometry(st *structure.Structure, ctx *plotContext) {
	var (
		scale  = ctx.unitsScale
		canvas = ctx.canvas
		config = ctx.config
	)

	canvas.Gstyle(
		fmt.Sprintf("stroke:%s;stroke-width:%d;fill:none",
			config.GeometryColor,
			config.GeometryWidth,
		),
	)

	var (
		startPoint, endPoint       *g2d.Point
		startX, startY, endX, endY int
	)

	for _, element := range st.Elements() {
		startPoint = scale.applyToPoint(element.StartPoint())
		endPoint = scale.applyToPoint((element.EndPoint()))
		startX, startY = int(startPoint.X()), int(startPoint.Y())
		endX, endY = int(endPoint.X()), int(endPoint.Y())

		canvas.Line(
			startX, startY, endX, endY,
			fmt.Sprintf("id=\"bar__%s\"", element.GetID()),
		)
	}

	for _, node := range st.GetAllNodes() {
		point := scale.applyToPoint(node.Position)
		x, y := int(point.X()), int(point.Y())

		canvas.Circle(
			x, y, config.NodeRadius,
			fmt.Sprintf("id=\"node__%s\"", node.GetID()),
		)
	}

	canvas.Gend()
}
