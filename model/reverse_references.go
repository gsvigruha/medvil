package model

import (
	"medvil/model/building"
	"medvil/model/social"
)

type ReverseReferences struct {
	BuildingToFarm        map[*building.Building]*social.Farm
	BuildingToWorkshop    map[*building.Building]*social.Workshop
	BuildingToTownhall    map[*building.Building]*social.Townhall
	BuildingToMarketplace map[*building.Building]*social.Marketplace
}

func BuildReverseReferences(m *Map) ReverseReferences {
	BuildingToFarm := make(map[*building.Building]*social.Farm)
	BuildingToWorkshop := make(map[*building.Building]*social.Workshop)
	BuildingToTownhall := make(map[*building.Building]*social.Townhall)
	BuildingToMarketplace := make(map[*building.Building]*social.Marketplace)
	for i := range m.Countries {
		country := m.Countries[i]
		for j := range country.Towns {
			town := country.Towns[j]
			BuildingToTownhall[town.Townhall.Household.Building] = town.Townhall
			BuildingToMarketplace[town.Marketplace.Building] = town.Marketplace
			for k := range town.Farms {
				BuildingToFarm[town.Farms[k].Household.Building] = town.Farms[k]
			}
			for k := range town.Workshops {
				BuildingToWorkshop[town.Workshops[k].Household.Building] = town.Workshops[k]
			}
		}
	}
	return ReverseReferences{
		BuildingToFarm:        BuildingToFarm,
		BuildingToWorkshop:    BuildingToWorkshop,
		BuildingToTownhall:    BuildingToTownhall,
		BuildingToMarketplace: BuildingToMarketplace,
	}
}
