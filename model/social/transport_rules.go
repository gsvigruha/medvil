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

func GetRandomBuildingXY(b *building.Building, m navigation.IMap, check func(navigation.Field) bool) (uint16, uint16, bool) {
	fields := b.GetBuildingXYs(true)
	if fields == nil {
		return 0, 0, false
	}
	var filteredFields [][2]uint16
	for i := range fields {
		f := fields[i]
		if check(*m.GetField(f[0]-1, f[1])) || check(*m.GetField(f[0], f[1]-1)) || check(*m.GetField(f[0]+1, f[1])) || check(*m.GetField(f[0], f[1]+1)) {
			filteredFields = append(filteredFields, f)
		}
	}
	idx := rand.Intn(len(filteredFields))
	return filteredFields[idx][0], filteredFields[idx][1], true
}
