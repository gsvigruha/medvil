package social

import (
	"math"
)

func NumBatchesSimple(totalQuantity, transportQuantity uint16) int {
	return NumBatches(totalQuantity, transportQuantity, transportQuantity)
}

func NumBatches(totalQuantity, minTransportQuantity, maxTransportQuantity uint16) int {
	if totalQuantity < minTransportQuantity {
		return 0
	}
	return int(math.Ceil(float64(totalQuantity) / float64(maxTransportQuantity)))
}
