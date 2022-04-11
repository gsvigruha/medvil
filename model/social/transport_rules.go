package social

import (
	"medvil/model/artifacts"
)

const WaterTransportQuantity = 10
const FoodTransportQuantity = 6
const ProductTransportMaxVolume = 6
const ExchangeTaskMaxVolume = 25

func ProductTransportQuantity(a *artifacts.Artifact) uint16 {
	return ProductTransportMaxVolume / a.V
}

func MinProductTransportQuantity(as []artifacts.Artifacts) uint16 {
	var q = uint16(65535)
	for _, a := range as {
		if ProductTransportQuantity(a.A) < q {
			q = ProductTransportQuantity(a.A)
		}
	}
	return q
}
