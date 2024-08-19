package math

import "github.com/angelsolaorbaiceta/inkgeom/nums"

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
func MakeIntCloseInterval(nums []int) IntCloseInterval {
	min, max := nums[0], nums[0]

	for _, n := range nums[1:] {
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
// The value is calculated as min + t*(max - min), rounded to the nearest integer.
func (i IntCloseInterval) ValueAt(t nums.TParam) int {
	return i.min + int(t.Value()*float64(i.max-i.min))
}
