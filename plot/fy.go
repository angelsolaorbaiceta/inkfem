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
)

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

		// Draw the line if there is enough space to draw the arrow.
		if math.AbsInt(loadY) > ctx.config.DistLoadArrowSize {
			canvas.Line(
				loadX, loadY, loadX, 0,
				fmt.Sprintf("marker-end=\"url(#%s)\"", loadArrowMarkerId),
				fmt.Sprintf("stroke=\"%s\"", ctx.config.DistLoadColor),
			)
		}
	}
}
