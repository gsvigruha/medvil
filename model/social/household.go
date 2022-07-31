package social

import (
	"math/rand"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/stats"
	"medvil/model/time"
	"medvil/model/vehicles"
	"strings"
)

const ReproductionRate = 1.0 / (24 * 30 * 12)
const StoragePerArea = 50
const HeatingBudgetRatio = 0.3
const ExtrasBudgetRatio = 0.2

var Log = artifacts.GetArtifact("log")
var Firewood = artifacts.GetArtifact("firewood")
var Tools = artifacts.GetArtifact("tools")

const LogToFirewood = 5

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

func (h *Household) getNextTask(e economy.Equipment) economy.Task {
	if len(h.Tasks) == 0 {
		return nil
	}
	var i = 0
	for i < len(h.Tasks) {
		t := h.Tasks[i]
		_, sok := t.(*economy.SellTask)
		_, bok := t.(*economy.BuyTask)
		if !sok && !bok && !t.Blocked() && !t.IsPaused() && t.Equipped(e) {
			break
		}
		i++
	}
	if i == len(h.Tasks) {
		return nil
	}
	t := h.Tasks[i]
	h.Tasks = append(h.Tasks[0:i], h.Tasks[i+1:]...)
	return t
}

func (h *Household) getExchangeTask(m navigation.IMap, vehicle *vehicles.Vehicle) *economy.ExchangeTask {
	mp := h.Town.Marketplace
	var maxVolume uint16 = ExchangeTaskMaxVolumePedestrian
	var buildingCheckFn = navigation.Field.BuildingNonExtension
	_, _, sailableMP := GetRandomBuildingXY(mp.Building, m, navigation.Field.Sailable)
	_, _, sailableH := GetRandomBuildingXY(h.Building, m, navigation.Field.Sailable)
	if vehicle != nil && sailableMP && sailableH {
		maxVolume = ExchangeTaskMaxVolumeBoat
		buildingCheckFn = navigation.Field.Sailable
	}

	mx, my, mok := GetRandomBuildingXY(mp.Building, m, buildingCheckFn)
	hx, hy, hok := GetRandomBuildingXY(h.Building, m, buildingCheckFn)
	if !hok || !mok {
		return nil
	}
	et := &economy.ExchangeTask{
		HomeF:          m.GetField(hx, hy),
		MarketF:        m.GetField(mx, my),
		Exchange:       mp,
		HouseholdR:     &h.Resources,
		HouseholdMoney: &h.Money,
		Vehicle:        vehicle,
		GoodsToBuy:     nil,
		GoodsToSell:    nil,
		TaskTag:        "",
	}
	var empty = true
	var tasks []economy.Task
	for _, ot := range h.Tasks {
		var combined = false
		bt, bok := ot.(*economy.BuyTask)
		if bok && !bt.Blocked() && !bt.IsPaused() && artifacts.GetVolume(et.GoodsToBuy) < maxVolume {
			et.AddBuyTask(bt)
			combined = true
		}
		st, sok := ot.(*economy.SellTask)
		if sok && !st.Blocked() && !st.IsPaused() && artifacts.GetVolume(et.GoodsToSell) < maxVolume {
			et.AddSellTask(st)
			combined = true
		}
		if !combined {
			tasks = append(tasks, ot)
		} else {
			empty = false
		}
	}
	if !empty {
		h.Tasks = tasks
		return et
	}
	return nil
}

func (h *Household) getNextTaskCombineExchange(m navigation.IMap, e economy.Equipment) economy.Task {
	vehicle := h.GetVehicle()
	et := h.getExchangeTask(m, vehicle)
	if et == nil && vehicle != nil {
		vehicle.SetInUse(false)
	}
	if et != nil {
		return et
	}
	return h.getNextTask(e)
}

func (h *Household) AddTask(t economy.Task) {
	h.Tasks = append(h.Tasks, t)
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
	water := artifacts.GetArtifact("water")
	if h.Resources.Get(water) < economy.MinFoodOrDrinkPerPerson*numP &&
		NumBatchesSimple(economy.MaxFoodOrDrinkPerPerson*numP, WaterTransportQuantity) > h.NumTasks("transport", "water") {
		hx, hy, ok := GetRandomBuildingXY(h.Building, m, navigation.Field.BuildingNonExtension)
		if ok {
			dest := m.FindDest(navigation.Location{X: hx, Y: hy, Z: 0}, economy.WaterDestination{}, navigation.TravellerTypePedestrian)
			if dest != nil {
				h.AddTask(&economy.TransportTask{
					PickupF:  dest,
					DropoffF: m.GetField(hx, hy),
					PickupR:  &dest.Terrain.Resources,
					DropoffR: &h.Resources,
					A:        water,
					Quantity: WaterTransportQuantity,
				})
			}
		}
	}
	mp := h.Town.Marketplace
	var numFoodBatchesNeeded = 0
	for _, a := range economy.Foods {
		if h.Resources.Get(a) < economy.MinFoodOrDrinkPerPerson*numP {
			numFoodBatchesNeeded += NumBatchesSimple(economy.BuyFoodOrDrinkPerPerson()*numP, FoodTransportQuantity)
		}
	}
	for _, a := range economy.Foods {
		if h.Resources.Get(a) < economy.MinFoodOrDrinkPerPerson*numP {
			tag := "food_shopping#" + a.Name
			if NumBatchesSimple(economy.BuyFoodOrDrinkPerPerson()*numP, FoodTransportQuantity) > h.NumTasks("exchange", tag) {
				needs := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: FoodTransportQuantity}}
				var maxPrice = h.Money / uint32(numFoodBatchesNeeded)
				if maxPrice > mp.Price(needs)*2 {
					maxPrice = mp.Price(needs) * 2
				}
				if h.Money >= mp.Price(needs) && mp.HasTraded(a) {
					h.AddTask(&economy.BuyTask{
						Exchange:       mp,
						HouseholdMoney: &h.Money,
						Goods:          needs,
						MaxPrice:       maxPrice,
						TaskTag:        tag,
					})
				}
			}
		}
	}
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
	if h.Resources.Get(Firewood) < h.heatingFuelNeeded() && h.Resources.Get(Log)*LogToFirewood < h.heatingFuelNeeded() {
		tag := "heating_fuel_shopping"
		if NumBatchesSimple(ProductTransportQuantity(Log), ProductTransportQuantity(Log)) > h.NumTasks("exchange", tag) {
			needs := []artifacts.Artifacts{artifacts.Artifacts{A: Log, Quantity: ProductTransportQuantity(Log)}}
			if h.Money >= mp.Price(needs) && mp.HasTraded(Log) {
				h.AddTask(&economy.BuyTask{
					Exchange:       mp,
					HouseholdMoney: &h.Money,
					Goods:          needs,
					MaxPrice:       uint32(float64(h.Money) * HeatingBudgetRatio),
					TaskTag:        tag,
				})
			}
		}
	}
	if h.Resources.Get(economy.Medicine) < numP {
		tag := "medicine_shopping"
		if NumBatchesSimple(ProductTransportQuantity(economy.Medicine), ProductTransportQuantity(economy.Medicine)) > h.NumTasks("exchange", tag) {
			needs := []artifacts.Artifacts{artifacts.Artifacts{A: economy.Medicine, Quantity: ProductTransportQuantity(economy.Medicine)}}
			if h.Money >= mp.Price(needs) && mp.HasTraded(economy.Medicine) {
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
	if h.Resources.Get(economy.Beer) < numP {
		tag := "beer_shopping"
		if NumBatchesSimple(ProductTransportQuantity(economy.Beer), ProductTransportQuantity(economy.Beer)) > h.NumTasks("exchange", tag) {
			needs := []artifacts.Artifacts{artifacts.Artifacts{A: economy.Beer, Quantity: ProductTransportQuantity(economy.Beer)}}
			if h.Money >= mp.Price(needs) && mp.HasTraded(economy.Beer) {
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
	if Calendar.Hour == 0 {
		for i := 0; i < len(h.Tasks); i++ {
			if h.Tasks[i].IsPaused() {
				h.Tasks[i].Pause(false)
			}
		}
	}
	if h.Resources.Get(Firewood) < h.heatingFuelNeededPerMonth() && h.Resources.Remove(Log, 1) > 0 {
		h.Resources.Add(Firewood, LogToFirewood)
	}
	if Calendar.Day == 1 && Calendar.Hour == 0 {
		if Calendar.Season() == time.Winter {
			wood := h.Resources.Remove(Firewood, h.heatingFuelNeededPerMonth())
			h.Heating = float64(wood) / float64(h.heatingFuelNeededPerMonth())
		} else {
			h.Heating = 1.0
		}
	}
}

func (h *Household) MaybeBuyBoat(Calendar *time.CalendarType, m navigation.IMap) {
	if h.numBoats() == 0 && h.Building.HasExtension(building.Deck) && h.NumTasks("factory_pickup", economy.BoatConstruction.Name) == 0 {
		factory := PickFactory(h.Town.Factories)
		if factory != nil && factory.Price(economy.BoatConstruction) < uint32(float64(h.Money)*ExtrasBudgetRatio) {
			ext, hx, hy := h.Building.GetExtensionWithCoords()
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

func (h *Household) numBoats() int {
	var n = 0
	for _, v := range h.Vehicles {
		if v.T == vehicles.Boat {
			n++
		}
	}
	return n
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

func (h *Household) ArtifactToSell(a *artifacts.Artifact, q uint16, isProduct bool) uint16 {
	if a.Name == "water" || a == Firewood {
		return 0
	}
	if a == Tools && !isProduct {
		return 0
	}
	var threshold = economy.MaxFoodOrDrinkPerPerson
	if isProduct {
		threshold = economy.ProductMaxFoodOrDrinkPerPerson
	}
	var result uint16
	if economy.IsFoodOrDrink(a) {
		if q > threshold*uint16(len(h.People)) {
			result = q - threshold*uint16(len(h.People))
		} else {
			return 0
		}
	} else {
		result = q
	}
	if a == Log {
		logs := h.heatingFuelNeeded()/LogToFirewood + 1
		if q > logs {
			result = q - logs
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

func countTags(task economy.Task, name, tag string) int {
	var i = 0
	taskTags := strings.Split(task.Tag(), ";")
	for _, taskTag := range taskTags {
		if task.Name() == name && strings.Contains(taskTag, tag) {
			i++
		}
	}
	return i
}

func (h *Household) NumTasks(name string, tag string) int {
	var i = 0
	for _, t := range h.Tasks {
		i += countTags(t, name, tag)
	}
	for _, p := range h.People {
		if p.Task != nil {
			i += countTags(p.Task, name, tag)
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
		Household: h,
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
			person.Household = dstH
			dstH.People = append(dstH.People, person)
			person.Task = &economy.GoHomeTask{F: m.GetField(dstH.Building.X, dstH.Building.Y), P: person}
			break
		}
	}
}
