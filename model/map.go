package model

import (
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

func (m *Map) ShortPath(sx, sy, ex, ey uint16, travellerType uint8) *navigation.Path {
	if sx == ex && sy == ey {
		return nil
	}
	p := FindShortPathBFS(m, sx, sy, ex, ey, travellerType)
	if p != nil {
		return &navigation.Path{L: p[1:]}
	}
	return nil
}
