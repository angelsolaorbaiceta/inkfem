package plot

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/structure"
	"github.com/angelsolaorbaiceta/inkfem/structure/load"
)

func drawLocalDistributedMzLoad(
	dLoad *load.DistributedLoad,
	bar *structure.Element,
	ctx *plotContext,
) {
	if dLoad.Term != load.MZ {
		panic("Invalid distributed load term. Expected Mz")
	}

	var (
		canvas    = ctx.canvas
		scale     = ctx.unitsScale
		loadScale = ctx.options.DistLoadScale

		startX = int(scale.applyToLength(bar.Length() * dLoad.StartT.Value()))
		endX   = int(scale.applyToLength(bar.Length() * dLoad.EndT.Value()))
		startY = int(dLoad.StartValue * loadScale)
		endY   = int(dLoad.EndValue * loadScale)

		x = []int{startX, startX, endX, endX}
		y = []int{0, startY, endY, 0}
	)

	canvas.Polygon(x, y)
	canvas.Text(
		0, 0,
		fmt.Sprintf("%.2f", dLoad.StartValue),
		textTransform(startX, startY),
		fmt.Sprintf("fill:%s", ctx.config.DistLoadColor),
	)
	canvas.Text(
		0, 0,
		fmt.Sprintf("%.2f", dLoad.EndValue),
		textTransform(endX, endY),
		fmt.Sprintf("fill:%s", ctx.config.DistLoadColor),
	)
}
