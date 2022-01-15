package social

import (
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/time"
)

type Household struct {
	People   []*Person
	Money    uint32
	Building *building.Building
	Town     *Town
	Tasks    []economy.Task
}

func (h *Household) HasTask() bool {
	return len(h.Tasks) > 0
}

func (h *Household) getNextTask() economy.Task {
	t := h.Tasks[0]
	h.Tasks = h.Tasks[1:]
	return t
}

func (h *Household) ElapseTime(Calendar *time.CalendarType) {
	for i := range h.People {
		person := h.People[i]
		if person.Task == nil && person.IsHome && h.HasTask() {
			person.Task = h.getNextTask()
		}
	}
}
