package social

import (
	"math"
	"math/rand"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/stats"
	"medvil/model/time"
	"medvil/model/vehicles"
)

const ReproductionRate = 1.0 / (24 * 30 * 12)
const TextileConsumptionRate = 1.0 / (24 * 30 * 12 * 5)
const StoragePerArea = 50
const ExtrasBudgetRatio = 0.25

var Log = artifacts.GetArtifact("log")
var Firewood = artifacts.GetArtifact("firewood")
var Tools = artifacts.GetArtifact("tools")
var Textile = artifacts.GetArtifact("textile")

const LogToFirewood = 5
const MinLog = 1

type Household struct {
	People          []*Person
	TargetNumPeople uint16
	Money           uint32
	Building        *building.Building
	Town            *Town
	Tasks           []economy.Task
	Vehicles        []*vehicles.Vehicle
	Resources       artifacts.Resources
	Heating         float64
}

func (h *Household) NextTask(m navigation.IMap, e economy.Equipment) economy.Task {
	return h.getNextTaskCombineExchange(m, e)
}

func (h *Household) getNextTaskCombineExchange(m navigation.IMap, e economy.Equipment) economy.Task {
	firstTask := FirstUnblockedTask(h, e)
	if firstTask != nil && IsExchangeBaseTask(firstTask) {
		vehicle := h.GetVehicle()
		et := GetExchangeTask(h, h.Town.Marketplace, m, vehicle)
		if et == nil && vehicle != nil {
			vehicle.SetInUse(false)
		}
		if et != nil {
			return et
		}
	}
	return GetNextTask(h, e)
}

func (h *Household) AddTask(t economy.Task) {
	h.Tasks = append(h.Tasks, t)
}

func (h *Household) AddPriorityTask(t economy.Task) {
	h.Tasks = append([]economy.Task{t}, h.Tasks...)
}

func (h *Household) GetTasks() []economy.Task {
	return h.Tasks
}

func (h *Household) SetTasks(tasks []economy.Task) {
	h.Tasks = tasks
}

func (h *Household) IncTargetNumPeople() {
	if h.TargetNumPeople < h.Building.Plan.Area()*2 {
		h.TargetNumPeople++
	}
}

func (h *Household) DecTargetNumPeople() {
	if h.TargetNumPeople > 0 {
		h.TargetNumPeople--
	}
}

func (h *Household) HasRoomForPeople() bool {
	return uint16(len(h.People)) < h.TargetNumPeople
}

func (h *Household) HasSurplusPeople() bool {
	return uint16(len(h.People)) > h.TargetNumPeople
}

func (h *Household) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	if &h.Town.Townhall.Household != h { // Not Townhall, needs better check
		if h.HasRoomForPeople() {
			h.Town.Townhall.Household.ReassignFirstPerson(h, m)
		}
		if len(h.People) >= 2 && rand.Float64() < ReproductionRate {
			if h.HasRoomForPeople() {
				h.People = append(h.People, h.NewPerson(m))
			} else if h.Town.Townhall.Household.HasRoomForPeople() {
				person := h.Town.Townhall.Household.NewPerson(m)
				h.Town.Townhall.Household.People = append(h.Town.Townhall.Household.People, person)
				person.Traveller.FX = h.Building.X
				person.Traveller.FY = h.Building.Y
				person.Task = &economy.GoHomeTask{F: m.GetField(h.Town.Townhall.Household.Building.X, h.Town.Townhall.Household.Building.Y), P: person}
			}
		}
		if h.HasSurplusPeople() && h.Town.Townhall.Household.HasRoomForPeople() {
			h.ReassignFirstPerson(&h.Town.Townhall.Household, m)
		}
	}
	numP := uint16(len(h.People))
	FindWaterTask(h, numP, m)
	mp := h.Town.Marketplace
	GetFoodTasks(h, numP, mp)
	numTools := h.Resources.Get(Tools) + h.PeopleWithTools()
	if numP > numTools && h.NumTasks("exchange", "tools_purchase") == 0 {
		needs := []artifacts.Artifacts{artifacts.Artifacts{A: Tools, Quantity: 1}}
		if h.Money >= mp.Price(needs) && mp.HasTraded(Tools) {
			h.AddTask(&economy.BuyTask{
				Exchange:       mp,
				HouseholdMoney: &h.Money,
				Goods:          needs,
				MaxPrice:       uint32(float64(h.Money) * ExtrasBudgetRatio),
				TaskTag:        "tools_purchase",
			})
		}
	}

	h.MaybeBuyExtras(Log, MinLog, "heating_fuel_shopping")
	h.MaybeBuyExtras(economy.Medicine, numP, "medicine_shopping")
	h.MaybeBuyExtras(economy.Beer, numP, "beer_shopping")
	h.MaybeBuyExtras(Textile, h.textileNeeded(), "textile_shopping")

	if Calendar.Hour == 0 {
		for i := 0; i < len(h.Tasks); i++ {
			if h.Tasks[i].IsPaused() {
				h.Tasks[i].Pause(false)
			}
		}
	}
	if h.Resources.Get(Firewood) < h.heatingFuelNeeded() && h.Resources.Remove(Log, 1) > 0 {
		h.Resources.Add(Firewood, LogToFirewood)
	}
	if h.Resources.Get(Textile) > 0 && rand.Float64() < TextileConsumptionRate*float64(numP) {
		h.Resources.Remove(Textile, 1)
	}
	if Calendar.Day == 1 && Calendar.Hour == 0 {
		if Calendar.Season() == time.Winter {
			wood := h.Resources.Remove(Firewood, h.heatingFuelNeededPerMonth())
			heating := float64(wood) / float64(h.heatingFuelNeededPerMonth())
			if h.HasEnoughTextile() {
				h.Heating = heating
			} else {
				h.Heating = math.Max(heating, 0.5)
			}
		} else {
			h.Heating = 1.0
		}
	}
}

func (h *Household) MaybeBuyExtras(a *artifacts.Artifact, threshold uint16, tag string) {
	mp := h.Town.Marketplace
	if h.Resources.Get(a) < threshold {
		if NumBatchesSimple(ProductTransportQuantity(a), ProductTransportQuantity(a)) > h.NumTasks("exchange", tag) {
			needs := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: ProductTransportQuantity(a)}}
			if h.Money >= mp.Price(needs) && mp.HasTraded(a) {
				h.AddTask(&economy.BuyTask{
					Exchange:       mp,
					HouseholdMoney: &h.Money,
					Goods:          needs,
					MaxPrice:       uint32(float64(h.Money) * ExtrasBudgetRatio),
					TaskTag:        tag,
				})
			}
		}
	}
}

func (h *Household) MaybeBuyBoat(Calendar *time.CalendarType, m navigation.IMap) {
	if h.numVehicles(vehicles.Boat) == 0 && h.Building.HasExtension(building.Deck) && h.NumTasks("factory_pickup", economy.BoatConstruction.Name) == 0 {
		factory := PickFactory(h.Town.Factories)
		if factory != nil && factory.Price(economy.BoatConstruction) < uint32(float64(h.Money)*ExtrasBudgetRatio) {
			ext, hx, hy := h.Building.GetExtensionWithCoords(building.Deck)
			fx, fy, fok := GetRandomBuildingXY(factory.Household.Building, m, navigation.Field.BuildingNonExtension)
			if ext != nil && ext.T == building.Deck && fok {
				order := factory.CreateOrder(economy.BoatConstruction, h)
				h.AddTask(&economy.FactoryPickupTask{
					PickupF:  m.GetField(fx, fy),
					DropoffF: m.GetField(hx, hy),
					Order:    order,
					TaskBase: economy.TaskBase{FieldCenter: true},
				})
			}
		}
	}
}

func (h *Household) MaybeBuyCart(Calendar *time.CalendarType, m navigation.IMap) {
	if h.numVehicles(vehicles.Cart) == 0 && h.NumTasks("factory_pickup", economy.CartConstruction.Name) == 0 {
		factory := PickFactory(h.Town.Factories)
		if factory != nil && factory.Price(economy.CartConstruction) < uint32(float64(h.Money)*ExtrasBudgetRatio) {
			hx, hy, _ := GetRandomBuildingXY(h.Building, m, navigation.Field.BuildingNonExtension)
			fx, fy, fok := GetRandomBuildingXY(factory.Household.Building, m, navigation.Field.BuildingNonExtension)
			if fok {
				order := factory.CreateOrder(economy.CartConstruction, h)
				h.AddTask(&economy.FactoryPickupTask{
					PickupF:  m.GetField(fx, fy),
					DropoffF: m.GetField(hx, hy),
					Order:    order,
					TaskBase: economy.TaskBase{FieldCenter: true},
				})
			}
		}
	}
}

func (h *Household) numVehicles(t *vehicles.VehicleType) int {
	var n = 0
	for _, v := range h.Vehicles {
		if v.T == t {
			n++
		}
	}
	return n
}

func (h *Household) textileNeeded() uint16 {
	return uint16(len(h.People)) + 1
}

func (h *Household) HasEnoughTextile() bool {
	return uint16(len(h.People)) <= h.Resources.Get(Textile)
}

func (h *Household) heatingFuelNeededPerMonth() uint16 {
	fuel := uint16(len(h.People) / 3)
	if fuel > 0 {
		return fuel
	}
	return 1
}

func (h *Household) heatingFuelNeeded() uint16 {
	return h.heatingFuelNeededPerMonth() * time.NumWinterMonths
}

func (h *Household) PeopleWithTools() uint16 {
	var n = uint16(0)
	for _, p := range h.People {
		if p.Equipment.Tool() {
			n++
		}
	}
	return n
}

func NotInputOrProduct(*artifacts.Artifact) bool {
	return false
}

func (h *Household) SellArtifacts(isInput func(*artifacts.Artifact) bool, isProduct func(*artifacts.Artifact) bool) {
	for a, q := range h.Resources.Artifacts {
		qToSell := h.ArtifactToSell(a, q, isInput(a), isProduct(a))
		if qToSell > 0 {
			tag := "sell_artifacts#" + a.Name
			if NumBatchesSimple(qToSell, ProductTransportQuantity(a)) > h.NumTasks("exchange", tag) {
				goods := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: ProductTransportQuantityWithLimit(a, qToSell)}}
				h.AddTask(&economy.SellTask{
					Exchange: h.Town.Marketplace,
					Goods:    goods,
					TaskTag:  tag,
				})
			}
		}
	}
}

func (h *Household) ArtifactToSell(a *artifacts.Artifact, q uint16, isInput bool, isProduct bool) uint16 {
	if isInput {
		return 0
	}
	p := uint16(len(h.People))
	if a.Name == "water" || a == Firewood {
		return 0
	}
	if a == Tools && !isProduct {
		return 0
	}
	var result uint16
	if economy.IsFoodOrDrink(a) {
		var threshold = economy.MaxFoodOrDrinkPerPerson
		if isProduct {
			threshold = economy.ProductMaxFoodOrDrinkPerPerson
		}
		if q > threshold*p {
			result = q - threshold*p
		} else {
			return 0
		}
	} else {
		result = q
	}
	if a == Log {
		if q >= MinLog {
			result = q - 1
		} else {
			return 0
		}
	}
	if a == Textile {
		textile := h.textileNeeded() + ProductTransportQuantity(Textile)
		if q > textile {
			result = q - textile
		} else {
			return 0
		}
	}
	if result >= ProductTransportQuantity(a) || h.Resources.Full() {
		return result
	}
	return 0
}

func (h *Household) HasFood() bool {
	return economy.HasFood(h.Resources)
}

func (h *Household) HasDrink() bool {
	return economy.HasDrink(h.Resources)
}

func (h *Household) HasMedicine() bool {
	return economy.HasMedicine(h.Resources)
}

func (h *Household) HasBeer() bool {
	return economy.HasBeer(h.Resources)
}

func (h *Household) NumTasks(name string, tag string) int {
	var i = 0
	for _, t := range h.Tasks {
		i += CountTags(t, name, tag)
	}
	for _, p := range h.People {
		if p.Task != nil {
			i += CountTags(p.Task, name, tag)
		}
	}
	return i
}

func (h *Household) NewPerson(m navigation.IMap) *Person {
	hx, hy, _ := GetRandomBuildingXY(h.Building, m, func(navigation.Field) bool { return true })
	return &Person{
		Food:      MaxPersonState,
		Water:     MaxPersonState,
		Happiness: MaxPersonState,
		Health:    MaxPersonState,
		Home:      h,
		Task:      nil,
		IsHome:    true,
		Traveller: &navigation.Traveller{
			FX: hx,
			FY: hy,
			FZ: 0,
			PX: 0,
			PY: 0,
			T:  navigation.TravellerTypePedestrian,
		},
		Equipment: &economy.NoEquipment{},
	}
}

func (h *Household) Filter(Calendar *time.CalendarType, m navigation.IMap) {
	var newPeople = make([]*Person, 0, len(h.People))
	for _, p := range h.People {
		if p.Health == 0 {
			m.GetField(p.Traveller.FX, p.Traveller.FY).UnregisterTraveller(p.Traveller)
			if p.Task != nil && !economy.IsPersonalTask(p.Task.Name()) {
				h.AddTask(p.Task)
			}
			p.releaseTask()
		} else {
			newPeople = append(newPeople, p)
		}
	}
	h.People = newPeople

	var newTasks = make([]economy.Task, 0, len(h.Tasks))
	for _, t := range h.Tasks {
		if !t.Expired(Calendar) {
			newTasks = append(newTasks, t)
		}
	}
	h.Tasks = newTasks
}

func (h *Household) AddVehicle(v *vehicles.Vehicle) {
	h.Vehicles = append(h.Vehicles, v)
}

func (h *Household) GetVehicle() *vehicles.Vehicle {
	for _, v := range h.Vehicles {
		if !v.InUse {
			v.SetInUse(true)
			return v
		}
	}
	return nil
}

func (h *Household) Stats() *stats.Stats {
	return &stats.Stats{
		Money:     h.Money,
		People:    uint32(len(h.People)),
		Buildings: 1,
		Artifacts: h.Resources.NumArtifacts(),
	}
}

func (srcH *Household) ReassignFirstPerson(dstH *Household, m navigation.IMap) {
	for pi, person := range srcH.People {
		if person.Task == nil {
			srcH.People = append(srcH.People[:pi], srcH.People[pi+1:]...)
			dstH.AssignPerson(person, m)
			person.Task = &economy.GoHomeTask{F: m.GetField(dstH.Building.X, dstH.Building.Y), P: person}
			break
		}
	}
}

func (h *Household) AssignPerson(person *Person, m navigation.IMap) {
	person.Home = h
	h.People = append(h.People, person)
}

func (h *Household) Field(m navigation.IMap) *navigation.Field {
	return m.GetField(h.Building.X, h.Building.Y)
}

func (h *Household) RandomField(m navigation.IMap, check func(navigation.Field) bool) *navigation.Field {
	x, y, ok := GetRandomBuildingXY(h.Building, m, check)
	if ok {
		return m.GetField(x, y)
	}
	return nil
}

func (h *Household) GetResources() *artifacts.Resources {
	return &h.Resources
}

func (h *Household) GetBuilding() *building.Building {
	return h.Building
}

func (h *Household) GetHeating() float64 {
	return h.Heating
}

func (h *Household) GetMoney() *uint32 {
	return &h.Money
}

func (h *Household) Destroy(m navigation.IMap) {
	dstH := &h.Town.Townhall.Household
	for _, person := range h.People {
		dstH.AssignPerson(person, m)
	}
	dstH.Money += h.Money
	m.GetField(h.Building.X, h.Building.Y).Terrain.Resources.AddResources(h.Resources)
}
