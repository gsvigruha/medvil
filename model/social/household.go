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
const ClothesConsumptionRate = 1.0 / (24 * 30 * 12 * 3)
const StoragePerArea = 100
const ExtrasBudgetRatio = 0.25
const BuildingBrokenRate = 1.0 / (24 * 30 * 12 * 15)
const FleeingRate = 1.0 / (24 * 30 * 12 * 3)

var Log = artifacts.GetArtifact("log")
var Firewood = artifacts.GetArtifact("firewood")
var Tools = artifacts.GetArtifact("tools")
var Textile = artifacts.GetArtifact("textile")
var Leather = artifacts.GetArtifact("leather")
var Clothes = artifacts.GetArtifact("clothes")
var IronBar = artifacts.GetArtifact("iron_bar")

const LogToFirewood = 5
const MinLog = 1

type Household struct {
	People          []*Person
	TargetNumPeople uint16
	Money           uint32
	Building        *building.Building
	Town            *Town `json:"-"`
	Tasks           []economy.Task
	Vehicles        []*vehicles.Vehicle
	Resources       *artifacts.Resources
	Heating         uint8
}

func (h *Household) NextTask(m navigation.IMap, e *economy.EquipmentType) economy.Task {
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
	if h.Town.Townhall.Household != h { // Not Townhall, needs better check
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
			h.ReassignFirstPerson(h.Town.Townhall.Household, m)
		}
	}
	if h.NumTasks("exchange", "market") <= len(h.People)/3 {
		CombineExchangeTasks(h, h.Town.Marketplace, m)
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
				Exchange:        mp,
				HouseholdWallet: h,
				Goods:           needs,
				MaxPrice:        uint32(float64(h.Money) * ExtrasBudgetRatio),
				TaskTag:         "tools_purchase",
			})
		}
	}

	h.MaybeBuyExtras(Log, MinLog, "heating_fuel_shopping")
	h.MaybeBuyExtras(economy.Medicine, numP, "medicine_shopping")
	h.MaybeBuyExtras(economy.Beer, numP, "beer_shopping")
	if mp.Prices[Textile] < mp.Prices[Leather] {
		h.MaybeBuyExtras(Textile, h.clothesNeeded(), "textile_shopping")
	} else {
		h.MaybeBuyExtras(Leather, h.clothesNeeded(), "leather_shopping")
	}

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
	if !h.HasEnoughClothes() {
		if h.Resources.Remove(Textile, 1) > 0 {
			h.Resources.Add(Clothes, 1)
		} else if h.Resources.Remove(Leather, 1) > 0 {
			h.Resources.Add(Clothes, 1)
		}
	}
	if h.Resources.Get(Clothes) > 0 && rand.Float64() < ClothesConsumptionRate*float64(numP) {
		h.Resources.Remove(Clothes, 1)
	}
	if Calendar.Day == 1 && Calendar.Hour == 0 {
		if Calendar.Season() == time.Winter {
			wood := h.Resources.Remove(Firewood, h.heatingFuelNeededPerMonth())
			if wood > 0 {
				h.Heating = 100
			} else if h.HasEnoughClothes() {
				h.Heating = 50
			} else {
				h.Heating = 0
			}
		} else {
			h.Heating = 100
		}
	}
	if rand.Float64() < BuildingBrokenRate {
		h.Building.Broken = true
	}
	if h.Building.Broken {
		if h.NumTasks("repair", "") == 0 {
			h.AddTask(&economy.RepairTask{
				B: h.Building,
				F: m.GetField(h.Building.X, h.Building.Y),
				R: h.GetResources(),
			})
		}

		needs := h.Resources.Needs(h.Building.Plan.RepairCost())
		if len(needs) > 0 && h.NumTasks("exchange", "repair_shopping") == 0 {
			if h.Money >= mp.Price(needs) {
				h.AddTask(&economy.BuyTask{
					Exchange:        mp,
					HouseholdWallet: h,
					Goods:           needs,
					MaxPrice:        uint32(float64(h.Money) * ExtrasBudgetRatio),
					TaskTag:         "repair_shopping",
				})
			}
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
					Exchange:        mp,
					HouseholdWallet: h,
					Goods:           needs,
					MaxPrice:        uint32(float64(h.Money) * ExtrasBudgetRatio),
					TaskTag:         tag,
				})
			}
		}
	}
}

func (h *Household) MaybeBuyBoat(Calendar *time.CalendarType, m navigation.IMap) {
	if h.Building.HasExtension(building.Deck) && h.numVehicles(vehicles.Boat) == 0 && h.NumTasks("factory_pickup", economy.BoatConstruction.Name) == 0 {
		factory := PickFactory(h.Town.Factories, building.Deck)
		if factory != nil && factory.Price(economy.BoatConstruction) < uint32(float64(h.Money)*ExtrasBudgetRatio) {
			order := factory.CreateOrder(economy.BoatConstruction, h)
			if order != nil {
				h.AddTask(&economy.FactoryPickupTask{
					// Factories need to be accessed via land even for boat pickups first
					PickupD:  factory.Household.Destination(building.NonExtension),
					DropoffD: h.Destination(building.Deck),
					Order:    order,
					TaskBase: economy.TaskBase{FieldCenter: true},
				})
			}
		}
	}
}

func (h *Household) MaybeBuyCart(Calendar *time.CalendarType, m navigation.IMap) {
	if h.numVehicles(vehicles.Cart) < len(h.People)/2 && h.NumTasks("factory_pickup", economy.CartConstruction.Name) == 0 {
		factory := PickFactory(h.Town.Factories, building.NonExtension)
		if factory != nil && factory.Price(economy.CartConstruction) < uint32(float64(h.Money)*ExtrasBudgetRatio) {
			order := factory.CreateOrder(economy.CartConstruction, h)
			if order != nil {
				h.AddTask(&economy.FactoryPickupTask{
					PickupD:  factory.Household.Destination(building.NonExtension),
					DropoffD: h.Destination(building.NonExtension),
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

func (h *Household) clothesNeeded() uint16 {
	return uint16(len(h.People)) + 1
}

func (h *Household) HasEnoughClothes() bool {
	return uint16(len(h.People)) <= h.Resources.Get(Clothes)
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
		if p.Equipment.Tool {
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
	if a.Name == "water" || a == Firewood || a == Clothes {
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
			result = q - MinLog
		} else {
			return 0
		}
	}
	if a == Textile || a == Leather {
		needed := h.clothesNeeded()
		if q > needed {
			result = q - needed
		} else {
			return 0
		}
	}
	if a == economy.Beer || a == economy.Medicine {
		if q >= p {
			result = q - p
		} else {
			return 0
		}
	}
	if a == Paper {
		if isProduct {
			result = q
		} else {
			return 0
		}
	}
	if h.Building.Broken {
		needed := artifacts.GetQuantity(h.Building.Plan.RepairCost(), a)
		if q > needed {
			result = q - needed
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
	return economy.HasFood(*h.Resources)
}

func (h *Household) HasDrink() bool {
	return economy.HasDrink(*h.Resources)
}

func (h *Household) HasMedicine() bool {
	return economy.HasMedicine(*h.Resources)
}

func (h *Household) HasBeer() bool {
	return economy.HasBeer(*h.Resources)
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
	f := h.RandomField(m, func(navigation.Field) bool { return true })
	traveller := &navigation.Traveller{
		FX: f.X,
		FY: f.Y,
		FZ: 0,
		PX: 0,
		PY: 0,
		T:  navigation.TravellerTypePedestrian,
	}
	traveller.InitPathElement(f)
	person := &Person{
		Food:      MaxPersonState,
		Water:     MaxPersonState,
		Happiness: MaxPersonState,
		Health:    MaxPersonState,
		Home:      h,
		Task:      nil,
		IsHome:    true,
		Traveller: traveller,
		Equipment: economy.NoEquipment,
	}
	traveller.Person = person
	return person
}

func (h *Household) Filter(Calendar *time.CalendarType, m IMap) {
	var newPeople = make([]*Person, 0, len(h.People))
	for _, p := range h.People {
		if p.Health == 0 {
			f := m.GetField(p.Traveller.FX, p.Traveller.FY)
			f.UnregisterTraveller(p.Traveller)
			if p.Task != nil && !economy.IsPersonalTask(p.Task.Name()) {
				h.AddTask(p.Task)
			}
			p.releaseTask()
			if p.Equipment.Tool || p.Equipment.Weapon {
				f.Terrain.Resources.Add(IronBar, 1)
			}
			h.Town.Country.SocietyStats.RegisterDeath()
		} else if p.Happiness == 0 && rand.Float64() < FleeingRate && h.Town.Country.T != CountryTypeOutlaw {
			if p.Task != nil && !economy.IsPersonalTask(p.Task.Name()) {
				h.AddTask(p.Task)
			}
			p.releaseTask()

			var town *Town
			var dist = 0.0
			for _, countryI := range m.GetCountries(CountryTypeOutlaw) {
				for _, townI := range countryI.Towns {
					distI := math.Abs(float64(h.Town.Townhall.Household.Building.X)-float64(townI.Townhall.Household.Building.X)) +
						math.Abs(float64(h.Town.Townhall.Household.Building.Y)-float64(townI.Townhall.Household.Building.Y))
					if town == nil || dist > distI {
						town = townI
						dist = distI
					}
				}
			}
			h.Town.Country.SocietyStats.RegisterDeparture()
			if town != nil {
				town.Townhall.Household.AssignPerson(p, m)
				p.Task = &economy.GoHomeTask{F: m.GetField(town.Townhall.Household.Building.X, town.Townhall.Household.Building.Y), P: p}
			} else {
				m.GetField(p.Traveller.FX, p.Traveller.FY).UnregisterTraveller(p.Traveller)
			}
		} else if guard := m.GetNearbyGuard(p.Traveller); guard != nil && h.Town.Country.T == CountryTypeOutlaw {
			th := guard.Home.GetTown().Townhall
			th.Household.AssignPerson(p, m)
			p.Task = &economy.GoHomeTask{F: m.GetField(th.Household.Building.X, th.Household.Building.Y), P: p}
		} else if p.Home == h {
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

	var newVehicles = make([]*vehicles.Vehicle, 0, len(h.Vehicles))
	for _, v := range h.Vehicles {
		if !v.Broken || v.InUse {
			newVehicles = append(newVehicles, v)
		} else if v.T == vehicles.Cart {
			h.Resources.Add(IronBar, 1)
		}
	}
	h.Vehicles = newVehicles
}

func (h *Household) AddVehicle(v *vehicles.Vehicle) {
	h.Vehicles = append(h.Vehicles, v)
}

func (h *Household) AllocateVehicle(waterOk bool) *vehicles.Vehicle {
	for _, v := range h.Vehicles {
		if !v.InUse && (!v.T.Water || waterOk) {
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

func (h *Household) HasFreePerson() bool {
	for _, person := range h.People {
		if person.Task == nil {
			return true
		}
	}
	return false
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
	return h.Resources
}

func (h *Household) GetBuilding() *building.Building {
	return h.Building
}

func (h *Household) GetHeating() uint8 {
	return h.Heating
}

func (h *Household) Spend(amount uint32) {
	h.Money -= amount
}

func (h *Household) Earn(amount uint32) {
	h.Money += amount
}

func (h *Household) GetMoney() uint32 {
	return h.Money
}

func (h *Household) Destroy(m navigation.IMap) {
	dstH := h.Town.Townhall.Household
	for _, person := range h.People {
		dstH.AssignPerson(person, m)
	}
	dstH.Money += h.Money
	m.GetField(h.Building.X, h.Building.Y).Terrain.Resources.AddResources(*h.Resources)
}

func (h *Household) Destination(extensionType *building.BuildingExtensionType) navigation.Destination {
	return &navigation.BuildingDestination{B: h.Building, ET: extensionType}
}

func (h *Household) PendingCosts() uint32 {
	return PendingCosts(h.Tasks)
}

func (h *Household) Broken() bool {
	return h.Building.Broken
}

func (h *Household) GetTown() *Town {
	return h.Town
}

func (h *Household) GetPeople() []*Person {
	return h.People
}
