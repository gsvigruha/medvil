package social

import (
	"math"
)

func NumBatchesSimple(totalQuantity, transportQuantity int) int {
	return NumBatches(totalQuantity, transportQuantity, transportQuantity)
}

func NumBatches(totalQuantity, minTransportQuantity, maxTransportQuantity int) int {
	if totalQuantity < minTransportQuantity {
		return 0
	}
	return int(math.Ceil(float64(totalQuantity) / float64(maxTransportQuantity)))
}
