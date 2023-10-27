package stats

import (
	"medvil/model/artifacts"
	"medvil/model/economy"
	"medvil/model/time"
	"reflect"
)

type HouseholdStats struct {
	Money     uint32
	Artifacts uint32
	People    uint32
	Buildings uint32
}

type Stats struct {
	Global   HouseholdStats
	Farm     HouseholdStats
	Workshop HouseholdStats
	Mine     HouseholdStats
	Gov      HouseholdStats
	Trader   HouseholdStats

	Deaths     uint32
	Departures uint32
	Poverty    uint32
	TradeM     map[*artifacts.Artifact]uint32
	TradeQ     map[*artifacts.Artifact]uint32
	PendingT   map[economy.Task]uint32
	CompletedT map[string]uint32
}

func (s *Stats) Init(pt map[economy.Task]uint32) {
	s.PendingT = pt
	s.CompletedT = make(map[string]uint32)
	s.TradeM = make(map[*artifacts.Artifact]uint32)
	s.TradeQ = make(map[*artifacts.Artifact]uint32)
	for _, a := range artifacts.All {
		s.TradeM[a] = 0
		s.TradeQ[a] = 0
	}
}

func (s *Stats) Reset() {
	s.Global.Reset()
	s.Farm.Reset()
	s.Workshop.Reset()
	s.Mine.Reset()
	s.Gov.Reset()
	s.Trader.Reset()
	s.Poverty = 0
}

func (s *HouseholdStats) Reset() {
	s.Money = 0
	s.Artifacts = 0
	s.People = 0
	s.Buildings = 0
}

func (s *HouseholdStats) Add(os *HouseholdStats) {
	s.Money += os.Money
	s.Artifacts += os.Artifacts
	s.People += os.People
	s.Buildings += os.Buildings
}

func (s *Stats) Add(os *Stats) {
	s.Global.Add(&os.Global)
	s.Farm.Add(&os.Farm)
	s.Workshop.Add(&os.Workshop)
	s.Mine.Add(&os.Mine)
	s.Gov.Add(&os.Gov)
	s.Trader.Add(&os.Trader)

	s.Deaths += os.Deaths
	s.Departures += os.Departures
	s.Poverty += os.Poverty

	for a, q := range os.TradeM {
		s.TradeM[a] = s.TradeM[a] + q
	}
	for a, q := range os.TradeQ {
		s.TradeQ[a] = s.TradeQ[a] + q
	}
	for t, q := range os.CompletedT {
		s.CompletedT[t] = s.CompletedT[t] + q
	}
}

func (s *Stats) RegisterDeath() {
	s.Deaths++
}

func (s *Stats) RegisterDeparture() {
	s.Departures++
}

func (s *Stats) RegisterTrade(a *artifacts.Artifact, unitPrice uint32, quantity uint16) {
	s.TradeM[a] = s.TradeM[a] + unitPrice*uint32(quantity)
	s.TradeQ[a] = s.TradeQ[a] + uint32(quantity)
}

func (s *Stats) StartTask(t economy.Task, calendar *time.CalendarType) {
	if s.PendingT != nil {
		s.PendingT[t] = calendar.DaysElapsed()
	}
}

func (s *Stats) DeleteTask(t economy.Task) {
	if s.PendingT != nil {
		if _, ok := s.PendingT[t]; ok {
			delete(s.PendingT, t)
		}
	}
}

func (s *Stats) FinishTask(t economy.Task, calendar *time.CalendarType) {
	if s.PendingT != nil {
		if start, ok := s.PendingT[t]; ok {
			aggrName := reflect.TypeOf(t).Elem().Name()
			if aggrTime, ok := s.CompletedT[aggrName]; ok {
				s.CompletedT[aggrName] = aggrTime + calendar.DaysElapsed() - start
			} else {
				s.CompletedT[aggrName] = calendar.DaysElapsed() - start
			}
			delete(s.PendingT, t)
		}
	}
}
