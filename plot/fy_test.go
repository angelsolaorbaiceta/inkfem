package plot

import (
	"bytes"
	"fmt"
	"testing"

	svg "github.com/ajstarks/svgo"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/stretchr/testify/assert"
)

func TestDrawLocalDistributedFyLoad(t *testing.T) {
	var (
		writer  bytes.Buffer
		context = &plotContext{
			canvas: svg.New(&writer),
			config: DefaultPlotConfig(),
			options: &StructurePlotOps{
				Scale:         1.0,
				DistLoadScale: 1.0,
				MinMargin:     0,
			},
			unitsScale: unitsScale(1.0),
		}
		barGeometry = g2d.MakeSegment(g2d.MakePoint(0, 0), g2d.MakePoint(100, 0))
	)

	t.Cleanup(func() {
	})

	t.Run("Constant Fy positive load", func(t *testing.T) {
		var (
			startT = nums.MakeTParam(0.2)
			endT   = nums.MakeTParam(0.8)
			value  = 200.0
			dLoad  = load.MakeDistributed(load.FY, true, startT, value, endT, value)

			wantPolygon = "<polygon points=\"20,0 20,-200 80,-200 80,0\"\\s*/>"
		)

		drawLocalDistributedFyLoad(dLoad, barGeometry, context)

		got := writer.String()

		fmt.Printf("Got: %s\n", got)

		assert.Regexp(t, wantPolygon, got)
	})
}
