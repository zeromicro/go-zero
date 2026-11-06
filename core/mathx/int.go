package mathx

// MaxInt returns the larger one of a and b.
// Deprecated: use builtin max instead.
func MaxInt(a, b int) int {
	return max(a, b)
}

// MinInt returns the smaller one of a and b.
// Deprecated: use builtin min instead.
func MinInt(a, b int) int {
	return min(a, b)
}
