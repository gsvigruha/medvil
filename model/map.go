package model

import (
	//"math/rand"
	"medvil/controller"
	"medvil/model/social"
)

type Map struct {
	SX        uint16
	SY        uint16
	Fields    [][]Field
	Countries []social.Country
}

func (m *Map) ElapseTime(Calendar *controller.CalendarType) {
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
