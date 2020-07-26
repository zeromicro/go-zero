package mathx

import "math"

func CalcEntropy(m map[interface{}]int, total int) float64 {
	var entropy float64

	for _, v := range m {
		proba := float64(v) / float64(total)
		entropy -= proba * math.Log2(proba)
	}

	return entropy / math.Log2(float64(len(m)))
}
