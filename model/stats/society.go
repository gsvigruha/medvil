package stats

import (
	"medvil/model/economy"
	"medvil/model/time"
	"reflect"
)

type SocietyStats struct {
	Deaths                 uint32
	Departures             uint32
	ProducedNum            uint32
	ExchangedNum           uint32
	ProducedPrice          uint32
	ExchangedPrice         uint32
	FoodPrice              uint32
	HouseholdItemPrice     uint32
	BuildingMaterialsPrice uint32
	PendingTasks           map[economy.Task]uint32
	CompletedTasks         map[string]uint32
}

func (s *SocietyStats) RegisterTrade(price uint32, quantity uint16) {
	s.ExchangedNum += uint32(quantity)
	s.ExchangedPrice += uint32(price)
}

func (s *SocietyStats) RegisterDeath() {
	s.Deaths++
}

func (s *SocietyStats) RegisterDeparture() {
	s.Departures++
}

func (s *SocietyStats) RegisterFoodPrices(items []uint32) {
	for _, item := range items {
		s.FoodPrice += item
	}
	s.FoodPrice /= uint32(len(items))
}

func (s *SocietyStats) RegisterBuildingMaterialsPrices(items []uint32) {
	for _, item := range items {
		s.BuildingMaterialsPrice += item
	}
	s.BuildingMaterialsPrice /= uint32(len(items))
}

func (s *SocietyStats) RegisterHouseholdItemPrices(items []uint32) {
	for _, item := range items {
		s.HouseholdItemPrice += item
	}
	s.HouseholdItemPrice /= uint32(len(items))
}

func (s *SocietyStats) StartTask(t economy.Task, calendar *time.CalendarType) {
	if s.PendingTasks != nil {
		s.PendingTasks[t] = calendar.DaysElapsed()
	}
}

func (s *SocietyStats) FinishTask(t economy.Task, calendar *time.CalendarType) {
	if s.PendingTasks != nil {
		if start, ok := s.PendingTasks[t]; ok {
			aggrName := reflect.TypeOf(t).Elem().Name()
			if aggrTime, ok := s.CompletedTasks[aggrName]; ok {
				s.CompletedTasks[aggrName] = aggrTime + calendar.DaysElapsed() - start
			} else {
				s.CompletedTasks[aggrName] = calendar.DaysElapsed() - start
			}
			delete(s.PendingTasks, t)
		}
	}
}
