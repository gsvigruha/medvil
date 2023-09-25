package social

import (
	"math"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/stats"
	"medvil/model/time"
	"medvil/model/vehicles"
)

type Expedition struct {
	Money           uint32
	People          []*Person
	TargetNumPeople uint16
	Vehicle         *vehicles.Vehicle
	Resources       *artifacts.Resources
	StorageTarget   map[*artifacts.Artifact]int
	Tasks           []economy.Task
	Town            *Town
}

const MaxDistanceFromTown = 10

func (e *Expedition) DistanceToTown() float64 {
	return math.Abs(float64(e.Town.Townhall.Household.Building.X)-float64(e.Vehicle.Traveller.FX)) +
		math.Abs(float64(e.Town.Townhall.Household.Building.Y)-float64(e.Vehicle.Traveller.FY))
}

func (e *Expedition) CloseToTown(m navigation.IMap) bool {
	if e.Vehicle.T.Water && !m.Shore(e.Vehicle.Traveller.FX, e.Vehicle.Traveller.FY) {
		return false
	}
	return e.DistanceToTown() <= MaxDistanceFromTown
}

func (e *Expedition) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	for _, person := range e.People {
		person.ElapseTime(Calendar, m)
	}

	numP := uint16(len(e.People))
	FindWaterTask(e, numP, m)

	if e.CloseToTown(m) {
		if e.Town.Settings.UseSupplier {
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
				if e.NumTasks("transport", economy.TransportTaskTag(pickupD, a)) == 0 {
					targetQ := uint16(e.StorageTarget[a])
					if q < targetQ {
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

		for i := 0; i < len(e.Tasks); i++ {
			if e.Tasks[i].IsPaused() {
				e.Tasks[i].Pause(false)
			}
		}
	}
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

func (e *Expedition) NumTasks(name string, tag string) int {
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
