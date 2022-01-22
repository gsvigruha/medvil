package model

import (
	"medvil/model/building"
	"medvil/model/social"
)

type ReverseReferences struct {
	BuildingToHousehold map[*building.Building]*social.Household
}

func BuildReverseReferences(m *Map) ReverseReferences {
	BuildingToHousehold := make(map[*building.Building]*social.Household)
	for i := range m.Countries {
		country := m.Countries[i]
		for j := range country.Towns {
			town := country.Towns[j]
			BuildingToHousehold[town.Townhall.Household.Building] = &town.Townhall.Household
			for k := range town.Farms {
				BuildingToHousehold[town.Farms[k].Household.Building] = &town.Farms[k].Household
			}
		}
	}
	return ReverseReferences{BuildingToHousehold: BuildingToHousehold}
}
