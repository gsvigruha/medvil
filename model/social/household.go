package social

import (
	"math/rand"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

const ReproductionRate = 1.0 / (24 * 30 * 12)
const StoragePerArea = 10

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
	for i := range h.People {
		person := h.People[i]
		if person.Task == nil && h.HasTask() {
			person.Task = h.getNextTask()
		}
	}
	if h.Town != nil { // Not Townhall, needs better check
		if h.HasRoomForPeople() && len(h.Town.Townhall.Household.People) > 0 {
			person := h.Town.Townhall.Household.People[0]
			h.Town.Townhall.Household.People = h.Town.Townhall.Household.People[1:]
			h.People = append(h.People, person)
			person.Household = h
			person.Task = &economy.GoHomeTask{F: m.GetField(h.Building.X, h.Building.Y), P: person}
		}
		if rand.Float64() < ReproductionRate {
			if len(h.People) >= 1 && h.HasRoomForPeople() {
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
}

func (h *Household) ArtifactToSell(a *artifacts.Artifact, q uint16) uint16 {
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

func (h *Household) NewPerson() *Person {
	return &Person{
		Food:      MaxPersonState,
		Water:     MaxPersonState,
		Happiness: MaxPersonState,
		Health:    MaxPersonState,
		Household: h,
		Task:      nil,
		IsHome:    true,
		Traveller: &navigation.Traveller{
			FX: h.Building.X,
			FY: h.Building.Y,
			FZ: 0,
			PX: 0,
			PY: 0,
		},
	}
}
