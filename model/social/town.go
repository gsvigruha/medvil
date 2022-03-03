package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
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
	Country       *Country
	Townhall      *Townhall
	Marketplace   *Marketplace
	Farms         []*Farm
	Workshops     []*Workshop
	Mines         []*Mine
	Constructions []*building.BuildingConstruction
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
	var constructions []*building.BuildingConstruction
	for k := range town.Constructions {
		construction := town.Constructions[k]
		if construction.IsComplete() {
			switch construction.T {
			case building.BuildingTypeMine:
				mine := &Mine{Household: Household{Building: construction.Building, Town: town}}
				mine.Household.Resources.VolumeCapacity = mine.Household.Building.Plan.Area() * StoragePerArea
				town.Mines = append(town.Mines, mine)
				m.SetBuildingUnits(construction.Building, false)
			}
		} else {
			constructions = append(constructions, construction)
		}
	}
	town.Constructions = constructions
}

const ConstructionTransportQuantity = 5

func (town *Town) CreateConstruction(b *building.Building, bt building.BuildingType, m navigation.IMap) {
	c := &building.BuildingConstruction{Building: b, RemainingCost: b.Plan.ConstructionCost(), T: bt, Storage: &artifacts.Resources{}}
	c.Storage.Init()
	town.Constructions = append(town.Constructions, c)

	buildingF := m.GetField(c.Building.X, c.Building.Y)
	var totalTasks uint16 = 0
	for _, a := range c.RemainingCost {
		var totalQ = a.Quantity
		totalTasks += totalQ
		for totalQ > 0 {
			var q uint16 = ConstructionTransportQuantity
			if totalQ < ConstructionTransportQuantity {
				q = totalQ
			}
			totalQ -= q
			town.Townhall.Household.AddTask(&economy.TransportTask{
				PickupF:  m.GetField(town.Townhall.Household.Building.X, town.Townhall.Household.Building.Y),
				DropoffF: buildingF,
				PickupR:  &town.Townhall.Household.Resources,
				DropoffR: c.Storage,
				A:        a.A,
				Quantity: q,
			})
		}
	}
	c.MaxProgress = totalTasks
	for i := uint16(0); i < totalTasks; i++ {
		town.Townhall.Household.AddTask(&economy.BuildingTask{
			F: buildingF,
			C: c,
		})
	}
}
