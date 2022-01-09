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
	for i := uint16(0); i < m.SX; i++ {
		for j := uint16(0); j < m.SY; j++ {
			//f := m.Fields[i][j]

		}
	}
}
