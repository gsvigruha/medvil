package model

import (
	"math/rand"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/model/terrain"
	"medvil/model/time"
)

type Map struct {
	SX        uint16
	SY        uint16
	Fields    [][]navigation.Field
	Countries []social.Country
}

func (m *Map) SpreadPlant(i, j uint16) {
	if i >= 0 && j >= 0 && i < m.SX && j < m.SY {
		if m.Fields[i][j].Plant == nil && m.Fields[i][j].Terrain.T == terrain.Water {
			m.Fields[i][j].Plant = &terrain.Plant{
				T:             &terrain.AllCropTypes[2],
				X:             uint16(i),
				Y:             uint16(j),
				BirthDateDays: uint32(1000*12*30 - rand.Intn(20*12*30)),
				Shape:         uint8(rand.Intn(10)),
			}
		}
	}
}

func (m *Map) ElapseTime(Calendar *time.CalendarType) {
	for i := range m.Countries {
		country := m.Countries[i]
		for j := range country.Towns {
			country.Towns[j].ElapseTime(Calendar, m)
		}
	}
	for i := uint16(0); i < m.SX; i++ {
		for j := uint16(0); j < m.SY; j++ {
			f := &m.Fields[i][j]
			if f.Plant != nil {
				f.Plant.ElapseTime(Calendar)
				if f.Plant.T.Name == "reed" && rand.Float64() < 0.001 {
					m.SpreadPlant(i-1, j)
					m.SpreadPlant(i, j-1)
					m.SpreadPlant(i+1, j)
					m.SpreadPlant(i, j+1)
				}
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

func (m *Map) AddBuildingConstruction(c *social.Country, x, y uint16, bp *building.BuildingPlan, bt building.BuildingType) bool {
	b := m.AddBuilding(x, y, bp, true)
	if b != nil {
		c.Towns[0].CreateBuildingConstruction(b, bt, m)
		return true
	} else {
		return false
	}
}

func (m *Map) AddRoadConstruction(c *social.Country, x, y uint16, rt *building.RoadType) bool {
	r := &building.Road{T: rt, Construction: true}
	m.GetField(x, y).Road = r
	c.Towns[0].CreateRoadConstruction(x, y, r, m)
	return true
}

func (m *Map) GetBuildingBaseFields(x, y uint16, bp *building.BuildingPlan) []navigation.FieldWithContext {
	var fields []navigation.FieldWithContext
	for i := uint16(0); i < 5; i++ {
		for j := uint16(0); j < 5; j++ {
			bx := int(x+i) - 2
			by := int(y+j) - 2
			if bp.BaseShape[i][j] != nil {
				if bx >= 0 && by >= 0 && bx < int(m.SX) && by < int(m.SY) {
					f := &m.Fields[bx][by]
					if !f.Buildable() {
						return nil
					} else {
						fields = append(fields, f)
					}
				} else {
					return nil
				}
			}
		}
	}
	return fields
}

func (m *Map) SetBuildingUnits(b *building.Building, construction bool) {
	bp := b.Plan
	for i := uint16(0); i < 5; i++ {
		for j := uint16(0); j < 5; j++ {
			bx := int(b.X+i) - 2
			by := int(b.Y+j) - 2
			if bp.BaseShape[i][j] != nil {
				m.Fields[bx][by].Building.BuildingUnits = b.ToBuildingUnits(uint8(i), uint8(j), construction)
				if !construction {
					m.Fields[bx][by].Building.RoofUnit = b.GetRoof(uint8(i), uint8(j))
				}
			}
		}
	}
}

func (m *Map) AddBuilding(x, y uint16, bp *building.BuildingPlan, construction bool) *building.Building {
	if m.GetBuildingBaseFields(x, y, bp) == nil {
		return nil
	}
	b := &building.Building{X: x, Y: y, Plan: bp}
	m.SetBuildingUnits(b, construction)
	return b
}

func (m *Map) ShortPath(sx, sy, ex, ey uint16, travellerType uint8) *navigation.Path {
	if sx == ex && sy == ey {
		return nil
	}
	p := FindShortPathBFS(m, sx, sy, m.GetField(ex, ey), travellerType)
	if p != nil {
		return &navigation.Path{F: p[1:]}
	}
	return nil
}

func (m *Map) FindDest(sx, sy uint16, dest navigation.Destination, travellerType uint8) *navigation.Field {
	p := FindShortPathBFS(m, sx, sy, dest, travellerType)
	if p != nil {
		return p[len(p)-1]
	}
	return nil
}
