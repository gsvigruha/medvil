package social

import (
	"math/rand"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/stats"
	"medvil/model/time"
	"strings"
)

const ReproductionRate = 1.0 / (24 * 30 * 12)
const StoragePerArea = 50
const ToolsBudgetRatio = 0.2
const HeatingBudgetRatio = 0.3

var Log = artifacts.GetArtifact("log")
var Tools = artifacts.GetArtifact("tools")

type Household struct {
	People          []*Person
	TargetNumPeople uint16
	Money           uint32
	Building        *building.Building
	Town            *Town
	Tasks           []economy.Task
	Resources       artifacts.Resources
	Heating         float64
}

func (h *Household) HasTask() bool {
	for i := range h.Tasks {
		if !h.Tasks[i].Blocked() && !h.Tasks[i].IsPaused() {
			return true
		}
	}
	return false
}

func (h *Household) getNextTask() economy.Task {
	var i = 0
	for i < len(h.Tasks) {
		if !h.Tasks[i].Blocked() && !h.Tasks[i].IsPaused() {
			break
		}
		i++
	}
	t := h.Tasks[i]
	h.Tasks = append(h.Tasks[0:i], h.Tasks[i+1:]...)
	return t
}

func (h *Household) getExchangeTask(m navigation.IMap) *economy.ExchangeTask {
	mp := h.Town.Marketplace
	mx, my := mp.Building.GetRandomBuildingXY()
	et := &economy.ExchangeTask{
		HomeF:          m.GetField(h.Building.X, h.Building.Y),
		MarketF:        m.GetField(mx, my),
		Exchange:       mp,
		HouseholdR:     &h.Resources,
		HouseholdMoney: &h.Money,
		GoodsToBuy:     nil,
		GoodsToSell:    nil,
		TaskTag:        "",
	}
	var empty = true
	var tasks []economy.Task
	for _, ot := range h.Tasks {
		var combined = false
		bt, bok := ot.(*economy.BuyTask)
		if bok && !bt.Blocked() && !bt.IsPaused() && artifacts.GetVolume(et.GoodsToBuy) < ExchangeTaskMaxVolume {
			et.AddBuyTask(bt)
			combined = true
		}
		st, sok := ot.(*economy.SellTask)
		if sok && !st.Blocked() && !st.IsPaused() && artifacts.GetVolume(et.GoodsToSell) < ExchangeTaskMaxVolume {
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

func (h *Household) getNextTaskCombineExchange(m navigation.IMap) economy.Task {
	et := h.getExchangeTask(m)
	if et != nil {
		return et
	}
	return h.getNextTask()
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
			for pi, person := range h.Town.Townhall.Household.People {
				if person.Task == nil {
					person.Reassign(h, pi, m)
					break
				}
			}
		}
		if len(h.People) >= 2 && rand.Float64() < ReproductionRate {
			if h.HasRoomForPeople() {
				h.People = append(h.People, h.NewPerson())
			} else if h.Town.Townhall.Household.HasRoomForPeople() {
				person := h.Town.Townhall.Household.NewPerson()
				h.Town.Townhall.Household.People = append(h.Town.Townhall.Household.People, person)
				person.Traveller.FX = h.Building.X
				person.Traveller.FY = h.Building.Y
				person.Task = &economy.GoHomeTask{F: m.GetField(h.Town.Townhall.Household.Building.X, h.Town.Townhall.Household.Building.Y), P: person}
			}
		}
		if h.HasSurplusPeople() && h.Town.Townhall.Household.HasRoomForPeople() {
			for pi, person := range h.People {
				if person.Task == nil {
					person.Reassign(&h.Town.Townhall.Household, pi, m)
					break
				}
			}
		}
	}
	numP := uint16(len(h.People))
	water := artifacts.GetArtifact("water")
	if h.Resources.Get(water) < economy.MinFoodOrDrinkPerPerson*numP &&
		NumBatchesSimple(economy.MaxFoodOrDrinkPerPerson*numP, WaterTransportQuantity) > h.NumTasks("transport", "water") {
		hx, hy := h.Building.GetRandomBuildingXY()
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
				MaxPrice:       uint32(float64(h.Money) * ToolsBudgetRatio),
				TaskTag:        "tools_purchase",
			})
		}
	}
	if h.Resources.Get(Log) < h.heatingFuelNeeded() {
		tag := "heating_fuel_shopping"
		if NumBatchesSimple(h.heatingFuelNeeded(), ProductTransportQuantity(Log)) > h.NumTasks("exchange", tag) {
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
	if Calendar.Hour == 0 {
		for i := 0; i < len(h.Tasks); i++ {
			if h.Tasks[i].IsPaused() {
				h.Tasks[i].Pause(false)
			}
		}
	}
	if Calendar.Day == 1 && Calendar.Hour == 0 {
		if Calendar.Season() == time.Winter {
			logs := h.Resources.Remove(Log, h.heatingFuelNeededPerMonth())
			h.Heating = float64(logs) / float64(h.heatingFuelNeededPerMonth())
		} else {
			h.Heating = 1.0
		}
	}
}

func (h *Household) heatingFuelNeededPerMonth() uint16 {
	fuel := uint16(len(h.People) / 5)
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
		if p.Tool {
			n++
		}
	}
	return n
}

func (h *Household) ArtifactToSell(a *artifacts.Artifact, q uint16, isProduct bool) uint16 {
	if a.Name == "water" {
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
		if q > h.heatingFuelNeeded() {
			result = q - h.heatingFuelNeeded()
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

func (h *Household) NewPerson() *Person {
	hx, hy := h.Building.GetRandomBuildingXY()
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
		},
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

func (h *Household) Stats() *stats.Stats {
	return &stats.Stats{
		Money:     h.Money,
		People:    uint32(len(h.People)),
		Buildings: 1,
		Artifacts: h.Resources.NumArtifacts(),
	}
}
