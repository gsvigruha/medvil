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
const MaxSubsidyRatio = 0.8

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

type TownSettings struct {
	RoadRepairs        bool
	WallRepairs        bool
	Trading            bool
	ArtifactCollection bool
	Coinage            bool
}

var DefaultTownSettings = TownSettings{
	RoadRepairs:        true,
	WallRepairs:        true,
	Trading:            true,
	ArtifactCollection: true,
	Coinage:            true,
}

var DefaultStorageTarget = map[string]int{
	"fruit":     20,
	"vegetable": 20,
	"bread":     20,
	"cube":      20,
	"brick":     20,
	"board":     20,
	"tile":      10,
	"thatch":    10,
	"log":       20,
	"textile":   20,
}

type Town struct {
	Country       *Country `json:"-"`
	Townhall      *Townhall
	Marketplace   *Marketplace
	Farms         []*Farm
	Workshops     []*Workshop
	Mines         []*Mine
	Factories     []*Factory
	Towers        []*Tower
	Walls         []*Wall
	Constructions []*building.Construction
	Transfers     *MoneyTransfers
	Roads         []*navigation.Field
	Settings      TownSettings
	Stats         *stats.Stats
	History       *stats.History
	Supplier      *Town
}

func (town *Town) Init() {
	defaultTransfers := TransferCategories{
		Rate:      30,
		Threshold: 200,
	}
	militaryTransfers := TransferCategories{
		Rate:      0,
		Threshold: 100,
	}
	town.Transfers = &MoneyTransfers{
		Farm:              defaultTransfers,
		Workshop:          defaultTransfers,
		Mine:              defaultTransfers,
		Factory:           defaultTransfers,
		Tower:             militaryTransfers,
		Trader:            defaultTransfers,
		MarketFundingRate: 70,
	}
	town.History = &stats.History{}
	town.ArchiveHistory()

	town.Townhall.StorageTarget = make(map[*artifacts.Artifact]int)
	for _, a := range artifacts.All {
		var amount int = 0
		if q, ok := town.Townhall.Household.Resources.Artifacts[a]; ok {
			amount = int(q)
		}
		town.Townhall.StorageTarget[a] = amount
	}
}

func (town *Town) ElapseTime(Calendar *time.CalendarType, m IMap) {
	s := town.Stats
	s.Reset()
	eoYear := (Calendar.Hour == 0 && Calendar.Day == 1 && Calendar.Month == 1)
	eoMonth := (Calendar.Hour == 0 && Calendar.Day == 1)
	if town.Marketplace != nil {
		town.Marketplace.ElapseTime(Calendar, m)
		s.Add(town.Marketplace.Stats())
		if eoMonth {
			town.Transfers.FundMarket(&town.Townhall.Household.Money, &town.Marketplace.Money)
		}
	}
	for l := range town.Townhall.Household.People {
		person := town.Townhall.Household.People[l]
		person.ElapseTime(Calendar, m)
	}
	town.Townhall.ElapseTime(Calendar, m)
	town.Townhall.Household.Filter(Calendar, m)
	town.Townhall.Filter(Calendar, m)
	s.Add(town.Townhall.Household.Stats())
	for _, trader := range town.Townhall.Traders {
		s.Add(trader.Stats())
	}
	for k := range town.Farms {
		farm := town.Farms[k]
		for l := range farm.Household.People {
			person := farm.Household.People[l]
			person.ElapseTime(Calendar, m)
		}
		farm.ElapseTime(Calendar, m)
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
		tower.Household.Filter(Calendar, m)
		s.Add(tower.Household.Stats())
	}
	if eoYear {
		CollectTax(town.Farms, town, town.Transfers.Farm)
		CollectTax(town.Mines, town, town.Transfers.Mine)
		CollectTax(town.Workshops, town, town.Transfers.Workshop)
		CollectTax(town.Towers, town, town.Transfers.Tower)
		CollectTax(town.Factories, town, town.Transfers.Factory)
		CollectTax(town.Townhall.Traders, town, town.Transfers.Trader)

		budget := uint32(float64(town.Townhall.Household.Money) * MaxSubsidyRatio)
		subsidyNeeded := (SumSubsidyNeeded(town.Farms, town.Transfers.Farm) +
			SumSubsidyNeeded(town.Mines, town.Transfers.Mine) +
			SumSubsidyNeeded(town.Workshops, town.Transfers.Workshop) +
			SumSubsidyNeeded(town.Towers, town.Transfers.Tower) +
			SumSubsidyNeeded(town.Factories, town.Transfers.Factory) +
			SumSubsidyNeeded(town.Townhall.Traders, town.Transfers.Trader))
		var ratio = 1.0
		if budget < subsidyNeeded {
			ratio = float64(budget) / float64(subsidyNeeded)
		}
		SendSubsidy(town.Farms, town, town.Transfers.Farm, ratio)
		SendSubsidy(town.Mines, town, town.Transfers.Mine, ratio)
		SendSubsidy(town.Workshops, town, town.Transfers.Workshop, ratio)
		SendSubsidy(town.Towers, town, town.Transfers.Tower, ratio)
		SendSubsidy(town.Factories, town, town.Transfers.Factory, ratio)
		SendSubsidy(town.Townhall.Traders, town, town.Transfers.Trader, ratio)
	}
	for k := range town.Walls {
		wall := town.Walls[k]
		wall.ElapseTime(Calendar, m)
	}
	var constructions []*building.Construction
	for k := range town.Constructions {
		construction := town.Constructions[k]
		if construction.IsComplete() {
			b := construction.Building
			field := m.GetField(construction.X, construction.Y)
			switch construction.T {
			case building.BuildingTypeMine:
				mine := &Mine{Household: &Household{Building: b, Town: town, Resources: &artifacts.Resources{}}}
				mine.Household.Resources.VolumeCapacity = b.Plan.Area() * StoragePerArea
				town.Mines = append(town.Mines, mine)
			case building.BuildingTypeWorkshop:
				w := &Workshop{Household: &Household{Building: b, Town: town, Resources: &artifacts.Resources{}}}
				w.Household.Resources.VolumeCapacity = b.Plan.Area() * StoragePerArea
				town.Workshops = append(town.Workshops, w)
			case building.BuildingTypeFarm:
				f := &Farm{Household: &Household{Building: b, Town: town, Resources: &artifacts.Resources{}}}
				f.Household.Resources.VolumeCapacity = b.Plan.Area() * StoragePerArea
				town.Farms = append(town.Farms, f)
			case building.BuildingTypeFactory:
				f := &Factory{Household: &Household{Building: b, Town: town, Resources: &artifacts.Resources{}}}
				f.Household.Resources.VolumeCapacity = b.Plan.Area() * StoragePerArea
				town.Factories = append(town.Factories, f)
			case building.BuildingTypeTower:
				t := &Tower{Household: &Household{Building: b, Town: town, Resources: &artifacts.Resources{}}}
				t.Household.Resources.VolumeCapacity = b.Plan.Area() * StoragePerArea
				town.Towers = append(town.Towers, t)
			case building.BuildingTypeWall, building.BuildingTypeGate:
				w := &Wall{Building: b, Town: town, F: field}
				town.Walls = append(town.Walls, w)
			case building.BuildingTypeTownhall:
				town.Country.CreateNewTown(b, town)
			case building.BuildingTypeMarket:
				town.Marketplace = &Marketplace{Town: town, Building: b}
				town.Marketplace.Storage.VolumeCapacity = b.Plan.Area() * StoragePerArea
				town.Marketplace.Init()
			case building.BuildingTypeRoad:
				construction.Road.Construction = false
				navigation.SetRoadConnections(m, field)
				town.Roads = append(town.Roads, field)
			case building.BuildingTypeCanal:
				field.Construction = false
				field.Terrain.T = terrain.Canal
			case building.BuildingTypeStatue:
				construction.Statue.Construction = false
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
	if Calendar.Day == 30 && Calendar.Hour == 0 && town.Townhall.Household.Resources.Remove(Paper, 1) > 0 {
		if town.Settings.RoadRepairs {
			for _, road := range town.Roads {
				if road.Road.Broken && town.Townhall.Household.NumTasks("repair", economy.BuildingTaskTag(road)) == 0 {
					f := m.GetField(road.X, road.Y)
					res := &artifacts.Resources{}
					res.Init(25)
					town.Townhall.Household.AddTask(&economy.RepairTask{
						Repairable: road.Road,
						Field:      f,
						Resources:  res,
					})
					town.AddTransportTasks(road.Road.RepairCost(), f, res, m)
				}
			}
		}
		if town.Settings.WallRepairs {
			for _, wall := range town.Walls {
				wf := m.GetField(wall.F.X, wall.F.Y)
				res := &artifacts.Resources{}
				res.Init(25)
				if wall.Building.Broken && town.Townhall.Household.NumTasks("repair", economy.BuildingTaskTag(wf)) == 0 {
					town.Townhall.Household.AddTask(&economy.RepairTask{
						Repairable: wall.Building,
						Field:      wall.F,
						Resources:  res,
					})
					town.AddTransportTasks(wall.Building.RepairCost(), wall.F, res, m)
				}
			}
		}
		if town.Settings.ArtifactCollection {
			for i := -TownhallMaxDistance; i <= TownhallMaxDistance; i++ {
				for j := -TownhallMaxDistance; j <= TownhallMaxDistance; j++ {
					f := m.GetField(uint16(int(town.Townhall.Household.Building.X)+i), uint16(int(town.Townhall.Household.Building.Y)+j))
					if f != nil && !f.Allocated && town.Townhall.FieldWithinDistance(f) {
						town.AddTransportTask(f)
					}
				}
			}
		}
	}
}

func (town *Town) AddTransportTask(f *navigation.Field) {
	for a, q := range f.Terrain.Resources.Artifacts {
		if f.Terrain.Resources.IsRealArtifact(a) && q > 0 {
			tag := economy.TransportTaskTag(f, a)
			if town.Townhall.Household.NumTasks("transport", tag) == 0 {
				town.Townhall.Household.AddTask(&economy.TransportTask{
					PickupD:        f,
					DropoffD:       town.Townhall.Household.Destination(building.NonExtension),
					PickupR:        f.Terrain.Resources,
					DropoffR:       town.Townhall.Household.Resources,
					A:              a,
					TargetQuantity: q,
				})
			}
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

func (town *Town) CreateStatueConstruction(x, y uint16, s *building.Statue, m navigation.IMap) {
	c := &building.Construction{X: x, Y: y, Statue: s, Cost: s.T.Cost, T: building.BuildingTypeStatue, Storage: &artifacts.Resources{}}
	c.Storage.Init(StoragePerArea)
	town.Constructions = append(town.Constructions, c)

	f := m.GetField(x, y)
	town.AddConstructionTasks(c, f, m)
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
	var dest navigation.Destination
	if c.Building != nil && c.Building.Broken {
		dest = buildingF
	} else {
		dest = buildingF.TopLocation()
	}
	totalTasks := town.AddTransportTasks(c.Cost, dest, c.Storage, m)
	if totalTasks == 0 {
		totalTasks = 1
	}
	c.MaxProgress = totalTasks
	for i := uint16(0); i < totalTasks; i++ {
		town.Townhall.Household.AddTask(&economy.BuildingTask{
			D: dest,
			C: c,
		})
	}
}

func (town *Town) AddTransportTasks(cost []artifacts.Artifacts, dest navigation.Destination, storage *artifacts.Resources, m navigation.IMap) uint16 {
	var totalTasks uint16 = 0
	for _, a := range cost {
		var totalQ = a.Quantity
		totalTasks += totalQ
		for totalQ > 0 {
			var q uint16 = ConstructionTransportQuantity
			if totalQ < ConstructionTransportQuantity {
				q = totalQ
			}
			totalQ -= q
			town.Townhall.Household.AddTask(&economy.TransportTask{
				PickupD:          m.GetField(town.Townhall.Household.Building.X, town.Townhall.Household.Building.Y),
				DropoffD:         dest,
				PickupR:          town.Townhall.Household.Resources,
				DropoffR:         storage,
				A:                a.A,
				TargetQuantity:   q,
				CompleteQuantity: true,
			})
		}
	}
	return totalTasks
}

func (town *Town) CreateDemolishTask(b *building.Building, r *building.Road, f *navigation.Field, m navigation.IMap) {
	if town.Townhall.Household.NumTasks("demolish", economy.DemolishTaskTag(f)) == 0 {
		town.Townhall.Household.AddTask(&economy.DemolishTask{
			Building: b,
			Road:     r,
			F:        f,
			Town:     town,
			M:        m,
		})
	}
}

func (town *Town) GetHouseholds() []*Household {
	var households []*Household
	for _, f := range town.Farms {
		households = append(households, f.Household)
	}
	for _, w := range town.Workshops {
		households = append(households, w.Household)
	}
	for _, m := range town.Mines {
		households = append(households, m.Household)
	}
	for _, f := range town.Factories {
		households = append(households, f.Household)
	}
	for _, t := range town.Towers {
		households = append(households, t.Household)
	}
	return households
}

func AddTransportTasksForField(field *navigation.Field, th *Townhall, m navigation.IMap) {
	for a, q := range field.Terrain.Resources.Artifacts {
		if q > 0 {
			th.Household.AddTask(&economy.TransportTask{
				PickupD:        field,
				DropoffD:       m.GetField(th.Household.Building.X, th.Household.Building.Y),
				PickupR:        field.Terrain.Resources,
				DropoffR:       th.Household.Resources,
				A:              a,
				TargetQuantity: q,
			})
		}
	}
}

func DestroyBuilding[H House](houses []H, b *building.Building, m navigation.IMap) []H {
	var newHouses []H
	for _, house := range houses {
		household := house.GetHome().(*Household)
		if household.Building == b {
			// Remove the building elements from the field
			for _, coords := range b.GetBuildingXYs(true) {
				m.GetField(coords[0], coords[1]).Building = navigation.FieldBuildingObjects{}
			}
			// Land used by the house to be destroyed should be unallocated
			for _, field := range house.GetFields() {
				field.Field().Allocated = false
			}
			household.Destroy(m)
			AddTransportTasksForField(m.GetField(b.X, b.Y), household.Town.Townhall, m)
		} else {
			newHouses = append(newHouses, house)
		}
	}
	return newHouses
}

func (town *Town) DestroyRoad(r *building.Road, m navigation.IMap) {
	var newRoads []*navigation.Field
	for _, road := range town.Roads {
		if road.Road == r {
			f := m.GetField(road.X, road.Y)
			f.Road = nil
			f.Statue = nil
		} else {
			newRoads = append(newRoads, road)
		}
	}
	town.Roads = newRoads
}

func (town *Town) DestroyBuilding(b *building.Building, m navigation.IMap) {
	switch b.Plan.BuildingType {
	case building.BuildingTypeFarm:
		town.Farms = DestroyBuilding(town.Farms, b, m)
	case building.BuildingTypeMine:
		town.Mines = DestroyBuilding(town.Mines, b, m)
	case building.BuildingTypeWorkshop:
		town.Workshops = DestroyBuilding(town.Workshops, b, m)
	case building.BuildingTypeFactory:
		town.Factories = DestroyBuilding(town.Factories, b, m)
	case building.BuildingTypeTower:
		town.Towers = DestroyBuilding(town.Towers, b, m)
	case building.BuildingTypeMarket:
		for _, coords := range town.Marketplace.Building.GetBuildingXYs(true) {
			m.GetField(coords[0], coords[1]).Building = navigation.FieldBuildingObjects{}
		}
		m.GetField(town.Marketplace.Building.X, town.Marketplace.Building.Y).Terrain.Resources.AddResources(town.Marketplace.Storage)
		town.Marketplace = nil
	case building.BuildingTypeWall, building.BuildingTypeGate:
		var newWalls []*Wall
		for _, wall := range town.Walls {
			if wall.Building == b {
				f := m.GetField(wall.Building.X, wall.Building.Y)
				f.Building = navigation.FieldBuildingObjects{}
				f.Terrain.Resources.AddAll(b.Plan.RepairCost())
				AddTransportTasksForField(f, town.Townhall, m)
			} else {
				newWalls = append(newWalls, wall)
			}
		}
		town.Walls = newWalls
	}
}

func (town *Town) ArchiveHistory() {
	var pt = make(map[economy.Task]uint32)
	if town.Stats != nil {
		town.History.Archive(town.Stats)
		pt = town.Stats.PendingTasks
	}
	town.Stats = &stats.Stats{}
	town.Stats.Init(pt)
}
