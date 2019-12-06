package process

import (
	"fmt"

	"github.com/angelsolaorbaiceta/inkgeom"
)

// PointSolutionValue is a tuple of T and Value.
type PointSolutionValue struct {
	T     inkgeom.TParam `json:"t"`
	Value float64        `json:"val"`
}

func (psv PointSolutionValue) String() string {
	return fmt.Sprintf("T = %f : %f ", psv.T.Value(), psv.Value)
}
