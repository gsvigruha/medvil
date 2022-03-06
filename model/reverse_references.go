package model

import (
	"medvil/model/building"
	"medvil/model/social"
)

type ReverseReferences struct {
	BuildingToFarm         map[*building.Building]*social.Farm
	BuildingToMine         map[*building.Building]*social.Mine
	BuildingToWorkshop     map[*building.Building]*social.Workshop
	BuildingToTownhall     map[*building.Building]*social.Townhall
	BuildingToMarketplace  map[*building.Building]*social.Marketplace
	BuildingToConstruction map[*building.Building]*building.Construction
}

func BuildReverseReferences(m *Map) ReverseReferences {
	BuildingToFarm := make(map[*building.Building]*social.Farm)
	BuildingToMine := make(map[*building.Building]*social.Mine)
	BuildingToWorkshop := make(map[*building.Building]*social.Workshop)
	BuildingToTownhall := make(map[*building.Building]*social.Townhall)
	BuildingToMarketplace := make(map[*building.Building]*social.Marketplace)
	BuildingToConstruction := make(map[*building.Building]*building.Construction)

	for i := range m.Countries {
		country := m.Countries[i]
		for j := range country.Towns {
			town := country.Towns[j]
			BuildingToTownhall[town.Townhall.Household.Building] = town.Townhall
			BuildingToMarketplace[town.Marketplace.Building] = town.Marketplace
			for k := range town.Farms {
				BuildingToFarm[town.Farms[k].Household.Building] = town.Farms[k]
			}
			for k := range town.Mines {
				BuildingToMine[town.Mines[k].Household.Building] = town.Mines[k]
			}
			for k := range town.Workshops {
				BuildingToWorkshop[town.Workshops[k].Household.Building] = town.Workshops[k]
			}
			for k := range town.Constructions {
				if town.Constructions[k].Building != nil {
					BuildingToConstruction[town.Constructions[k].Building] = town.Constructions[k]
				}
			}
		}
	}
	return ReverseReferences{
		BuildingToFarm:         BuildingToFarm,
		BuildingToMine:         BuildingToMine,
		BuildingToWorkshop:     BuildingToWorkshop,
		BuildingToTownhall:     BuildingToTownhall,
		BuildingToMarketplace:  BuildingToMarketplace,
		BuildingToConstruction: BuildingToConstruction,
	}
}
