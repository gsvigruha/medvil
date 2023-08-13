package maps

import (
	"medvil/model"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/social"
)

type CountryConf struct {
	TownhallPlan    string
	MarketplacePlan string
	FarmPlan        string
	WorkshopPlan    string
	TownhallRes     map[string]uint16
	MarketplaceRes  map[string]uint16
	People          uint16
	Money           uint32
	Village         bool
}

var PlayerConf = CountryConf{
	TownhallPlan:    "samples/building/townhouse_1.building.json",
	MarketplacePlan: "samples/building/marketplace_1.building.json",
	TownhallRes: map[string]uint16{
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
	MarketplaceRes: map[string]uint16{
		"vegetable": 50,
		"bread":     20,
		"log":       20,
		"textile":   30,
	},
	People:  5,
	Money:   2000,
	Village: false,
}

var OutlawConf = CountryConf{
	TownhallPlan:    "samples/building/outlaw_townhouse_1.building.json",
	MarketplacePlan: "samples/building/outlaw_marketplace_1.building.json",
	FarmPlan:        "samples/building/outlaw_farm_1.building.json",
	WorkshopPlan:    "samples/building/outlaw_workshop_1.building.json",
	TownhallRes: map[string]uint16{
		"fruit":     20,
		"vegetable": 20,
		"bread":     20,
		"log":       10,
	},
	MarketplaceRes: map[string]uint16{
		"vegetable": 50,
		"sheep":     5,
		"log":       20,
		"textile":   30,
	},
	People:  8,
	Money:   1000,
	Village: true,
}

func addFarm(conf CountryConf, town *social.Town, x, y int, m *model.Map) {
	farmB := &building.Building{
		Plan: building.BuildingPlanFromJSON(conf.FarmPlan),
		X:    uint16(x),
		Y:    uint16(y),
	}
	farmB.Plan.BuildingType = building.BuildingTypeFarm
	AddBuilding(farmB, m)
	farm := &social.Farm{Household: &social.Household{Building: farmB, Town: town}}
	farm.Household.TargetNumPeople = 4
	farm.Household.Resources.VolumeCapacity = farm.Household.Building.Plan.Area() * social.StoragePerArea
	addFarmLand(farm, economy.FarmFieldUseTypePasture, -1, 0, m)
	addFarmLand(farm, economy.FarmFieldUseTypePasture, -1, 1, m)
	addFarmLand(farm, economy.FarmFieldUseTypePasture, -1, -1, m)
	addFarmLand(farm, economy.FarmFieldUseTypePasture, -2, 1, m)
	addFarmLand(farm, economy.FarmFieldUseTypePasture, -2, 0, m)
	addFarmLand(farm, economy.FarmFieldUseTypePasture, -2, -1, m)
	addFarmLand(farm, economy.FarmFieldUseTypeVegetables, 1, 0, m)
	addFarmLand(farm, economy.FarmFieldUseTypeVegetables, 1, 1, m)
	addFarmLand(farm, economy.FarmFieldUseTypeVegetables, 1, -1, m)
	addFarmLand(farm, economy.FarmFieldUseTypeVegetables, 0, 1, m)
	addFarmLand(farm, economy.FarmFieldUseTypeForestry, 1, -2, m)
	addFarmLand(farm, economy.FarmFieldUseTypeForestry, 0, -2, m)
	addFarmLand(farm, economy.FarmFieldUseTypeOrchard, -1, -2, m)
	addFarmLand(farm, economy.FarmFieldUseTypeOrchard, -2, -2, m)
	town.Farms = append(town.Farms, farm)
}

func addFarmLand(farm *social.Farm, useType uint8, dx, dy int, m *model.Map) {
	x := uint16(int(farm.Household.Building.X) + dx)
	y := uint16(int(farm.Household.Building.Y) + dy)
	farm.Land = append(farm.Land,
		social.FarmLand{
			X:       x,
			Y:       y,
			UseType: useType,
			F:       m.GetField(x, y),
		},
	)
}

func GenerateCountry(conf CountryConf, m *model.Map) bool {
	tx, ty := findStartingLocation(m)
	if tx == 0 && ty == 0 {
		return false
	}

	townhall := &building.Building{
		Plan: building.BuildingPlanFromJSON(conf.TownhallPlan),
		X:    uint16(tx - 2),
		Y:    uint16(ty),
	}
	townhall.Plan.BuildingType = building.BuildingTypeTownhall

	AddBuilding(townhall, m)
	marketplace := &building.Building{
		Plan: building.BuildingPlanFromJSON(conf.MarketplacePlan),
		X:    uint16(tx + 2),
		Y:    uint16(ty),
	}
	marketplace.Plan.BuildingType = building.BuildingTypeMarket
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
	town.Townhall.Household.Money = conf.Money
	town.Marketplace.Money = conf.Money
	for i := range town.Townhall.Household.People {
		town.Townhall.Household.People[i] = town.Townhall.Household.NewPerson(m)
	}
	{
		res := &town.Townhall.Household.Resources
		for a, q := range conf.TownhallRes {
			res.Add(artifacts.GetArtifact(a), q)
		}
		town.Init()
	}
	{
		town.Marketplace.Init()
		res := &town.Marketplace.Storage
		for a, q := range conf.MarketplaceRes {
			res.Add(artifacts.GetArtifact(a), q)
		}
	}

	if conf.Village {
		addFarm(conf, town, tx+2, ty-4, m)
		addFarm(conf, town, tx+2, ty+4, m)
		{
			workshopB := &building.Building{
				Plan: building.BuildingPlanFromJSON(conf.WorkshopPlan),
				X:    uint16(tx - 1),
				Y:    uint16(ty - 4),
			}
			workshopB.Plan.BaseShape[2][2].Extension = &building.BuildingExtension{T: building.Workshop}
			workshopB.Plan.BuildingType = building.BuildingTypeWorkshop
			AddBuilding(workshopB, m)
			workshop := &social.Workshop{Household: &social.Household{Building: workshopB, Town: town}}
			workshop.Household.TargetNumPeople = 2
			workshop.Household.Resources.VolumeCapacity = workshop.Household.Building.Plan.Area() * social.StoragePerArea
			workshop.Manufacture = economy.GetManufacture("butchering")
			town.Workshops = append(town.Workshops, workshop)
		}
		{
			workshopB := &building.Building{
				Plan: building.BuildingPlanFromJSON(conf.WorkshopPlan),
				X:    uint16(tx - 1),
				Y:    uint16(ty - 5),
			}
			workshopB.Plan.BaseShape[2][2].Extension = &building.BuildingExtension{T: building.Workshop}
			workshopB.Plan.BuildingType = building.BuildingTypeWorkshop
			AddBuilding(workshopB, m)
			workshop := &social.Workshop{Household: &social.Household{Building: workshopB, Town: town}}
			workshop.Household.TargetNumPeople = 2
			workshop.Household.Resources.VolumeCapacity = workshop.Household.Building.Plan.Area() * social.StoragePerArea
			workshop.Manufacture = economy.GetManufacture("sewing")
			town.Workshops = append(town.Workshops, workshop)
		}

		town.Townhall.Household.TargetNumPeople = 2
	}

	return true
}
