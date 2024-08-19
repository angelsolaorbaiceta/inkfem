package math

// AbsInt returns the absolute value of an integer.
func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// IntMax returns the maximum of a list of integers.
func IntMax(nums ...int) int {
	max := nums[0]

	for _, n := range nums[1:] {
		if n > max {
			max = n
		}
	}

	return max
}

// IntMin returns the minimum of a list of integers.
func IntMin(nums ...int) int {
	min := nums[0]

	for _, n := range nums[1:] {
		if n < min {
			min = n
		}
	}

	return min
}

// IntIsBetweenCloseRange returns true if x is between min(a, b) and max(a, b),
// inclusive. The order of a and b doesn't matter.
func IntIsBetweenCloseRange(x, a, b int) bool {
	var (
		min = IntMin(a, b)
		max = IntMax(a, b)
	)

	return x >= min && x <= max
}
