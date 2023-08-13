package maps

import (
	"medvil/model"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/social"
)

type CountryConf struct {
	TownhallPlan    string
	MarketplacePlan string
	Res             map[string]uint16
	People          uint16
}

var PlayerConf = CountryConf{
	TownhallPlan:    "samples/building/townhouse_1.building.json",
	MarketplacePlan: "samples/building/marketplace_1.building.json",
	Res: map[string]uint16{
		"fruit":     50,
		"vegetable": 50,
		"bread":     20,
		"cube":      50,
		"brick":     50,
		"board":     40,
		"tile":      20,
		"thatch":    10,
		"log":       20,
		"textile":   30,
	},
	People: 5,
}

var OutlawConf = CountryConf{
	TownhallPlan:    "samples/building/outlaw_townhouse_1.building.json",
	MarketplacePlan: "samples/building/outlaw_marketplace_1.building.json",
	Res: map[string]uint16{
		"fruit":     20,
		"vegetable": 20,
		"bread":     20,
		"log":       10,
	},
	People: 3,
}

func GenerateCountry(conf CountryConf, m *model.Map) {
	tx, ty := findStartingLocation(m)

	townhall := &building.Building{
		Plan: building.BuildingPlanFromJSON(conf.TownhallPlan),
		X:    uint16(tx - 2),
		Y:    uint16(ty),
	}
	AddBuilding(townhall, m)
	marketplace := &building.Building{
		Plan: building.BuildingPlanFromJSON(conf.MarketplacePlan),
		X:    uint16(tx + 2),
		Y:    uint16(ty),
	}
	AddBuilding(marketplace, m)

	country := &social.Country{Towns: []*social.Town{&social.Town{}}}
	m.Countries = append(m.Countries, country)
	town := country.Towns[0]
	town.Country = country
	town.Townhall = &social.Townhall{Household: &social.Household{Building: townhall, Town: town}}
	town.Marketplace = &social.Marketplace{Building: marketplace, Town: town}
	town.Townhall.Household.People = make([]*social.Person, conf.People)
	town.Townhall.Household.TargetNumPeople = conf.People
	town.Townhall.Household.Resources.VolumeCapacity = town.Townhall.Household.Building.Plan.Area() * social.StoragePerArea
	town.Townhall.Household.Building.Plan.BuildingType = building.BuildingTypeTownhall
	town.Marketplace.Building.Plan.BuildingType = building.BuildingTypeMarket
	town.Townhall.Household.Money = 2000
	town.Marketplace.Money = 2000
	for i := range town.Townhall.Household.People {
		town.Townhall.Household.People[i] = town.Townhall.Household.NewPerson(m)
	}
	{
		res := &town.Townhall.Household.Resources
		for a, q := range conf.Res {
			res.Add(artifacts.GetArtifact(a), q)
		}
		town.Init()
	}
	{
		town.Marketplace.Init()
		res := &town.Marketplace.Storage
		res.Add(artifacts.GetArtifact("vegetable"), 50)
		res.Add(artifacts.GetArtifact("bread"), 20)
		res.Add(artifacts.GetArtifact("log"), 20)
		res.Add(artifacts.GetArtifact("textile"), 30)
	}
}
