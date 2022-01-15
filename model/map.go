package model

import (
	//"math/rand"
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
				farm.ElapseTime(Calendar)
				for l := range farm.Household.People {
					person := farm.Household.People[l]
					person.ElapseTime(Calendar, m)
				}
			}
		}
	}
}

func (m *Map) GetField(x uint16, y uint16) *navigation.Field {
	return &m.Fields[x][y]
}
