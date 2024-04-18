package mathx

type Numerical interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// AtLeast returns the greater of x or lower.
func AtLeast[T Numerical](x, lower T) T {
	if x < lower {
		return lower
	}
	return x
}

// AtMost returns the smaller of x or upper.
func AtMost[T Numerical](x, upper T) T {
	if x > upper {
		return upper
	}
	return x
}

// Between returns the value of x clamped to the range [lower, upper].
func Between[T Numerical](x, lower, upper T) T {
	if x < lower {
		return lower
	}
	if x > upper {
		return upper
	}
	return x
}
