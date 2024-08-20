package plot

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

var (
	// fxDistLoadLinePositions are the positions where the distributed Fx load lines
	// with the arrowheads are drawn. We draw them at 20% intervals of the length.
	fxDistLoadLinePositions = []nums.TParam{
		nums.MakeTParam(0.2),
		nums.MakeTParam(0.4),
		nums.MakeTParam(0.6),
		nums.MakeTParam(0.8),
	}
)

// The Fx load is represented as a polygon with the following vertices:
//
//  1. The start node, which is at (0, 0) due to the transformation.
//  2. With x = 0, move y by the value of the load's start value.
//  3. With x = endX, move y by the value of the load's end value.
//  4. The end node, which is at (endX, 0) due to the transformation.
//
// By convention, the value of the Fx load is drawn in the Y axis, with the
// arrows pointing in the X axis direction, according to the load's sign.
// Positive values point to the X axis direction, while negative values point
// to the opposite direction.
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

		xLength = scale.applyToLength(bar.Length())
		startX  = int(scale.applyToLength(bar.Length() * dLoad.StartT.Value()))
		endX    = int(scale.applyToLength(bar.Length() * dLoad.EndT.Value()))
		startY  = int(dLoad.StartValue * loadScale)
		endY    = int(dLoad.EndValue * loadScale)

		x = []int{startX, startX, endX, endX}
		y = []int{0, startY, endY, 0}
	)

	canvas.Polygon(x, y)

	// To draw the arrow lines, we first need to determine the Y yInterval between
	// where the arrow lines are drawn. This is the two Y limit coordinates of
	// the polygon.
	var (
		yInterval            = math.MakeIntCloseInterval([]int{0, startY, endY})
		lineYPos             int
		lineXStart, lineXEnd int
	)

	loadEq, err := dLoad.AsEquation(xLength, loadScale)
	if err != nil {
		panic(err)
	}

	for _, t := range fxDistLoadLinePositions {
		// The Y coordinate where the line is drawn
		lineYPos = yInterval.ValueAt(t)

		loadX, err := loadEq.XAt(float64(lineYPos))
		if err != nil {
			continue
		}

		// The X coordinate where the line starts.
		// If the Y position is between the load's start value and 0, the line starts
		// at the load's start T position. Otherwise, it starts at the load's slope.
		if math.IntIsBetweenCloseRange(lineYPos, startY, 0) {
			lineXStart = startX
		} else {
			lineXStart = int(loadX)
		}

		// The X coordinate where the line ends. It has to be either the end of the
		// bar (when the line ends at the load's end node) or a point in the load's
		// polygon.
		if math.IntIsBetweenCloseRange(lineYPos, endY, 0) {
			lineXEnd = endX
		} else {
			lineXEnd = int(loadX)
		}

		// If the value of the load is negative, switch the start and end points.
		if lineYPos < 0 {
			lineXStart, lineXEnd = lineXEnd, lineXStart
		}

		if math.AbsInt(lineXEnd-lineXStart) > ctx.config.DistLoadArrowSize {
			canvas.Line(
				lineXStart, lineYPos, lineXEnd, lineYPos,
				fmt.Sprintf("marker-end=\"url(#%s)\"", loadArrowMarkerId),
				fmt.Sprintf("stroke=\"%s\"", ctx.config.DistLoadColor),
			)
		}
	}
}
