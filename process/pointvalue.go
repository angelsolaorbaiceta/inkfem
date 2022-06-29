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

func (psv PointSolutionValue) Equals(other PointSolutionValue) bool {
	return psv.T.Equals(other.T) && nums.FloatsEqual(psv.Value, other.Value)
}

// isSameAsLast returns true if the last element in the slice is the same as the value
// passed as second argument.
func isSameAsLast(values []PointSolutionValue, value PointSolutionValue) bool {
	if len(values) == 0 {
		return false
	}

	return values[len(values)-1].Equals(value)
}

// appendIfNotSameAsLast adds the value to the slice if it's not the same as the last element.
func appendIfNotSameAsLast(values []PointSolutionValue, value PointSolutionValue) []PointSolutionValue {
	if !isSameAsLast(values, value) {
		values = append(values, value)
	}

	return values
}
