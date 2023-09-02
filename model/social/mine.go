package social

import (
	"encoding/json"
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/terrain"
	"medvil/model/time"
)

const MineMaxDistance = 10

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
	Household *Household
	Land      []MineLand
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
	if useType == economy.MineFieldUseTypeStone && f.Terrain.T == terrain.Rock {
		return true
	}
	if useType == economy.MineFieldUseTypeClay && f.Terrain.T == terrain.Mud {
		return true
	}
	if useType == economy.MineFieldUseTypeIron && f.Terrain.T == terrain.IronBog {
		return true
	}
	if useType == economy.MineFieldUseTypeGold && f.Terrain.T == terrain.Gold {
		return true
	}
	return false
}

func (m *Mine) ElapseTime(Calendar *time.CalendarType, imap navigation.IMap) {
	m.Household.ElapseTime(Calendar, imap)
	if Calendar.Hour == 0 {
		numP := len(m.Household.People)
		for _, land := range m.Land {
			m.AddTransportTask(land, imap)
			tag := economy.MiningTaskTag(land.F, land.UseType)
			if m.Household.NumTasks("mining", tag) < numP {
				if CheckMineUseType(land.UseType, land.F) && land.F.Terrain.Resources.IsEmpty() {
					m.Household.AddTask(&economy.MiningTask{F: land.F, UseType: land.UseType})
				}
			}
		}
	}

	if m.Household.Town.Marketplace != nil {
		m.Household.SellArtifacts(NotInputOrProduct, NotInputOrProduct)
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
