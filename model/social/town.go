package social

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

const FarmTaxRate float64 = 0.2
const WorkshopTaxRate float64 = 0.1
const MineTaxRate float64 = 0.1

type JSONBuilding struct {
	Plan string
	X    uint16
	Y    uint16
}

type JSONFarm struct {
	Land       [][]uint16
	Building   JSONBuilding
	Population uint8
	Money      uint32
}

type Town struct {
	Country     *Country
	Townhall    *Townhall
	Marketplace *Marketplace
	Farms       []*Farm
	Workshops   []*Workshop
	Mines       []*Mine
}

func (town *Town) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	taxing := (Calendar.Hour == 0 && Calendar.Day == 1 && Calendar.Month == 1)
	town.Marketplace.ElapseTime(Calendar, m)
	for l := range town.Townhall.Household.People {
		person := town.Townhall.Household.People[l]
		person.ElapseTime(Calendar, m)
	}
	town.Townhall.ElapseTime(Calendar, m)
	town.Townhall.Household.Filter(Calendar, m)
	for k := range town.Farms {
		farm := town.Farms[k]
		for l := range farm.Household.People {
			person := farm.Household.People[l]
			person.ElapseTime(Calendar, m)
		}
		farm.ElapseTime(Calendar, m)
		if taxing && farm.Household.Money > 0 {
			tax := uint32(FarmTaxRate * float64(farm.Household.Money))
			farm.Household.Money -= tax
			town.Townhall.Household.Money += tax
		}
		farm.Household.Filter(Calendar, m)
	}
	for k := range town.Workshops {
		workshop := town.Workshops[k]
		for l := range workshop.Household.People {
			person := workshop.Household.People[l]
			person.ElapseTime(Calendar, m)
		}
		workshop.ElapseTime(Calendar, m)
		if taxing && workshop.Household.Money > 0 {
			tax := uint32(WorkshopTaxRate * float64(workshop.Household.Money))
			workshop.Household.Money -= tax
			town.Townhall.Household.Money += tax
		}
		workshop.Household.Filter(Calendar, m)
	}
	for k := range town.Mines {
		mine := town.Mines[k]
		for l := range mine.Household.People {
			person := mine.Household.People[l]
			person.ElapseTime(Calendar, m)
		}
		mine.ElapseTime(Calendar, m)
		if taxing && mine.Household.Money > 0 {
			tax := uint32(MineTaxRate * float64(mine.Household.Money))
			mine.Household.Money -= tax
			town.Townhall.Household.Money += tax
		}
		mine.Household.Filter(Calendar, m)
	}
}
