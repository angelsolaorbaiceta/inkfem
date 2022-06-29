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

func (psv PointSolutionValue) Equals(other PointSolutionValue, epsilon float64) bool {
	return psv.T.Equals(other.T) && nums.FloatsEqualEps(psv.Value, other.Value, epsilon)
}

// isSameAsLast returns true if the last element in the slice is the same as the value
// passed as second argument.
//
// A value is the same if it is at the same T position and has the same value compared
// using the given epsilon.
func isSameAsLast(values []PointSolutionValue, value PointSolutionValue, epsilon float64) bool {
	if len(values) == 0 {
		return false
	}

	return values[len(values)-1].Equals(value, epsilon)
}

// appendIfNotSameAsLast adds the value to the slice if it's not the same as the last element.
//
// A value is the same if it is at the same T position and has the same value compared
// using the given epsilon.
func appendIfNotSameAsLast(values []PointSolutionValue, value PointSolutionValue, epsilon float64) []PointSolutionValue {
	if !isSameAsLast(values, value, epsilon) {
		values = append(values, value)
	}

	return values
}
