package plot

import (
	"fmt"

	svg "github.com/ajstarks/svgo"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

// drawGeometry draws the geometry of the structural elements in the given canvas, including
// them inside a group with the adequate affine transformation that results in the y-axis
// pointing upwards and the chosen drawing scale.
func drawGeometry(
	canvas *svg.SVG,
	st *structure.Structure,
	config *plotConfig,
) {
	canvas.Gstyle(
		fmt.Sprintf("stroke:%s;stroke-width:%d", config.GeometryColor, config.GeometryWidth),
	)

	var startPoint, endPoint *g2d.Point

	for _, element := range st.Elements() {
		startPoint = element.StartPoint()
		endPoint = element.EndPoint()

		canvas.Line(
			int(startPoint.X()), int(startPoint.Y()),
			int(endPoint.X()), int(endPoint.Y()),
			"vector-effect=\"non-scaling-stroke\"",
		)
	}

	canvas.Gend()
}
