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
	Factories     []*Factory
	Towers        []*Tower
	Constructions []*building.Construction
	Stats         *stats.Stats
	Transfers     *MoneyTransfers
	Roads         []*navigation.Field
}

func (town *Town) Init() {
	defaultTransfers := TransferCategories{
		TaxRate:      30,
		TaxThreshold: 300,
		Subsidy:      200,
	}
	militaryTransfers := TransferCategories{
		TaxRate:      0,
		TaxThreshold: 0,
		Subsidy:      100,
	}
	town.Transfers = &MoneyTransfers{
		Farm:              defaultTransfers,
		Workshop:          defaultTransfers,
		Mine:              defaultTransfers,
		Factory:           defaultTransfers,
		Tower:             militaryTransfers,
		MarketFundingRate: 70,
	}
	town.Stats = &stats.Stats{}

	town.Townhall.StorageTarget = make(map[*artifacts.Artifact]*int)
	for _, a := range artifacts.All {
		var amount int = 0
		if q, ok := town.Townhall.Household.Resources.Artifacts[a]; ok {
			amount = int(q)
		}
		town.Townhall.StorageTarget[a] = &amount
	}
}

func (town *Town) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	s := &stats.Stats{}
	eoYear := (Calendar.Hour == 0 && Calendar.Day == 1 && Calendar.Month == 1)
	eoMonth := (Calendar.Hour == 0 && Calendar.Day == 1)
	town.Marketplace.ElapseTime(Calendar, m)
	s.Add(town.Marketplace.Stats())
	if eoMonth {
		town.Transfers.FundMarket(&town.Townhall.Household.Money, &town.Marketplace.Money)
	}
	for l := range town.Townhall.Household.People {
		person := town.Townhall.Household.People[l]
		person.ElapseTime(Calendar, m)
	}
	town.Townhall.ElapseTime(Calendar, m)
	town.Townhall.Household.Filter(Calendar, m)
	town.Townhall.Filter(Calendar, m)
	s.Add(town.Townhall.Household.Stats())
	for k := range town.Farms {
		farm := town.Farms[k]
		for l := range farm.Household.People {
			person := farm.Household.People[l]
			person.ElapseTime(Calendar, m)
		}
		farm.ElapseTime(Calendar, m)
		if eoYear {
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
		if eoYear {
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
		if eoYear {
			town.Transfers.Mine.Transfer(&town.Townhall.Household.Money, &mine.Household.Money)
		}
		mine.Household.Filter(Calendar, m)
		s.Add(mine.Household.Stats())
	}
	for k := range town.Factories {
		factory := town.Factories[k]
		for l := range factory.Household.People {
			person := factory.Household.People[l]
			person.ElapseTime(Calendar, m)
		}
		factory.ElapseTime(Calendar, m)
		if eoYear {
			town.Transfers.Factory.Transfer(&town.Townhall.Household.Money, &factory.Household.Money)
		}
		factory.Household.Filter(Calendar, m)
		s.Add(factory.Household.Stats())
	}
	for k := range town.Towers {
		tower := town.Towers[k]
		for l := range tower.Household.People {
			person := tower.Household.People[l]
			person.ElapseTime(Calendar, m)
		}
		tower.ElapseTime(Calendar, m)
		if eoYear {
			town.Transfers.Tower.Transfer(&town.Townhall.Household.Money, &tower.Household.Money)
		}
		tower.Household.Filter(Calendar, m)
		s.Add(tower.Household.Stats())
	}
	var constructions []*building.Construction
	for k := range town.Constructions {
		construction := town.Constructions[k]
		if construction.IsComplete() {
			b := construction.Building
			field := m.GetField(construction.X, construction.Y)
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
			case building.BuildingTypeFactory:
				f := &Factory{Household: Household{Building: b, Town: town}}
				f.Household.Resources.VolumeCapacity = b.Plan.Area() * StoragePerArea
				town.Factories = append(town.Factories, f)
			case building.BuildingTypeTower:
				t := &Tower{Household: Household{Building: b, Town: town}}
				t.Household.Resources.VolumeCapacity = b.Plan.Area() * StoragePerArea
				town.Towers = append(town.Towers, t)
			case building.BuildingTypeRoad:
				if construction.Road.Construction {
					construction.Road.Construction = false
					navigation.SetRoadConnections(m, field)
					town.Roads = append(town.Roads, field)
				} else if construction.Road.Broken {
					construction.Road.Broken = false
				}
			case building.BuildingTypeCanal:
				field.Construction = false
				field.Terrain.T = terrain.Canal
			}
			if b != nil {
				m.SetBuildingUnits(b, false)
				for _, coords := range b.GetBuildingXYs(false) {
					bf := m.GetField(coords[0], coords[1])
					navigation.SetRoadConnectionsForNeighbors(m, bf)
					navigation.SetBuildingDeckForNeighbors(m, bf)
					navigation.SetWallConnections(m, bf)
				}
			}
		} else {
			constructions = append(constructions, construction)
		}
	}
	town.Constructions = constructions
	town.Stats = s
	for _, road := range town.Roads {
		if road.Road.Broken && town.Townhall.Household.NumTasks("building", economy.BuildingTaskTag(road)) == 0 {
			c := &building.Construction{
				Road:    road.Road,
				X:       road.X,
				Y:       road.Y,
				Cost:    road.Road.T.Cost,
				T:       building.BuildingTypeRoad,
				Storage: &artifacts.Resources{},
			}
			town.Constructions = append(town.Constructions, c)
			town.AddConstructionTasks(c, road, m)
		}
	}
}

func (town *Town) CreateRoadConstruction(x, y uint16, r *building.Road, m navigation.IMap) {
	c := &building.Construction{X: x, Y: y, Road: r, Cost: r.T.Cost, T: building.BuildingTypeRoad, Storage: &artifacts.Resources{}}
	c.Storage.Init(StoragePerArea)
	town.Constructions = append(town.Constructions, c)

	roadF := m.GetField(x, y)
	roadF.Allocated = true
	town.AddConstructionTasks(c, roadF, m)
}

func (town *Town) CreateBuildingConstruction(b *building.Building, m navigation.IMap) {
	bt := b.Plan.BuildingType
	c := &building.Construction{X: b.X, Y: b.Y, Building: b, Cost: b.Plan.ConstructionCost(), T: bt, Storage: &artifacts.Resources{}}
	c.Storage.Init((b.Plan.Area() + b.Plan.RoofArea()) * StoragePerArea)
	town.Constructions = append(town.Constructions, c)

	buildingF := m.GetField(b.X, b.Y)
	town.AddConstructionTasks(c, buildingF, m)
}

func (town *Town) CreateIncrementalBuildingConstruction(b *building.Building, cost []artifacts.Artifacts, m navigation.IMap) {
	bt := b.Plan.BuildingType
	c := &building.Construction{X: b.X, Y: b.Y, Building: b, Cost: cost, T: bt, Storage: &artifacts.Resources{}}
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

func (town *Town) CreateLevelingTask(f *navigation.Field, taskType uint8, m navigation.IMap) {
	if town.Townhall.Household.NumTasks("terraform", economy.TerraformTaskTag(f)) == 0 {
		f.Construction = true
		town.Townhall.Household.AddTask(&economy.TerraformTask{
			F: f,
			M: m,
			T: taskType,
		})
	}
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

func (town *Town) CreateDemolishTask(b *building.Building) {
}

func (town *Town) GetHouseholds() []*Household {
	var households []*Household
	for _, f := range town.Farms {
		households = append(households, &f.Household)
	}
	for _, w := range town.Workshops {
		households = append(households, &w.Household)
	}
	for _, m := range town.Mines {
		households = append(households, &m.Household)
	}
	for _, f := range town.Factories {
		households = append(households, &f.Household)
	}
	for _, t := range town.Towers {
		households = append(households, &t.Household)
	}
	return households
}
