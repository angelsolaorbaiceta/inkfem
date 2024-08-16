package plot

import "fmt"

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
