package plot

import (
	"io"

	svg "github.com/ajstarks/svgo"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// StructurePlotOps are the options that can be passed to the StructureToSVG function to tweak
// how the structure definition is drawn.
type StructurePlotOps struct {
	Scale float64
}

// StructureToSVG generates an SVG diagram representing the structure's definition and writes
// the result to the given writer.
func StructureToSVG(st *structure.Structure, options StructurePlotOps, w io.Writer) {
	width := 500
	height := 500
	canvas := svg.New(w)
	canvas.Start(width, height)
	canvas.Circle(width/2, height/2, 100)
	canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")
	canvas.End()
}
