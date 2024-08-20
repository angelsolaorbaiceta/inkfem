package plot

import (
	"bytes"
	"strings"
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

			wantPolygon = "<polygon points=\"20,0 20,-200 80,-200 80,0\" />"
			wantArrow1  = "<line x1=\"20\" y1=\"-200\" x2=\"20\" y2=\"0\" marker-end=\"url(#loadArrow)\" stroke=\"#558B2F\" />"
			wantArrow2  = "<line x1=\"30\" y1=\"-200\" x2=\"30\" y2=\"0\" marker-end=\"url(#loadArrow)\" stroke=\"#558B2F\" />"
			wantArrow3  = "<line x1=\"40\" y1=\"-200\" x2=\"40\" y2=\"0\" marker-end=\"url(#loadArrow)\" stroke=\"#558B2F\" />"
			wantArrow4  = "<line x1=\"50\" y1=\"-200\" x2=\"50\" y2=\"0\" marker-end=\"url(#loadArrow)\" stroke=\"#558B2F\" />"
			wantArrow5  = "<line x1=\"60\" y1=\"-200\" x2=\"60\" y2=\"0\" marker-end=\"url(#loadArrow)\" stroke=\"#558B2F\" />"
			wantArrow6  = "<line x1=\"70\" y1=\"-200\" x2=\"70\" y2=\"0\" marker-end=\"url(#loadArrow)\" stroke=\"#558B2F\" />"
			wantArrow7  = "<line x1=\"80\" y1=\"-200\" x2=\"80\" y2=\"0\" marker-end=\"url(#loadArrow)\" stroke=\"#558B2F\" />"
		)

		drawLocalDistributedFyLoad(dLoad, barGeometry, context)

		// split lines
		var (
			gotLines  = strings.Split(writer.String(), "\n")
			gotPoly   = gotLines[0]
			gotArrow1 = gotLines[3]
			gotArrow2 = gotLines[4]
			gotArrow3 = gotLines[5]
			gotArrow4 = gotLines[6]
			gotArrow5 = gotLines[7]
			gotArrow6 = gotLines[8]
			gotArrow7 = gotLines[9]
		)

		assert.Equal(t, wantPolygon, gotPoly)
		assert.Equal(t, wantArrow1, gotArrow1)
		assert.Equal(t, wantArrow2, gotArrow2)
		assert.Equal(t, wantArrow3, gotArrow3)
		assert.Equal(t, wantArrow4, gotArrow4)
		assert.Equal(t, wantArrow5, gotArrow5)
		assert.Equal(t, wantArrow6, gotArrow6)
		assert.Equal(t, wantArrow7, gotArrow7)
	})
}
