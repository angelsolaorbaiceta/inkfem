package process

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

// PointSolutionValue is a tuple of a T parameter and a value.
type PointSolutionValue struct {
	T     nums.TParam
	Value float64
}

func (psv PointSolutionValue) String() string {
	return fmt.Sprintf("%f : %f", psv.T.Value(), psv.Value)
}
