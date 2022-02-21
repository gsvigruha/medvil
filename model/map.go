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
			for k := range town.Workshops {
				workshop := town.Workshops[k]
				for l := range workshop.Household.People {
					person := workshop.Household.People[l]
					person.ElapseTime(Calendar, m)
				}
				workshop.ElapseTime(Calendar, m)
			}
		}
	}
	for i := uint16(0); i < m.SX; i++ {
		for j := uint16(0); j < m.SY; j++ {
			f := &m.Fields[i][j]
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
		f.Household.Resources.VolumeCapacity = f.Household.Building.Plan.Area() * social.StoragePerArea
		t.Farms = append(t.Farms, f)
		return true
	} else {
		return false
	}
}

func (m *Map) AddWorkshop(c *social.Country, x, y uint16, bp *building.BuildingPlan) bool {
	b := m.AddBuilding(x, y, bp)
	if b != nil {
		t := c.Towns[0]
		w := &social.Workshop{Household: social.Household{Building: b, Town: t}}
		w.Household.Resources.VolumeCapacity = w.Household.Building.Plan.Area() * social.StoragePerArea
		t.Workshops = append(t.Workshops, w)
		return true
	} else {
		return false
	}
}

func (m *Map) GetBuildingBaseFields(x, y uint16, bp *building.BuildingPlan) []navigation.FieldWithContext {
	var fields []navigation.FieldWithContext
	for i := uint16(0); i < 5; i++ {
		for j := uint16(0); j < 5; j++ {
			bx := int(x+i) - 2
			by := int(y+j) - 2
			if bp.BaseShape[i][j] {
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

func (m *Map) AddBuilding(x, y uint16, bp *building.BuildingPlan) *building.Building {
	b := building.Building{X: x, Y: y, Plan: bp}
	if m.GetBuildingBaseFields(x, y, bp) == nil {
		return nil
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
