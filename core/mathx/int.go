package mathx

// MaxInt returns the larger one of a and b.
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinInt returns the smaller one of a and b.
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}
