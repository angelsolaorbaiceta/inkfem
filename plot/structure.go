package plot

import (
	"io"

	svg "github.com/ajstarks/svgo"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
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
	)

	canvas.Start(int(rectBounds.Width()), int(rectBounds.Height()))
	// canvas.Circle(width/2, height/2, 100)
	// canvas.Text(width/2, height/2, "Hello, SVG", "text-anchor:middle;font-size:30px;fill:white")
	canvas.End()
}

func structureRectBounds(st *structure.Structure, options StructurePlotOps) *g2d.Rect {
	nodePositions := make([]*g2d.Point, 0, st.NodesCount())

	for _, node := range st.Nodes {
		nodePositions = append(nodePositions, node.Position)
	}

	rect, err := g2d.MakeRectContaining(nodePositions)
	if err != nil {
		panic("Can't compute the structure's rectangular bounds")
	}

	rect, err = rect.WithScaledSize(options.Scale)
	if err != nil {
		panic("Can't compute the structure's rectangular bounds")
	}

	rect, err = rect.WithMargins(100.0, 100.0)
	if err != nil {
		panic("Can't compute the structure's rectangular bounds")
	}

	return rect
}
