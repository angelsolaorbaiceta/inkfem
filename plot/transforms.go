package plot

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkfem/math"
	"github.com/angelsolaorbaiceta/inkfem/structure"
)

// transformMatrix returns a string representation of a 2D transformation matrix.
// The transformation matrix has the following values:
//
//	| sx  0   tx |
//	|  0  sy  ty |
//	|  0  0   1  |
//
// Where:
//   - sx is the scaling factor in the x-axis.
//   - sy is the scaling factor in the y-axis.
//   - tx is the translation factor in the x-axis.
//   - ty is the translation factor in the y-axis.
//
// Note that this function doesn't use shear values, as it's not needed for the
// transformations in the plot.
func transformMatrix(sx, sy, tx, ty float64) string {
	return fmt.Sprintf("matrix(%f,0,0,%f,%f,%f)", sx, sy, tx, ty)
}

// rotation returns a string representation of a rotation transformation.
func rotation(angleInRad float64) string {
	angleInDeg := math.RadToDeg(angleInRad)
	return fmt.Sprintf("rotate(%f)", angleInDeg)
}

// translate returns a string representation of a translation transformation.
func translate(tx, ty float64) string {
	return fmt.Sprintf("translate(%f %f)", tx, ty)
}

// scale returns a string representation of a scaling transformation.
func scale(sx, sy float64) string {
	return fmt.Sprintf("scale(%f %f)", sx, sy)
}

// transformToLocalBar returns a string representation of the transformation
// needed to set the reference frame of the canvas to the reference frame of the bar.
// The origin is set in the bar's start point, and the X axis is aligned with the
// bar's direction.
func transformToLocalBar(bar *structure.Element, scale unitsScale) string {
	var (
		angle     = bar.RefFrame().AngleInRadsFromX()
		origin    = scale.applyToPoint(bar.StartPoint())
		translate = translate(origin.X(), origin.Y())
		rotation  = rotation(angle)
	)

	return fmt.Sprintf("%s %s", translate, rotation)
}
