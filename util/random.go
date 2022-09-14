package util

import (
	"math/rand"
)

func RandomIndexWeighted(weights []float64) int {
	var sum = 0.0
	for _, w := range weights {
		sum += w
	}
	cutoff := rand.Float64() * sum
	var cnt = 0.0
	for i, w := range weights {
		cnt += w
		if cnt >= cutoff {
			return i
		}
	}
	// Should not happen
	return -1
}
