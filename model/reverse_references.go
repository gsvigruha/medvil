package model

import (
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/social"
)

type ReverseReferences struct {
	BuildingToFarm         map[*building.Building]*social.Farm
	BuildingToMine         map[*building.Building]*social.Mine
	BuildingToWorkshop     map[*building.Building]*social.Workshop
	BuildingToFactory      map[*building.Building]*social.Factory
	BuildingToTower        map[*building.Building]*social.Tower
	BuildingToTownhall     map[*building.Building]*social.Townhall
	BuildingToMarketplace  map[*building.Building]*social.Marketplace
	BuildingToConstruction map[*building.Building]*building.Construction
	TravellerToPerson      map[*navigation.Traveller]*social.Person
	TravellerToTrader      map[*navigation.Traveller]*social.Trader
}

func AddPeople(TravellerToPerson map[*navigation.Traveller]*social.Person, h *social.Household) {
	for l := range h.People {
		p := h.People[l]
		TravellerToPerson[p.Traveller] = p
	}
}

func BuildReverseReferences(m *Map) ReverseReferences {
	BuildingToFarm := make(map[*building.Building]*social.Farm)
	BuildingToMine := make(map[*building.Building]*social.Mine)
	BuildingToWorkshop := make(map[*building.Building]*social.Workshop)
	BuildingToFactory := make(map[*building.Building]*social.Factory)
	BuildingToTower := make(map[*building.Building]*social.Tower)
	BuildingToTownhall := make(map[*building.Building]*social.Townhall)
	BuildingToMarketplace := make(map[*building.Building]*social.Marketplace)
	BuildingToConstruction := make(map[*building.Building]*building.Construction)
	TravellerToPerson := make(map[*navigation.Traveller]*social.Person)
	TravellerToTrader := make(map[*navigation.Traveller]*social.Trader)

	for i := range m.Countries {
		country := m.Countries[i]
		for j := range country.Towns {
			town := country.Towns[j]
			BuildingToTownhall[town.Townhall.Household.Building] = town.Townhall
			AddPeople(TravellerToPerson, town.Townhall.Household)
			for k := range town.Townhall.Traders {
				t := town.Townhall.Traders[k]
				TravellerToTrader[t.Person.Traveller] = t
			}
			if town.Marketplace != nil {
				BuildingToMarketplace[town.Marketplace.Building] = town.Marketplace
			}
			for k := range town.Farms {
				BuildingToFarm[town.Farms[k].Household.Building] = town.Farms[k]
				AddPeople(TravellerToPerson, town.Farms[k].Household)
			}
			for k := range town.Mines {
				BuildingToMine[town.Mines[k].Household.Building] = town.Mines[k]
				AddPeople(TravellerToPerson, town.Mines[k].Household)
			}
			for k := range town.Workshops {
				BuildingToWorkshop[town.Workshops[k].Household.Building] = town.Workshops[k]
				AddPeople(TravellerToPerson, town.Workshops[k].Household)
			}
			for k := range town.Factories {
				BuildingToFactory[town.Factories[k].Household.Building] = town.Factories[k]
				AddPeople(TravellerToPerson, town.Factories[k].Household)
			}
			for k := range town.Towers {
				BuildingToTower[town.Towers[k].Household.Building] = town.Towers[k]
				AddPeople(TravellerToPerson, town.Towers[k].Household)
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
		BuildingToFactory:      BuildingToFactory,
		BuildingToTower:        BuildingToTower,
		BuildingToTownhall:     BuildingToTownhall,
		BuildingToMarketplace:  BuildingToMarketplace,
		BuildingToConstruction: BuildingToConstruction,
		TravellerToPerson:      TravellerToPerson,
		TravellerToTrader:      TravellerToTrader,
	}
}
