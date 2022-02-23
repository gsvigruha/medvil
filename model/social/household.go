package social

import (
	"math/rand"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
	"strings"
)

const ReproductionRate = 1.0 / (24 * 30 * 12)
const StoragePerArea = 25
const waterTransportQuantity = 10
const foodTransportQuantity = 5

type Household struct {
	People          []*Person
	TargetNumPeople uint16
	Money           uint32
	Building        *building.Building
	Town            *Town
	Tasks           []economy.Task
	Resources       artifacts.Resources
}

func (h *Household) HasTask() bool {
	for i := range h.Tasks {
		if !h.Tasks[i].Blocked() {
			return true
		}
	}
	return false
}

func (h *Household) getNextTask() economy.Task {
	var i = 0
	for i < len(h.Tasks) {
		if !h.Tasks[i].Blocked() {
			break
		}
		i++
	}
	t := h.Tasks[i]
	h.Tasks = append(h.Tasks[0:i], h.Tasks[i+1:]...)
	return t
}

func (h *Household) AddTask(t economy.Task) {
	h.Tasks = append(h.Tasks, t)
}

func (h *Household) IncTargetNumPeople() {
	if h.TargetNumPeople < h.Building.Plan.Area()*2 {
		h.TargetNumPeople++
	}
}

func (h *Household) HasRoomForPeople() bool {
	return uint16(len(h.People)) < h.TargetNumPeople
}

func (h *Household) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	if &h.Town.Townhall.Household != h { // Not Townhall, needs better check
		if h.HasRoomForPeople() && len(h.Town.Townhall.Household.People) > 0 {
			person := h.Town.Townhall.Household.People[0]
			h.Town.Townhall.Household.People = h.Town.Townhall.Household.People[1:]
			h.People = append(h.People, person)
			person.Household = h
			person.Task = &economy.GoHomeTask{F: m.GetField(h.Building.X, h.Building.Y), P: person}
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
	}
	numP := uint16(len(h.People))
	water := artifacts.GetArtifact("water")
	if h.Resources.Get(water) < economy.MinFoodOrDrinkPerPerson*numP {
		hx, hy := h.Building.GetRandomBuildingXY()
		dest := m.FindDest(hx, hy, economy.WaterDestination{}, navigation.TravellerTypePedestrian)
		if dest != nil {
			if int(economy.MaxFoodOrDrinkPerPerson*numP)/waterTransportQuantity > h.NumTasks("transport", "water") {
				h.AddTask(&economy.TransportTask{
					PickupF:  dest,
					DropoffF: m.GetField(hx, hy),
					PickupR:  &dest.Terrain.Resources,
					DropoffR: &h.Resources,
					A:        water,
					Quantity: waterTransportQuantity,
				})
			}
		}
	}
	mp := h.Town.Marketplace
	for _, a := range economy.Foods {
		if h.Resources.Get(a) < economy.MinFoodOrDrinkPerPerson*numP {
			tag := "food_shopping#" + a.Name
			if int(economy.BuyFoodOrDrinkPerPerson()*numP)/foodTransportQuantity > h.NumTasks("exchange", tag) {
				needs := []artifacts.Artifacts{artifacts.Artifacts{A: a, Quantity: foodTransportQuantity}}
				if h.Money >= mp.Price(needs) && mp.CanBuy(needs) {
					mx, my := mp.Building.GetRandomBuildingXY()
					hx, hy := h.Building.GetRandomBuildingXY()
					h.AddTask(&economy.ExchangeTask{
						HomeF:          m.GetField(hx, hy),
						MarketF:        m.GetField(mx, my),
						Exchange:       mp,
						HouseholdR:     &h.Resources,
						HouseholdMoney: &h.Money,
						GoodsToBuy:     needs,
						GoodsToSell:    nil,
						TaskTag:        tag,
					})
				}
			}
		}
	}
}

func (h *Household) ArtifactToSell(a *artifacts.Artifact, q uint16) uint16 {
	if a.Name == "water" {
		return 0
	}
	if economy.IsFoodOrDrink(a) {
		if q > economy.MinFoodOrDrinkPerPerson*uint16(len(h.People)) {
			return q - economy.MinFoodOrDrinkPerPerson*uint16(len(h.People))
		} else {
			return 0
		}
	}
	return q
}

func (h *Household) HasFood() bool {
	return economy.HasFood(h.Resources)
}

func (h *Household) HasDrink() bool {
	return economy.HasDrink(h.Resources)
}

func (h *Household) NumTasks(name string, tag string) int {
	var i = 0
	for _, t := range h.Tasks {
		if t.Name() == name && strings.Contains(t.Tag(), tag) {
			i++
		}
	}
	for _, p := range h.People {
		if p.Task != nil && p.Task.Name() == name && strings.Contains(p.Task.Tag(), tag) {
			i++
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

func (h *Household) FilterPeople(m navigation.IMap) {
	var newPeople = make([]*Person, 0, len(h.People))
	for _, p := range h.People {
		if p.Health == 0 {
			m.GetField(p.Traveller.FX, p.Traveller.FY).UnregisterTraveller(p.Traveller)			
		} else {
			newPeople = append(newPeople, p)
		}
	}
	h.People = newPeople
}
