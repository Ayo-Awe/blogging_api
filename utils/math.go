package utils

func ClampInt(val, min, max int) int {
	if val < min {
		return min
	}

	if val > max {
		return max
	}

	return val
}
