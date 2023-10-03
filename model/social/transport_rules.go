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

func GetRandomBuildingXY(b *building.Building, m navigation.IMap, check func(navigation.Field) bool) (uint16, uint16, bool) {
	fields := b.GetBuildingXYs(true)
	if fields == nil {
		return 0, 0, false
	}
	var filteredFields [][2]uint16 = nil
	for i := range fields {
		f := fields[i]
		if check(*m.GetField(f[0], f[1])) {
			filteredFields = append(filteredFields, f)
		}
	}
	if filteredFields == nil {
		return 0, 0, false
	}
	idx := rand.Intn(len(filteredFields))
	return filteredFields[idx][0], filteredFields[idx][1], true
}

func WithinDistance(b *building.Building, f *navigation.Field, d uint16) bool {
	return WithinDistanceCoords(b.X, b.Y, f.X, f.Y, d)
}

func WithinDistanceCoords(x1, y1, x2, y2, d uint16) bool {
	return (int(x1)-int(x2))*(int(x1)-int(x2))+(int(y1)-int(y2))*(int(y1)-int(y2)) <= int(d)*int(d)
}
