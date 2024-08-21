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

func TestDrawLocalDistributedFxLoad(t *testing.T) {
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

	t.Run("Constant Fx positive load, partial length", func(t *testing.T) {
		var (
			writer  bytes.Buffer
			context = makeContext(&writer, 1.0)
			startT  = nums.MakeTParam(0.2)
			endT    = nums.MakeTParam(0.8)
			dLoad   = load.MakeDistributed(load.FX, true, startT, 200, endT, 200)

			wantPolygon   = "<polygon points=\"20,0 20,200 80,200 80,0\" />"
			wantArrowYPos = []int{40, 80, 120, 160}
			wantArrows    = make([]string, len(wantArrowYPos))
		)

		for i, y := range wantArrowYPos {
			wantArrows[i] = fmt.Sprintf(
				"<line x1=\"20\" y1=\"%d\" x2=\"80\" y2=\"%d\" marker-end=\"url(#%s)\" stroke=\"%s\" />",
				y, y, arrowId, stroke,
			)
		}

		drawLocalDistributedFxLoad(dLoad, barGeometry, context)

		var (
			gotLines  = strings.Split(writer.String(), "\n")
			gotPoly   = gotLines[0]
			gotArrows = gotLines[3 : len(wantArrows)+3]
		)

		assert.Equal(t, wantPolygon, gotPoly)
		assert.Equal(t, wantArrows, gotArrows)
	})

	t.Run("Decreasing Fx positive load, full length", func(t *testing.T) {
		var (
			writer  bytes.Buffer
			context = makeContext(&writer, 1.0)
			startT  = nums.MinT
			endT    = nums.MaxT
			dLoad   = load.MakeDistributed(load.FX, true, startT, 200, endT, 120)

			wantPolygon      = "<polygon points=\"0,0 0,200 100,120 100,0\" />"
			wantArrowYPos    = []int{40, 80, 120, 160}
			wantArrowXEndPos = []int{100, 100, 100, 50}
			wantArrows       = make([]string, len(wantArrowYPos))
		)

		for i, y := range wantArrowYPos {
			x := wantArrowXEndPos[i]
			wantArrows[i] = fmt.Sprintf(
				"<line x1=\"0\" y1=\"%d\" x2=\"%d\" y2=\"%d\" marker-end=\"url(#%s)\" stroke=\"%s\" />",
				y, x, y, arrowId, stroke,
			)
		}

		drawLocalDistributedFxLoad(dLoad, barGeometry, context)
		fmt.Println(writer.String())

		var (
			gotLines  = strings.Split(writer.String(), "\n")
			gotPoly   = gotLines[0]
			gotArrows = gotLines[3 : len(wantArrows)+3]
		)

		assert.Equal(t, wantPolygon, gotPoly)
		assert.Equal(t, wantArrows, gotArrows)
	})
}
