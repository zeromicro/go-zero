package mathx

// MaxInt returns the larger one of an and b.
func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

// MinInt returns the smaller one of an and b.
func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}
