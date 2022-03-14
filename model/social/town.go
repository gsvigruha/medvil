package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/stats"
	"medvil/model/terrain"
	"medvil/model/time"
)

const ConstructionTransportQuantity = 5

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
	Constructions []*building.Construction
	Stats         *stats.Stats
	Transfers     *MoneyTransfers
}

func (town *Town) Init() {
	defaultTransfers := TransferCategories{
		TaxRate:      20,
		TaxThreshold: 100,
		Subsidy:      100,
	}
	town.Transfers = &MoneyTransfers{
		Farm:     defaultTransfers,
		Workshop: defaultTransfers,
		Mine:     defaultTransfers,
	}
	town.Stats = &stats.Stats{}
}

func (town *Town) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	s := &stats.Stats{}
	taxing := (Calendar.Hour == 0 && Calendar.Day == 1 && Calendar.Month == 1)
	town.Marketplace.ElapseTime(Calendar, m)
	s.Add(town.Marketplace.Stats())
	for l := range town.Townhall.Household.People {
		person := town.Townhall.Household.People[l]
		person.ElapseTime(Calendar, m)
	}
	town.Townhall.ElapseTime(Calendar, m)
	town.Townhall.Household.Filter(Calendar, m)
	s.Add(town.Townhall.Household.Stats())
	for k := range town.Farms {
		farm := town.Farms[k]
		for l := range farm.Household.People {
			person := farm.Household.People[l]
			person.ElapseTime(Calendar, m)
		}
		farm.ElapseTime(Calendar, m)
		if taxing {
			town.Transfers.Farm.Transfer(&town.Townhall.Household.Money, &farm.Household.Money)
		}
		farm.Household.Filter(Calendar, m)
		s.Add(farm.Household.Stats())
	}
	for k := range town.Workshops {
		workshop := town.Workshops[k]
		for l := range workshop.Household.People {
			person := workshop.Household.People[l]
			person.ElapseTime(Calendar, m)
		}
		workshop.ElapseTime(Calendar, m)
		if taxing {
			town.Transfers.Workshop.Transfer(&town.Townhall.Household.Money, &workshop.Household.Money)
		}
		workshop.Household.Filter(Calendar, m)
		s.Add(workshop.Household.Stats())
	}
	for k := range town.Mines {
		mine := town.Mines[k]
		for l := range mine.Household.People {
			person := mine.Household.People[l]
			person.ElapseTime(Calendar, m)
		}
		mine.ElapseTime(Calendar, m)
		if taxing {
			town.Transfers.Mine.Transfer(&town.Townhall.Household.Money, &mine.Household.Money)
		}
		mine.Household.Filter(Calendar, m)
		s.Add(mine.Household.Stats())
	}
	var constructions []*building.Construction
	for k := range town.Constructions {
		construction := town.Constructions[k]
		if construction.IsComplete() {
			b := construction.Building
			switch construction.T {
			case building.BuildingTypeMine:
				mine := &Mine{Household: Household{Building: b, Town: town}}
				mine.Household.Resources.VolumeCapacity = b.Plan.Area() * StoragePerArea
				town.Mines = append(town.Mines, mine)
			case building.BuildingTypeWorkshop:
				w := &Workshop{Household: Household{Building: b, Town: town}}
				w.Household.Resources.VolumeCapacity = b.Plan.Area() * StoragePerArea
				town.Workshops = append(town.Workshops, w)
			case building.BuildingTypeFarm:
				f := &Farm{Household: Household{Building: b, Town: town}}
				f.Household.Resources.VolumeCapacity = b.Plan.Area() * StoragePerArea
				town.Farms = append(town.Farms, f)
			case building.BuildingTypeRoad:
				construction.Road.Construction = false
			case building.BuildingTypeCanal:
				f := m.GetField(construction.X, construction.Y)
				f.Construction = false
				f.Terrain.T = terrain.Water
				f.Terrain.Resources.Add(artifacts.GetArtifact("water"), artifacts.InfiniteQuantity)
			}
			if b != nil {
				m.SetBuildingUnits(b, false)
			}
		} else {
			constructions = append(constructions, construction)
		}
	}
	town.Constructions = constructions
	town.Stats = s
}

func (town *Town) CreateRoadConstruction(x, y uint16, r *building.Road, m navigation.IMap) {
	c := &building.Construction{X: x, Y: y, Road: r, Cost: r.T.Cost, T: building.BuildingTypeRoad, Storage: &artifacts.Resources{}}
	c.Storage.Init(StoragePerArea)
	town.Constructions = append(town.Constructions, c)

	roadF := m.GetField(x, y)
	roadF.Allocated = true
	town.AddConstructionTasks(c, roadF, m)
}

func (town *Town) CreateBuildingConstruction(b *building.Building, bt building.BuildingType, m navigation.IMap) {
	c := &building.Construction{X: b.X, Y: b.Y, Building: b, Cost: b.Plan.ConstructionCost(), T: bt, Storage: &artifacts.Resources{}}
	c.Storage.Init((b.Plan.Area() + b.Plan.RoofArea()) * StoragePerArea)
	town.Constructions = append(town.Constructions, c)

	buildingF := m.GetField(b.X, b.Y)
	town.AddConstructionTasks(c, buildingF, m)
}

func (town *Town) CreateInfraConstruction(x, y uint16, it *building.InfraType, m navigation.IMap) {
	c := &building.Construction{X: x, Y: y, Cost: it.Cost, T: it.BT, Storage: &artifacts.Resources{}}
	c.Storage.Init(StoragePerArea)
	town.Constructions = append(town.Constructions, c)

	f := m.GetField(x, y)
	f.Allocated = true
	f.Construction = true
	town.AddConstructionTasks(c, f, m)
}

func (town *Town) AddConstructionTasks(c *building.Construction, buildingF *navigation.Field, m navigation.IMap) {
	var totalTasks uint16 = 0
	for _, a := range c.Cost {
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
	if totalTasks == 0 {
		totalTasks = 1
	}
	c.MaxProgress = totalTasks
	for i := uint16(0); i < totalTasks; i++ {
		town.Townhall.Household.AddTask(&economy.BuildingTask{
			F: buildingF,
			C: c,
		})
	}
}
