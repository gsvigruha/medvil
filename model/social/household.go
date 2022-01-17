package social

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
	//"fmt"
)

type Household struct {
	People    []*Person
	Money     uint32
	Building  *building.Building
	Town      *Town
	Tasks     []economy.Task
	Resources artifacts.Resources
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

func (h *Household) ElapseTime(Calendar *time.CalendarType) {
	for i := range h.People {
		person := h.People[i]
		if person.Task == nil && h.HasTask() {
			person.Task = h.getNextTask()
		}
	}
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
