package process

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

// PointSolutionValue is a tuple of T and Value.
type PointSolutionValue struct {
	T     nums.TParam
	Value float64
}

func (psv PointSolutionValue) String() string {
	return fmt.Sprintf("%f : %f ", psv.T.Value(), psv.Value)
}
