package plot

import (
	"fmt"
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
	)

	canvas.Start(int(rectBounds.Width()), int(rectBounds.Height()))

	canvas.Def()
	defineExtConstrainGroundPattern(canvas, config)
	canvas.DefEnd()

	// The canvas is scaled so the Y axis points upwards and the scale is applied.
	// The origin is set at the bottom left corner of the canvas, including a margin.
	canvas.Gtransform(
		fmt.Sprintf(
			"matrix(%f,0,0,%f,%d,%f)",
			options.Scale,
			-options.Scale,
			options.MinMargin,
			rectBounds.Height()-float64(options.MinMargin),
		),
	)
	drawGeometry(canvas, st, config, unitsScale)
	drawExternalConstraints(canvas, st, config, unitsScale)
	canvas.Gend()
	canvas.End()
}
