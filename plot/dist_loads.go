package plot

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

var (
	// fyDistLoadLinePositions are the positions where the distributed Fy load lines
	// with the arrowheads are drawn. We draw them at 10% intervals of the length.
	fyDistLoadLinePositions = []nums.TParam{
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
	// fxDistLoadLinePositions are the positions where the distributed Fx load lines
	// with the arrowheads are drawn. We draw them at 20% intervals of the length.
	fxDistLoadLinePositions = []nums.TParam{
		nums.MakeTParam(0.2),
		nums.MakeTParam(0.4),
		nums.MakeTParam(0.6),
		nums.MakeTParam(0.8),
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

	// To draw the arrow lines, we first need to determine the Y interval between
	// where the arrow lines are drawn. This is the two Y limit coordinates of
	// the polygon.
	var (
		interval             = math.MakeIntCloseInterval([]int{0, startY, endY})
		lineYPos             int
		lineXStart, lineXEnd int
	)
	for _, t := range fxDistLoadLinePositions {
		// The Y coordinate where the line is drawn
		lineYPos = interval.ValueAt(t)

		// The X coordinate where the line starts.
		// If the Y position is between the load's start value and 0, the line starts
		// at the start node.
		if math.IntIsBetweenCloseRange(lineYPos, startY, 0) {
			lineXStart = 0
		} else {
			lineXStart = endX
		}

		// The X coordinate where the line ends. It has to be either the end of the
		// bar (when the line ends at the load's end node) or a point in the load's
		// polygon.
		lineXEnd = endX

		if lineXEnd-lineXStart > ctx.config.DistLoadArrowSize {
			canvas.Line(
				lineXStart, lineYPos, lineXEnd, lineYPos,
				fmt.Sprintf("marker-end=\"url(#%s)\"", loadArrowMarkerId),
				fmt.Sprintf("stroke=\"%s\"", ctx.config.DistLoadColor),
			)
		}
	}
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

		startX = int(scale.applyToLength(bar.Length() * dLoad.StartT.Value()))
		endX   = int(scale.applyToLength(bar.Length() * dLoad.EndT.Value()))
		startY = int(-dLoad.StartValue * loadScale)
		endY   = int(-dLoad.EndValue * loadScale)

		x = []int{startX, startX, endX, endX}
		y = []int{0, startY, endY, 0}
	)

	canvas.Polygon(x, y)
	canvas.Text(
		0, 0,
		fmt.Sprintf("%.2f", dLoad.StartValue),
		textTransform(startX, startY), fmt.Sprintf("fill:%s", ctx.config.DistLoadColor),
	)
	canvas.Text(
		0, 0,
		fmt.Sprintf("%.2f", dLoad.EndValue),
		textTransform(endX, endY), fmt.Sprintf("fill:%s", ctx.config.DistLoadColor),
	)

	for _, t := range fyDistLoadLinePositions {
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
