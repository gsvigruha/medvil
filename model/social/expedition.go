package social

import (
	"math"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/stats"
	"medvil/model/terrain"
	"medvil/model/time"
	"medvil/model/vehicles"
)

type Expedition struct {
	Name             string
	Money            uint32
	People           []*Person
	TargetNumPeople  uint16
	Vehicle          *vehicles.Vehicle
	Resources        *artifacts.Resources
	StorageTarget    map[*artifacts.Artifact]int
	Tasks            []economy.Task
	Town             *Town
	DestinationField *navigation.Field
	Constructions    []*building.Construction
	Autopilot        bool
}

const MaxDistanceFromTown = 15
const StorageFullRatio = 0.8

func (e *Expedition) DistanceToTown(town *Town) float64 {
	return math.Abs(float64(town.Townhall.Household.Building.X)-float64(e.Vehicle.Traveller.FX)) +
		math.Abs(float64(town.Townhall.Household.Building.Y)-float64(e.Vehicle.Traveller.FY))
}

func (e *Expedition) CloseToTown(town *Town, m navigation.IMap) bool {
	if e.Vehicle.T.Water && !m.Shore(e.Vehicle.Traveller.FX, e.Vehicle.Traveller.FY) {
		return false
	}
	return e.DistanceToTown(town) <= MaxDistanceFromTown
}

func (e *Expedition) ValidCampSpot(f *navigation.Field, m navigation.IMap) bool {
	return ((e.Vehicle.T.Water && f.Sailable() && m.Shore(f.X, f.Y)) || (!e.Vehicle.T.Water && f.Drivable()))
}

func (e *Expedition) PickRandomSpotNearTownhall(town *Town, m navigation.IMap) *navigation.Field {
	for d := 5; d <= MaxDistanceFromTown; d++ {
		for i := 0; i < 5; i++ {
			f := m.RandomSpot(town.Townhall.Household.Building.X, town.Townhall.Household.Building.Y, d)
			if f != nil && e.ValidCampSpot(f, m) {
				return f
			}
		}
	}
	return nil
}

func (e *Expedition) PickRandomNearBySpot(m navigation.IMap) *navigation.Field {
	f := e.PickRandomSpotNearTownhall(e.Town, m)
	if f != nil {
		return f
	}
	for i := 0; i < 20; i++ {
		f := m.RandomSpot(e.Vehicle.Traveller.FX, e.Vehicle.Traveller.FY, 5)
		if f != nil && e.ValidCampSpot(f, m) {
			return f
		}
	}
	return nil
}

func (e *Expedition) NeedsFood() bool {
	for _, a := range economy.Foods {
		if NumFoodBatchesNeeded(e, uint16(len(e.People)), a) == 0 {
			return false
		}
	}
	return true
}

func (e *Expedition) StorageFull() bool {
	for a, targetQ := range e.StorageTarget {
		if targetQ > 0 {
			q := e.Resources.Get(a)
			if float64(q)/float64(targetQ) < StorageFullRatio {
				return false
			}
		}
	}
	return true
}

func (e *Expedition) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	currentF := m.GetField(e.Vehicle.Traveller.FX, e.Vehicle.Traveller.FY)
	if e.DestinationField == nil && currentF != nil && !e.Vehicle.T.Water && !currentF.Drivable() {
		e.DestinationField = e.PickRandomNearBySpot(m)
	}

	if e.DestinationField != nil {
		for i := 0; i < len(e.Tasks); i++ {
			e.Tasks[i].Pause(true)
		}
		if e.IsEveryoneBoarded() {
			e.Vehicle.Traveller.UseVehicle(e.Vehicle)

			if e.Vehicle.Traveller.IsAtDestination(e.DestinationField) {
				if !e.Vehicle.Traveller.IsOnFieldCenter() {
					e.Vehicle.Traveller.MoveToCenter(m)
				} else {
					e.DestinationField = nil
				}
			} else {
				hasPath, computing := e.Vehicle.Traveller.EnsurePath(e.DestinationField, m)
				if hasPath {
					e.Vehicle.Traveller.Move(m)
				} else if !computing {
					e.DestinationField = nil // no path, cancel destination
				}
			}

			for _, p := range e.People {
				p.Traveller.MoveWith(m, e.Vehicle.Traveller)
			}
		}
	}

	for _, person := range e.People {
		person.ElapseTime(Calendar, m)
	}

	numP := uint16(len(e.People))
	FindWaterTask(e, numP, true, m)

	if e.CloseToTown(e.Town, m) && e.DestinationField == nil {
		srcH := e.Town.Townhall.Household
		if e.HasRoomForPeople() {
			srcH.ReassignFirstPerson(e, len(e.Tasks) == 0, m)
		}
		for _, a := range artifacts.All {
			var q uint16 = 0
			if storageQ, ok := e.Resources.Artifacts[a]; ok {
				q = storageQ
			}
			pickupD := m.GetField(srcH.Building.X, srcH.Building.Y)
			targetQ := uint16(e.StorageTarget[a])
			if q < targetQ {
				tasksNeeded := (targetQ - q) / ProductTransportQuantity(a)
				if e.NumTasks("transport", economy.TransportTaskTag(pickupD, a)) < int(tasksNeeded) {
					e.AddTask(&economy.TransportTask{
						PickupD:        pickupD,
						DropoffD:       &navigation.TravellerDestination{T: e.Vehicle.Traveller},
						PickupR:        srcH.Resources,
						DropoffR:       e.Resources,
						A:              a,
						TargetQuantity: ProductTransportQuantity(a),
					})
				}
			}
		}
	}

	if e.DestinationField == nil && Calendar.Hour == 0 && Calendar.Day == 1 {
		for i := 0; i < len(e.Tasks); i++ {
			if e.Tasks[i].IsPaused() {
				e.Tasks[i].Pause(false)
			}
		}
	}

	if Calendar.Hour == 0 && Calendar.Day == 1 && e.Autopilot {
		if e.DestinationField == nil && e.NeedsFood() && !e.CloseToTown(e.Town, m) && e.Resources.Remove(Paper, 1) > 0 {
			e.DestinationField = e.PickRandomSpotNearTownhall(e.Town, m)
		}
		if e.DestinationField == nil && e.StorageFull() && e.CloseToTown(e.Town, m) && e.Resources.Remove(Paper, 1) > 0 {
			for _, town := range e.Town.Country.Towns {
				if town.Supplier == e {
					e.DestinationField = e.PickRandomSpotNearTownhall(town, m)
					break
				}
			}
		}
	}

	var constructions []*building.Construction
	for k := range e.Constructions {
		construction := e.Constructions[k]
		if construction.IsComplete() {
			b := construction.Building
			switch construction.T {
			case building.BuildingTypeTownhall:
				e.Town.Country.CreateNewTown(b, e)
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
	e.Constructions = constructions
}

func (e *Expedition) IsEveryoneBoarded() bool {
	for _, p := range e.People {
		if !p.IsHome {
			return false
		}
	}
	return true
}

func (e *Expedition) HasRoomForPeople() bool {
	return uint16(len(e.People)) < e.TargetNumPeople
}

func (e *Expedition) AddTask(task economy.Task) {
	e.Tasks = append(e.Tasks, task)
}

func (e *Expedition) AddPriorityTask(task economy.Task) {
	e.Tasks = append([]economy.Task{task}, e.Tasks...)
}

func (e *Expedition) GetTasks() []economy.Task {
	return e.Tasks
}

func (e *Expedition) SetTasks(tasks []economy.Task) {
	e.Tasks = tasks
}

func (e *Expedition) HasFood() bool {
	return economy.HasFood(*e.Resources)
}

func (e *Expedition) HasDrink() bool {
	return economy.HasDrink(*e.Resources)
}

func (e *Expedition) HasMedicine() bool {
	return economy.HasMedicine(*e.Resources)
}

func (e *Expedition) HasBeer() bool {
	return economy.HasBeer(*e.Resources)
}

func (e *Expedition) Field(m navigation.IMap) *navigation.Field {
	return m.GetField(e.Vehicle.Traveller.FX, e.Vehicle.Traveller.FY)
}

func (e *Expedition) RandomField(m navigation.IMap, check func(navigation.Field) bool) *navigation.Field {
	return e.Field(m)
}

func (e *Expedition) NextTask(m navigation.IMap, et *economy.EquipmentType) economy.Task {
	return GetNextTask(e, et)
}

func (e *Expedition) GetResources() *artifacts.Resources {
	return e.Resources
}

func (e *Expedition) GetBuilding() *building.Building {
	return nil
}

func (e *Expedition) GetHeating() uint8 {
	return 100
}

func (e *Expedition) HasEnoughClothes() bool {
	return true
}

func (e *Expedition) AddVehicle(v *vehicles.Vehicle) {
}

func (e *Expedition) AllocateVehicle(waterOk bool) *vehicles.Vehicle {
	return e.Vehicle
}

func (e *Expedition) NumTasks(name string, tag economy.Tag) int {
	var i = 0
	for _, t := range e.Tasks {
		i += CountTags(t, name, tag)
	}
	for _, p := range e.People {
		if p.Task != nil {
			i += CountTags(p.Task, name, tag)
		}
	}
	return i
}

func (e *Expedition) Spend(amount uint32) {
	e.Money -= amount
}

func (e *Expedition) Earn(amount uint32) {
	e.Money += amount
}

func (e *Expedition) GetMoney() uint32 {
	return e.Money
}

func (e *Expedition) Destination(extensionType *building.BuildingExtensionType) navigation.Destination {
	return &navigation.TravellerDestination{T: e.Vehicle.Traveller}
}

func (e *Expedition) Stats() *stats.HouseholdStats {
	return &stats.HouseholdStats{
		Money:     e.Money,
		People:    uint32(len(e.People)),
		Buildings: 0,
		Artifacts: e.Resources.NumArtifacts(),
	}
}

func (e *Expedition) PendingCosts() uint32 {
	return PendingCosts(e.Tasks)
}

func (e *Expedition) Broken() bool {
	return false
}

func (e *Expedition) GetTown() *Town {
	return e.Town
}

func (e *Expedition) GetPeople() []*Person {
	return e.People
}

func (e *Expedition) GetHome() Home {
	return e
}

func (e *Expedition) GetExchange() economy.Exchange {
	return nil
}

func (e *Expedition) IsHomeVehicle() bool {
	return true
}

func (e *Expedition) IsBoatEnabled() bool {
	return true
}

func (e *Expedition) AssignPerson(person *Person, m navigation.IMap) {
	person.Home = e
	e.People = append(e.People, person)
}

func (e *Expedition) IncTargetNumPeople() {
	if e.TargetNumPeople < e.Vehicle.T.MaxPeople {
		e.TargetNumPeople++
	}
}

func (e *Expedition) DecTargetNumPeople() {
	if e.TargetNumPeople > 0 {
		e.TargetNumPeople--
	}
}

func (e *Expedition) Filter(Calendar *time.CalendarType, m navigation.IMap) {
	var newPeople = make([]*Person, 0, len(e.People))
	for _, p := range e.People {
		f := m.GetField(p.Traveller.FX, p.Traveller.FY)
		if p.Health > 0 && p.Home == e {
			newPeople = append(newPeople, p)
		} else {
			if p.Task != nil && !economy.IsPersonalTask(p.Task.Name()) {
				e.AddTask(p.Task)
			}
			f.UnregisterTraveller(p.Traveller)
			e.Town.Stats.RegisterDeath()
		}
	}
	e.People = newPeople

	var newTasks = make([]economy.Task, 0, len(e.Tasks))
	for _, t := range e.Tasks {
		if !t.Expired(Calendar) {
			newTasks = append(newTasks, t)
		}
	}
	e.Tasks = newTasks
}

func (e *Expedition) IsPersonVisible() bool {
	return false
}

func (e *Expedition) ReassignFirstPerson(dstH Home, assingTask bool, m navigation.IMap) {
	if len(e.People) > 1 && e.CloseToTown(dstH.GetTown(), m) {
		for pi, person := range e.People {
			if person.Task == nil {
				e.People = append(e.People[:pi], e.People[pi+1:]...)
				dstH.AssignPerson(person, m)
				if assingTask {
					person.Task = &economy.GoHomeTask{D: dstH.Destination(building.NonExtension), P: person}
				}
				break
			}
		}
	}
}

func (e *Expedition) FieldWithinDistance(field *navigation.Field) bool {
	return WithinDistanceCoords(e.Vehicle.Traveller.FX, e.Vehicle.Traveller.FY, field.X, field.Y, TownhallMaxDistance)
}

func (e *Expedition) CreateBuildingConstruction(b *building.Building, m navigation.IMap) {
	CreateBuildingConstruction(e, b, m)
}

func (e *Expedition) AddConstruction(c *building.Construction) {
	e.Constructions = append(e.Constructions, c)
}

func (e *Expedition) BuildMarketplaceEnabled() bool {
	return false
}

func (e *Expedition) BuildHousesEnabled() bool {
	return false
}

func (e *Expedition) CreateLevelingTask(f *navigation.Field, taskType uint8, m navigation.IMap) {
	if e.NumTasks("terraform", economy.TerraformTaskTag(f)) == 0 {
		f.Construction = true
		e.AddTask(&economy.TerraformTask{
			F: f,
			M: m,
			T: taskType,
		})
	}
}

func (e *Expedition) CheckDestinationField(f *navigation.Field) bool {
	if e.Vehicle.T.Water && f.Terrain.T == terrain.Water && f.Sailable() {
		return true
	} else if !e.Vehicle.T.Water && f.Terrain.T == terrain.Grass && f.Drivable() {
		return true
	}
	return false
}

func (e *Expedition) GetName() string {
	return e.Name
}
