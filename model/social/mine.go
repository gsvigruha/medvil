package social

import (
	"encoding/json"
	"math/rand"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/terrain"
	"medvil/model/time"
	"medvil/util"
)

const MineMaxDistance = 15

type MineLand struct {
	X       uint16
	Y       uint16
	UseType uint8
	F       *navigation.Field
}

func (l MineLand) Field() *navigation.Field {
	return l.F
}

func (l MineLand) Context() string {
	switch l.UseType {
	case economy.MineFieldUseTypeStone:
		return "stone"
	case economy.MineFieldUseTypeClay:
		return "clay"
	case economy.MineFieldUseTypeIron:
		return "iron_ore"
	case economy.MineFieldUseTypeGold:
		return "gold_ore"
	}
	return ""
}

type Mine struct {
	Household  *Household
	Land       []MineLand
	AutoSwitch bool
	Optimize   bool
}

func (m *Mine) UnmarshalJSON(data []byte) error {
	var j map[string]json.RawMessage
	if err := json.Unmarshal(data, &j); err != nil {
		return err
	}
	if err := json.Unmarshal(j["household"], &m.Household); err != nil {
		return err
	}
	var l [][]uint16
	if err := json.Unmarshal(j["land"], &l); err != nil {
		return err
	}
	m.Land = make([]MineLand, len(l))
	for i := range l {
		m.Land[i].X = l[i][0]
		m.Land[i].Y = l[i][1]
		m.Land[i].UseType = uint8(l[i][2])
	}
	return nil
}

func (m *Mine) AddTransportTask(l MineLand, imap navigation.IMap) {
	home := imap.GetField(m.Household.Building.X, m.Household.Building.Y)
	for a, q := range l.F.Terrain.Resources.Artifacts {
		if l.F.Terrain.Resources.IsRealArtifact(a) && q > 0 {
			tag := economy.TransportTaskTag(l.F, a)
			if m.Household.NumTasks("transport", tag) == 0 {
				m.Household.AddTask(&economy.TransportTask{
					PickupD:        l.F,
					DropoffD:       home,
					PickupR:        l.F.Terrain.Resources,
					DropoffR:       m.Household.Resources,
					A:              a,
					TargetQuantity: q,
				})
			}
		}
	}
}

func CheckMineUseType(useType uint8, f *navigation.Field) bool {
	if f.Deposit == nil {
		return false
	}
	if useType == economy.MineFieldUseTypeStone && f.Deposit.T == terrain.Rock {
		return true
	}
	if useType == economy.MineFieldUseTypeClay && f.Deposit.T == terrain.Mud {
		return true
	}
	if useType == economy.MineFieldUseTypeIron && f.Deposit.T == terrain.IronBog {
		return true
	}
	if useType == economy.MineFieldUseTypeGold && f.Deposit.T == terrain.Gold {
		return true
	}
	return false
}

func (m *Mine) ElapseTime(Calendar *time.CalendarType, imap navigation.IMap) {
	m.Household.ElapseTime(Calendar, imap)

	if Calendar.Day == 30 && Calendar.Hour == 0 && Calendar.Month == 12 {
		m.Optimize = m.AutoSwitch && m.Household.Resources.Remove(Paper, 1) > 0
	}

	if Calendar.Hour == 0 && len(m.Land) > 0 {
		for _, land := range m.Land {
			m.AddTransportTask(land, imap)
		}
		numP := len(m.Household.People)
		if m.Household.NumTasks("mining", economy.EmptyTag) < numP {
			if m.Optimize && m.Household.Town.Marketplace != nil {
				var profits []float64
				for _, land := range m.Land {
					profits = append(profits, float64(m.Household.Town.Marketplace.Prices[economy.MineUseTypeArtifact(land.UseType)]))
				}
				land := m.Land[util.RandomIndexWeighted(profits)]
				m.Household.AddTask(&economy.MiningTask{F: land.F, UseType: land.UseType})
			} else {
				land := m.Land[rand.Intn(len(m.Land))]
				m.Household.AddTask(&economy.MiningTask{F: land.F, UseType: land.UseType})
			}
		}
	}

	if m.Household.Town.Marketplace != nil {
		m.Household.SellArtifacts(NotInputOrProduct, NotInputOrProduct)
		m.Household.MaybeBuyPaper(m.AutoSwitch)
		m.Household.MaybeBuyBoat(Calendar, imap)
		m.Household.MaybeBuyCart(Calendar, imap)
	}
}

func (m *Mine) GetFields() []navigation.FieldWithContext {
	fields := make([]navigation.FieldWithContext, len(m.Land))
	for i := range m.Land {
		fields[i] = m.Land[i]
	}
	return fields
}

func (m *Mine) FieldWithinDistance(field *navigation.Field) bool {
	return WithinDistance(m.Household.Building, field, MineMaxDistance)
}

func (m *Mine) GetHome() Home {
	return m.Household
}

func (m *Mine) GetLandDistribution() map[uint8]int {
	result := make(map[uint8]int)
	for _, land := range m.Land {
		if cnt, ok := result[land.UseType]; ok {
			result[land.UseType] = cnt + 1
		} else {
			result[land.UseType] = 1
		}
	}
	return result
}
