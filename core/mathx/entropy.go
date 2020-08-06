package mathx

import "math"

func CalcEntropy(m map[interface{}]int) float64 {
	var entropy float64

	var total int
	for _, v := range m {
		total += v
	}

	for _, v := range m {
		proba := float64(v) / float64(total)
		entropy -= proba * math.Log2(proba)
	}

	return entropy / math.Log2(float64(len(m)))
}
