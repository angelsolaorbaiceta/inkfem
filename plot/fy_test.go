package plot

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"

	svg "github.com/ajstarks/svgo"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
	"github.com/stretchr/testify/assert"
)

func TestDrawLocalDistributedFyLoad(t *testing.T) {
	var makeContext = func(w io.Writer, loadScale float64) *plotContext {
		return &plotContext{
			canvas: svg.New(w),
			config: DefaultPlotConfig(),
			options: &StructurePlotOps{
				Scale:         1.0,
				DistLoadScale: loadScale,
				MinMargin:     0,
			},
			unitsScale: unitsScale(1.0),
		}
	}

	var (
		barGeometry = g2d.MakeSegment(g2d.MakePoint(0, 0), g2d.MakePoint(100, 0))
		arrowId     = loadArrowMarkerId
		stroke      = DefaultPlotConfig().DistLoadColor
	)

	t.Run("Constant Fy positive load, partial length", func(t *testing.T) {
		var (
			writer  bytes.Buffer
			context = makeContext(&writer, 1.0)
			startT  = nums.MakeTParam(0.2)
			endT    = nums.MakeTParam(0.8)
			value   = 200.0
			dLoad   = load.MakeDistributed(load.FY, true, startT, value, endT, value)

			wantPolygon   = "<polygon points=\"20,0 20,-200 80,-200 80,0\" />"
			wantArrowXPos = []int{20, 30, 40, 50, 60, 70, 80}
			wantArrows    = make([]string, len(wantArrowXPos))
		)

		for i, x := range wantArrowXPos {
			wantArrows[i] = fmt.Sprintf(
				"<line x1=\"%d\" y1=\"-200\" x2=\"%d\" y2=\"0\" marker-end=\"url(#%s)\" stroke=\"%s\" />",
				x, x, arrowId, stroke,
			)
		}

		drawLocalDistributedFyLoad(dLoad, barGeometry, context)

		var (
			gotLines  = strings.Split(writer.String(), "\n")
			gotPoly   = gotLines[0]
			gotArrows = gotLines[3:10]
		)

		assert.Equal(t, wantPolygon, gotPoly)
		assert.Equal(t, wantArrows, gotArrows)
	})

	t.Run("Fy positive to negative load, full length", func(t *testing.T) {
		var (
			writer  bytes.Buffer
			context = makeContext(&writer, 1.0)
			startT  = nums.MinT
			endT    = nums.MaxT
			startV  = 200.0
			endV    = -400.0
			dLoad   = load.MakeDistributed(load.FY, true, startT, startV, endT, endV)

			wantPolygon = "<polygon points=\"0,0 0,-200 100,400 100,0\" />"
			// The arrow at x=30 doesn't have enough space to draw the arrow.
			wantArrowXPos = []int{10, 20, 40, 50, 60, 70, 80, 90}
			wantArrows    = make([]string, len(wantArrowXPos))
		)

		for i, x := range wantArrowXPos {
			y := 6*x - 200

			wantArrows[i] = fmt.Sprintf(
				"<line x1=\"%d\" y1=\"%d\" x2=\"%d\" y2=\"0\" marker-end=\"url(#%s)\" stroke=\"%s\" />",
				x, y, x, arrowId, stroke,
			)
		}

		drawLocalDistributedFyLoad(dLoad, barGeometry, context)

		var (
			gotLines  = strings.Split(writer.String(), "\n")
			gotPoly   = gotLines[0]
			gotArrows = gotLines[3:11]
		)

		assert.Equal(t, wantPolygon, gotPoly)
		assert.Equal(t, wantArrows, gotArrows)
	})

	t.Run("Fy positive to negative load, full length, scaled", func(t *testing.T) {
		var (
			scale   = 5
			writer  bytes.Buffer
			context = makeContext(&writer, float64(scale))
			startT  = nums.MinT
			endT    = nums.MaxT
			startV  = 200.0
			endV    = -400.0
			dLoad   = load.MakeDistributed(load.FY, true, startT, startV, endT, endV)

			wantPolygon   = "<polygon points=\"0,0 0,-1000 100,2000 100,0\" />"
			wantArrowXPos = []int{10, 20, 30, 40, 50, 60, 70, 80, 90}
			wantArrows    = make([]string, len(wantArrowXPos))
		)

		for i, x := range wantArrowXPos {
			y := scale * (6*x - 200)

			wantArrows[i] = fmt.Sprintf(
				"<line x1=\"%d\" y1=\"%d\" x2=\"%d\" y2=\"0\" marker-end=\"url(#%s)\" stroke=\"%s\" />",
				x, y, x, arrowId, stroke,
			)
		}

		drawLocalDistributedFyLoad(dLoad, barGeometry, context)

		var (
			gotLines  = strings.Split(writer.String(), "\n")
			gotPoly   = gotLines[0]
			gotArrows = gotLines[3:12]
		)

		assert.Equal(t, wantPolygon, gotPoly)
		assert.Equal(t, wantArrows, gotArrows)
	})
}
