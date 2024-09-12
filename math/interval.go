package math

import (
	"math"

	"github.com/angelsolaorbaiceta/inkgeom/nums"
)

// IntCloseInterval represents a closed interval of integers.
type IntCloseInterval struct {
	min int
	max int
}

// Min returns the left endpoint of the interval.
func (i IntCloseInterval) Min() int {
	return i.min
}

// Max returns the right endpoint of the interval.
func (i IntCloseInterval) Max() int {
	return i.max
}

// MakeIntCloseInterval creates a new IntCloseInterval from a list of integers.
func MakeIntCloseInterval(numbers []int) IntCloseInterval {
	min, max := numbers[0], numbers[0]

	for _, n := range numbers[1:] {
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}

	return IntCloseInterval{min: min, max: max}
}

// ValueAt returns the value at a given parameter t in the interval.
func (i IntCloseInterval) ValueAt(t nums.TParam) int {
	steps := t.Value() * float64(i.max-i.min) / nums.MaxT.Value()
	return i.min + int(math.Round(steps))
}
