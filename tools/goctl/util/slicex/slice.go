package slicex

func Unique[T comparable](s []T) []T {
	if len(s) <= 1 {
		return s
	}
	existedMap := make(map[T]bool)
	writeIndex := 0

	for readIndex := 0; readIndex < len(s); readIndex++ {
		if !existedMap[s[readIndex]] {
			existedMap[s[readIndex]] = true
			s[writeIndex] = s[readIndex]
			writeIndex++
		}
	}
	return s[:writeIndex]
}
