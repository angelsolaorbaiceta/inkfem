package plot

import (
	"io"

	svg "github.com/ajstarks/svgo"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

const (
	diagonalLinesPatternId = "diagonalLines"
)

// StructurePlotOps are the options that can be passed to the StructureToSVG
// function to tweak how the structure is drawn.
type StructurePlotOps struct {
	Scale     float64
	MinMargin int
}

// plotContext is a structure that holds the context for the plot drawing functions.
// It includes the canvas where the drawing is done, the configuration for the plot,
// and the units scale to draw the elements with the right size.
type plotContext struct {
	canvas     *svg.SVG
	config     *plotConfig
	unitsScale unitsScale
}

// StructureToSVG generates an SVG diagram representing the structure's definition
// and writes the result to the given writer.
func StructureToSVG(st *structure.Structure, options *StructurePlotOps, w io.Writer) {
	var (
		// The units scale is a scale factor to account for the units used to define the
		// structure. It makes sure that the size of the bars is adequate for the plot.
		unitsScale = determineUnitsScale(st)
		rectBounds = structureRectBounds(st, options, unitsScale)
		canvas     = svg.New(w)
		config     = defaultPlotConfig()
		ctx        = plotContext{
			canvas:     canvas,
			config:     config,
			unitsScale: unitsScale,
		}
	)

	canvas.Start(int(rectBounds.Width()), int(rectBounds.Height()))

	canvas.Def()
	defineExtConstrainGroundPattern(canvas, config)
	canvas.DefEnd()

	// The canvas is scaled so the Y axis points upwards and the scale is applied.
	// The origin is set at the bottom left corner of the canvas, including a margin.
	canvas.Gtransform(
		transformMatrix(
			options.Scale, -options.Scale,
			float64(options.MinMargin), rectBounds.Height()-float64(options.MinMargin),
		),
	)
	drawGeometry(st, &ctx)
	drawExternalConstraints(st, &ctx)
	canvas.Gend()
	canvas.End()
}
