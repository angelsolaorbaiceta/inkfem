package plot

import (
	"fmt"
	"io"
	"math"

	svg "github.com/ajstarks/svgo"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

const (
	diagonalLinesPatternId = "diagonalLines"
	loadArrowMarkerId      = "loadArrow"
)

// StructurePlotOps are the options that can be passed to the StructureToSVG
// function to tweak how the structure is drawn.
//
// These options are passed by the user when calling the plot command.
type StructurePlotOps struct {
	// Scale is the scale factor applied to the structure geometry.
	Scale float64
	// DistLoadScale is the scale factor applied to the distributed loads.
	DistLoadScale float64
	// MinMargin is the minimum margin between the structure and the canvas border.
	MinMargin int
}

func (ops StructurePlotOps) ApplyDrawingScale(value int) int {
	return int(math.Ceil(float64(value) * ops.Scale))
}

// plotContext is a structure that holds the context for the plot drawing functions.
// It includes the canvas where the drawing is done, the configuration for the plot,
// and the units scale to draw the elements with the right size.
type plotContext struct {
	canvas     *svg.SVG
	config     *PlotConfig
	options    *StructurePlotOps
	unitsScale unitsScale
}

// StructureToSVG generates an SVG diagram representing the structure's definition
// and writes the result to the given writer.
func StructureToSVG(
	st *structure.Structure,
	options *StructurePlotOps,
	config *PlotConfig,
	w io.Writer,
) {
	var (
		// The units scale is a scale factor to account for the units used to define the
		// structure. It makes sure that the size of the bars is adequate for the plot.
		unitsScale = determineUnitsScale(st)
		rectBounds = structureRectBounds(st, options, unitsScale)
		canvas     = svg.New(w)
		ctx        = plotContext{
			canvas:     canvas,
			config:     config,
			options:    options,
			unitsScale: unitsScale,
		}
	)

	canvas.Start(int(rectBounds.Width()), int(rectBounds.Height()))

	canvas.Def()
	defineExtConstrainGroundPattern(canvas, config)
	defineLoadArrowMarker(canvas, config)
	canvas.DefEnd()

	// The canvas is scaled so the Y axis points upwards and the scale is applied.
	// The origin is set at the bottom left corner of the canvas, including a margin.
	canvas.Gtransform(
		transformMatrix(
			options.Scale, -options.Scale,
			float64(options.MinMargin), rectBounds.Height()-float64(options.MinMargin),
		),
	)
	drawLoads(st, &ctx)
	drawGeometry(st, &ctx)
	drawExternalConstraints(st, &ctx)
	canvas.Gend()
	canvas.End()
}

func defineLoadArrowMarker(canvas *svg.SVG, config *PlotConfig) {
	size := config.DistLoadArrowSize

	canvas.Marker(
		loadArrowMarkerId,
		9,          // refX
		2,          // refY
		size, size, // markerWidth, markerHeight
		"orient=\"auto\"",
		"markerUnits=\"userSpaceOnUse\"",
		"viewBox=\"-2 -2 12 12\"",
	)
	canvas.Path(
		"M0,0 L9,2 L0,4",
		"vector-effect=\"non-scaling-stroke\"",
		fmt.Sprintf(
			"stroke:%s;stroke-width:%d;fill:none",
			config.DistLoadColor, config.DistLoadWidth,
		),
	)
	canvas.MarkerEnd()
}
