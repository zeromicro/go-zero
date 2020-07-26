package internal

import "strings"

const keySeparator = ","

func TotalWeights(c []NodeConf) int {
	var weights int

	for _, node := range c {
		if node.Weight < 0 {
			node.Weight = 0
		}
		weights += node.Weight
	}

	return weights
}

func formatKeys(keys []string) string {
	return strings.Join(keys, keySeparator)
}
