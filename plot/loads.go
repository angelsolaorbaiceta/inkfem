package plot

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

var (
	// distLoadLinePositions are the positions where the distributed load lines
	// with the arrowheads are drawn.
	distLoadLinePositions = []nums.TParam{
		nums.MakeTParam(0.1),
		nums.MakeTParam(0.2),
		nums.MakeTParam(0.3),
		nums.MakeTParam(0.4),
		nums.MakeTParam(0.5),
		nums.MakeTParam(0.6),
		nums.MakeTParam(0.7),
		nums.MakeTParam(0.8),
		nums.MakeTParam(0.9),
	}
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
		drawDistributedLoad(dLoad, bar, ctx)
	}

	canvas.Gend()
}

func drawDistributedLoad(
	dLoad *load.DistributedLoad,
	bar *structure.Element,
	ctx *plotContext,
) {
	if dLoad.IsInLocalCoords {
		drawLocalDistributedLoad(dLoad, bar, ctx)
	} else {
		// TODO: draw distributed load in global coords
	}
}

// The load is represented as a polygon with the following vertices:
//
//  1. The start node, which is at (0, 0) due to the transformation.
//  2. With x = 0, move y by the negative value of the load's start value.
func drawLocalDistributedLoad(
	dLoad *load.DistributedLoad,
	bar *structure.Element,
	ctx *plotContext,
) {
	var (
		canvas    = ctx.canvas
		scale     = ctx.unitsScale
		loadScale = ctx.options.DistLoadScale

		endX   = int(scale.applyToLength(bar.Length()))
		startY = int(-dLoad.StartValue * loadScale)
		endY   = int(-dLoad.EndValue * loadScale)

		x = []int{0, 0, endX, endX}
		y = []int{0, startY, endY, 0}
	)

	canvas.Polygon(x, y)
	canvas.Text(0, 0, fmt.Sprintf("%.2f", dLoad.StartValue), textTransform(0, startY), fmt.Sprintf("fill:%s", ctx.config.DistLoadColor))
	canvas.Text(0, 0, fmt.Sprintf("%.2f", dLoad.EndValue), textTransform(endX, endY), fmt.Sprintf("fill:%s", ctx.config.DistLoadColor))

	for _, t := range distLoadLinePositions {
		var (
			scaledLength = scale.applyToLength(bar.Length())
			loadX        = int(scaledLength * t.Value())
			loadY        = int(-dLoad.ValueAt(t) * loadScale)
		)

		if math.AbsInt(loadY) > ctx.config.DistLoadArrowSize {
			canvas.Line(
				loadX, loadY, loadX, 0,
				fmt.Sprintf("marker-end=\"url(#%s)\"", loadArrowMarkerId),
				fmt.Sprintf("stroke=\"%s\"", ctx.config.DistLoadColor),
			)
		}
	}
}

func textTransform(x, y int) string {
	var (
		translate = translate(float64(x), float64(y))
		scale     = scale(1, -1)
	)

	return fmt.Sprintf("transform=\"%s %s\"", translate, scale)
}
