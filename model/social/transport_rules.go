package social

import (
	"math/rand"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/navigation"
)

const WaterTransportQuantity = 10
const FoodTransportQuantity = 6
const ProductTransportMaxVolume = 6
const ExchangeTaskMaxVolumePedestrian = 25

func ProductTransportQuantityWithLimit(a *artifacts.Artifact, maxQ uint16) uint16 {
	var q = ProductTransportQuantity(a)
	if q > maxQ {
		return maxQ
	}
	return q
}

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

func WithinDistance(b *building.Building, f *navigation.Field, d uint16) bool {
	return (b.X-f.X)*(b.X-f.X)+(b.Y-f.Y)*(b.Y-f.Y) <= d*d
}
