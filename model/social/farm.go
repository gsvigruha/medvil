package social

import (
	"encoding/json"
	"medvil/controller"
	"medvil/model/economy"
	"medvil/model/navigation"
)

type FarmLand struct {
	X uint16
	Y uint16
	F navigation.IField
}

type Farm struct {
	Household Household
	Land      []FarmLand
	Tasks     []*economy.AgriculturalTask
}

func (f *Farm) UnmarshalJSON(data []byte) error {
	var j map[string]json.RawMessage
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	if err := json.Unmarshal(j["household"], &f.Household); err != nil {
		return err
	}
	var l [][]uint16
	if err := json.Unmarshal(j["land"], &l); err != nil {
		return err
	}
	f.Land = make([]FarmLand, len(l))
	for i := range l {
		f.Land[i].X = l[i][0]
		f.Land[i].Y = l[i][1]
	}
	return nil
}

func (f *Farm) HasTask() bool {
	return len(f.Tasks) > 0
}

func (f *Farm) getNextTask() economy.Task {
	t := f.Tasks[0]
	f.Tasks = f.Tasks[1:]
	return t
}

func (f *Farm) ElapseTime(Calendar *controller.CalendarType) {
	for i := range f.Household.People {
		person := f.Household.People[i]
		if person.Task == nil && person.IsHome && f.HasTask() {
			person.Task = f.getNextTask()
		}
	}
}
