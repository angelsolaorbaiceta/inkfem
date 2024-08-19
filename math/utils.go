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
