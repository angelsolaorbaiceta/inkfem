package plot

import (
	"fmt"
	"io"

	svg "github.com/ajstarks/svgo"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// StructurePlotOps are the options that can be passed to the StructureToSVG function to tweak
// how the structure definition is drawn.
type StructurePlotOps struct {
	Scale     float64
	MinMargin int
}

// StructureToSVG generates an SVG diagram representing the structure's definition and writes
// the result to the given writer.
func StructureToSVG(st *structure.Structure, options StructurePlotOps, w io.Writer) {
	var (
		rectBounds = structureRectBounds(st, options)
		canvas     = svg.New(w)
		config     = defaultPlotConfig()
	)

	canvas.Start(int(rectBounds.Width()), int(rectBounds.Height()))
	canvas.Gtransform(
		fmt.Sprintf(
			"matrix(%f,0,0,%f,%d,%f)",
			options.Scale,
			-options.Scale,
			options.MinMargin,
			rectBounds.Height()-float64(options.MinMargin),
		),
	)
	drawGeometry(canvas, st, &config)
	drawExternalConstraints(canvas, st, &config)
	canvas.Gend()
	canvas.End()
}
