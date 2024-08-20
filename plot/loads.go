package plot

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/g2d"
)

func drawLoads(st *structure.Structure, ctx *plotContext) {
	for _, bar := range st.Elements() {
		if bar.HasLoadsApplied() {
			drawLoadForBar(bar, ctx)
		}
	}
}

func drawLoadForBar(bar *structure.Element, ctx *plotContext) {
	// Set a group whose reference frame is that of the bar, which goes from the
	// start node to the end node.
	var (
		canvas = ctx.canvas
		scale  = ctx.unitsScale
	)

	canvas.Gtransform(transformToLocalBar(bar, scale))
	drawDistributedLoads(bar, ctx)
	canvas.Gend()
}

// drawDistributedLoads draws all the distributed loads of a bar element.
func drawDistributedLoads(bar *structure.Element, ctx *plotContext) {
	var (
		canvas = ctx.canvas
		config = ctx.config
	)

	canvas.Gstyle(
		fmt.Sprintf(
			"stroke-width:%d;stroke:%s;fill:%s",
			config.DistLoadWidth, config.DistLoadColor, config.DistLoadFillColor,
		),
	)

	for _, dLoad := range bar.DistributedLoads {
		drawDistributedLoad(dLoad, bar.Geometry(), ctx)
	}

	canvas.Gend()
}

// drawDistributedLoad draws a distributed load in the bar element.
func drawDistributedLoad(
	dLoad *load.DistributedLoad,
	barGeometry *g2d.Segment,
	ctx *plotContext,
) {
	if dLoad.IsInLocalCoords {
		switch dLoad.Term {
		case load.FX:
			drawLocalDistributedFxLoad(dLoad, barGeometry, ctx)
		case load.FY:
			drawLocalDistributedFyLoad(dLoad, barGeometry, ctx)
		case load.MZ:
			drawLocalDistributedMzLoad(dLoad, barGeometry, ctx)
		}
	} else {
		// TODO: draw distributed load in global coords
	}
}
