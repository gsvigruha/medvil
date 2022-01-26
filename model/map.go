package model

import (
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/model/time"
)

type Map struct {
	SX        uint16
	SY        uint16
	Fields    [][]navigation.Field
	Countries []social.Country
}

func (m *Map) ElapseTime(Calendar *time.CalendarType) {
	for i := range m.Countries {
		country := m.Countries[i]
		for j := range country.Towns {
			town := country.Towns[j]
			for l := range town.Townhall.Household.People {
				person := town.Townhall.Household.People[l]
				person.ElapseTime(Calendar, m)
			}
			town.Townhall.ElapseTime(Calendar, m)
			for k := range town.Farms {
				farm := town.Farms[k]
				for l := range farm.Household.People {
					person := farm.Household.People[l]
					person.ElapseTime(Calendar, m)
				}
				farm.ElapseTime(Calendar, m)
			}
		}
	}
	for i := uint16(0); i < m.SX; i++ {
		for j := uint16(0); j < m.SY; j++ {
			f := m.Fields[i][j]
			if f.Plant != nil {
				f.Plant.ElapseTime(Calendar)
				if f.Plant.T.IsAnnual() && Calendar.Season() == time.Winter {
					f.Plant = nil
				}
			}
		}
	}
}

func (m *Map) GetField(x uint16, y uint16) *navigation.Field {
	return &m.Fields[x][y]
}

func (m *Map) ReverseReferences() *ReverseReferences {
	rr := BuildReverseReferences(m)
	return &rr
}

func (m *Map) AddFarm(c *social.Country, x, y uint16, bp *building.BuildingPlan) bool {
	b := m.AddBuilding(x, y, bp)
	if b != nil {
		t := c.Towns[0]
		f := &social.Farm{Household: social.Household{Building: b, Town: t}}
		t.Farms = append(t.Farms, f)
		return true
	} else {
		return false
	}
}

func (m *Map) AddBuilding(x, y uint16, bp *building.BuildingPlan) *building.Building {
	b := building.Building{X: x, Y: y, Plan: bp}
	for i := uint16(0); i < 5; i++ {
		for j := uint16(0); j < 5; j++ {
			bx := int(b.X+i) - 2
			by := int(b.Y+j) - 2
			if bx >= 0 && by >= 0 {
				if bp.BaseShape[i][j] && (!m.Fields[bx][by].Building.Empty() || m.Fields[bx][by].Plant != nil) {
					return nil
				}
			} else {
				return nil
			}
		}
	}
	for i := uint16(0); i < 5; i++ {
		for j := uint16(0); j < 5; j++ {
			bx := int(b.X+i) - 2
			by := int(b.Y+j) - 2
			if bp.BaseShape[i][j] {
				m.Fields[bx][by].Building.BuildingUnits = b.ToBuildingUnits(uint8(i), uint8(j))
				m.Fields[bx][by].Building.RoofUnit = b.GetRoof(uint8(i), uint8(j))
			}
		}
	}
	return &b
}

func (m *Map) ShortPath(sx, sy, ex, ey uint16, travellerType uint8) *navigation.Path {
	if sx == ex && sy == ey {
		return nil
	}
	p := FindShortPathBFS(m, sx, sy, ex, ey, travellerType)
	if p != nil {
		return &navigation.Path{F: p[1:]}
	}
	return nil
}
