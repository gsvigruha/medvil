package stats

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/time"
	"reflect"
)

type Stats struct {
	Money     uint32
	Artifacts uint32
	People    uint32
	Buildings uint32

	Deaths           uint32
	Departures       uint32
	TradeMoneyAmount map[*artifacts.Artifact]uint32
	TradeQuantity    map[*artifacts.Artifact]uint32
	PendingTasks     map[economy.Task]uint32
	CompletedTasks   map[string]uint32
}

func (s *Stats) Init(pt map[economy.Task]uint32) {
	s.PendingTasks = pt
	s.CompletedTasks = make(map[string]uint32)
	s.TradeMoneyAmount = make(map[*artifacts.Artifact]uint32)
	s.TradeQuantity = make(map[*artifacts.Artifact]uint32)
	for _, a := range artifacts.All {
		s.TradeMoneyAmount[a] = 0
		s.TradeQuantity[a] = 0
	}
}

func (s *Stats) Reset() {
	s.Money = 0
	s.Artifacts = 0
	s.People = 0
	s.Buildings = 0
}

func (s *Stats) Add(os *Stats) {
	s.Money += os.Money
	s.Artifacts += os.Artifacts
	s.People += os.People
	s.Buildings += os.Buildings

	s.Deaths += os.Deaths
	s.Departures += os.Departures

	for a, q := range os.TradeMoneyAmount {
		s.TradeMoneyAmount[a] = s.TradeMoneyAmount[a] + q
	}
	for a, q := range os.TradeQuantity {
		s.TradeQuantity[a] = s.TradeQuantity[a] + q
	}
	for t, q := range os.CompletedTasks {
		s.CompletedTasks[t] = s.CompletedTasks[t] + q
	}
}

func (s *Stats) RegisterDeath() {
	s.Deaths++
}

func (s *Stats) RegisterDeparture() {
	s.Departures++
}

func (s *Stats) RegisterTrade(a *artifacts.Artifact, unitPrice uint32, quantity uint16) {
	s.TradeMoneyAmount[a] = s.TradeMoneyAmount[a] + unitPrice*uint32(quantity)
	s.TradeQuantity[a] = s.TradeQuantity[a] + uint32(quantity)
}

func (s *Stats) StartTask(t economy.Task, calendar *time.CalendarType) {
	if s.PendingTasks != nil {
		s.PendingTasks[t] = calendar.DaysElapsed()
	}
}

func (s *Stats) FinishTask(t economy.Task, calendar *time.CalendarType) {
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
