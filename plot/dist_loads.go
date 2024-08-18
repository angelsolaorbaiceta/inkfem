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
	// with the arrowheads are drawn. We draw them at 10% intervals of the length.
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
		switch dLoad.Term {
		case load.FX:
			drawLocalDistributedFxLoad(dLoad, bar, ctx)
		case load.FY:
			drawLocalDistributedFyLoad(dLoad, bar, ctx)
		}
	} else {
		// TODO: draw distributed load in global coords
	}
}

// The Fx load is represented as a polygon with the following vertices:
//
//  1. The start node, which is at (0, 0) due to the transformation.
//  2. With x = 0, move y by the value of the load's start value.
//  3. With x = endX, move y by the value of the load's end value.
//  4. The end node, which is at (endX, 0) due to the transformation.
//
// By convention, the value of the Fx load is drawn in the Y axis, with the arrows
// pointing in the X axis direction, according to the load's sign.
func drawLocalDistributedFxLoad(
	dLoad *load.DistributedLoad,
	bar *structure.Element,
	ctx *plotContext,
) {
	if dLoad.Term != load.FX {
		panic(fmt.Sprintf("Invalid distributed load term: %s. Expected Fx", dLoad.Term))
	}

	var (
		canvas    = ctx.canvas
		scale     = ctx.unitsScale
		loadScale = ctx.options.DistLoadScale

		endX   = int(scale.applyToLength(bar.Length()))
		startY = int(dLoad.StartValue * loadScale)
		endY   = int(dLoad.EndValue * loadScale)

		x = []int{0, 0, endX, endX}
		y = []int{0, startY, endY, 0}
	)

	canvas.Polygon(x, y)
}

// The Fy load is represented as a polygon with the following vertices:
//
//  1. The start node, which is at (0, 0) due to the transformation.
//  2. With x = 0, move y by the negative value of the load's start value.
//  3. With x = endX, move y by the negative value of the load's end value.
//  4. The end node, which is at (endX, 0) due to the transformation.
//
// By convention, Fy loads are drawn with the arrows pointing towards the bar.
// Thus, the polygon is drawn in the opposite direction of the load, for which
// we need to invert the sign of the load values.
func drawLocalDistributedFyLoad(
	dLoad *load.DistributedLoad,
	bar *structure.Element,
	ctx *plotContext,
) {
	if dLoad.Term != load.FY {
		panic(fmt.Sprintf("Invalid distributed load term: %s. Expected Fy", dLoad.Term))
	}

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
