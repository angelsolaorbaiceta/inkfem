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

// StructurePlotOps are the options that can be passed to the StructureToSVG function to tweak
// how the structure definition is drawn.
type StructurePlotOps struct {
	Scale     float64
	MinMargin int
}

// StructureToSVG generates an SVG diagram representing the structure's definition and writes
// the result to the given writer.
func StructureToSVG(st *structure.Structure, options *StructurePlotOps, w io.Writer) {
	var (
		rectBounds = structureRectBounds(st, options)
		canvas     = svg.New(w)
		config     = defaultPlotConfig()
	)

	canvas.Start(int(rectBounds.Width()), int(rectBounds.Height()))

	canvas.Def()
	defineExtConstrainGroundPattern(canvas, config)
	canvas.DefEnd()

	canvas.Gtransform(
		fmt.Sprintf(
			"matrix(%f,0,0,%f,%d,%f)",
			options.Scale,
			-options.Scale,
			options.MinMargin,
			rectBounds.Height()-float64(options.MinMargin),
		),
	)
	drawGeometry(canvas, st, config)
	drawExternalConstraints(canvas, st, config)
	canvas.Gend()
	canvas.End()
}

func defineExtConstrainGroundPattern(canvas *svg.SVG, config *plotConfig) {
	canvas.Pattern(diagonalLinesPatternId, 0, 0, 10, 10, "user", "patternTransform=\"rotate(30)\"")
	canvas.Line(0, 0, 0, 10, fmt.Sprintf("stroke:%s;stroke-width:%d", config.ExternalConstColor, config.ExternalConstWidth))
	canvas.PatternEnd()
}
