package model

import (
	"math/rand"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/model/terrain"
	"medvil/model/time"
)

const ReedSpreadRate = 0.0001
const GrassGrowRate = 0.0001

type Map struct {
	SX        uint16
	SY        uint16
	Fields    [][]navigation.Field
	Countries []*social.Country
}

func (m *Map) SpreadPlant(i, j uint16) {
	if i >= 0 && j >= 0 && i < m.SX && j < m.SY {
		if m.Fields[i][j].Plant == nil && m.Fields[i][j].Road == nil && m.Shore(i, j) {
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
				if f.Plant.T.Name == "reed" && rand.Float64() < ReedSpreadRate {
					m.SpreadPlant(i-1, j)
					m.SpreadPlant(i, j-1)
					m.SpreadPlant(i+1, j)
					m.SpreadPlant(i, j+1)
				}
				if f.Plant.T.IsAnnual() && Calendar.Season() == time.Winter {
					f.Plant = nil
				}
			}
			if f.Animal != nil {
				f.Animal.ElapseTime(Calendar)
				if Calendar.Season() == time.Winter && !f.Animal.Corralled {
					f.Animal = nil
				}
			}
			if f.Terrain.T == terrain.Canal && m.HasNeighborField(f.X, f.Y, terrain.Water) {
				f.Terrain.T = terrain.Water
				f.Terrain.Resources.Add(artifacts.GetArtifact("water"), artifacts.InfiniteQuantity)
				navigation.SetRoadConnectionsForNeighbors(m, f)
				navigation.SetBuildingDeckForNeighbors(m, f)
			}
			if f.Plant == nil && f.Terrain.T == terrain.Dirt && rand.Float64() < GrassGrowRate && Calendar.Season() == time.Winter {
				f.Terrain.T = terrain.Grass
			}
		}
	}
}

func (m *Map) GetNField(x uint16, dx int, y uint16, dy int) *navigation.Field {
	return m.GetField(uint16(int(x)+dx), uint16(int(y)+dy))
}

func (m *Map) GetField(x uint16, y uint16) *navigation.Field {
	if x >= m.SX || y >= m.SY {
		return nil
	}
	return &m.Fields[x][y]
}

func (m *Map) ReverseReferences() *ReverseReferences {
	rr := BuildReverseReferences(m)
	return &rr
}

func (m *Map) AddBuildingConstruction(town *social.Town, x, y uint16, bp *building.BuildingPlan, direction uint8) bool {
	b := m.AddBuilding(x, y, bp.Copy(), true, direction)
	if b != nil {
		town.CreateBuildingConstruction(b, m)
		return true
	} else {
		return false
	}
}

func (m *Map) AddRoadConstruction(town *social.Town, x, y uint16, rt *building.RoadType) bool {
	r := &building.Road{T: rt, Construction: true, Broken: false}
	m.GetField(x, y).Road = r
	town.CreateRoadConstruction(x, y, r, m)
	return true
}

func (m *Map) AddInfraConstruction(town *social.Town, x, y uint16, it *building.InfraType) bool {
	town.CreateInfraConstruction(x, y, it, m)
	return true
}

func (m *Map) AddLevelingTask(town *social.Town, x, y uint16, taskType uint8) bool {
	town.CreateLevelingTask(m.GetField(x, y), taskType, m)
	return true
}

func (m *Map) AddWallRampConstruction(town *social.Town, x, y uint16) bool {
	f := m.GetField(x, y)
	b := f.Building.GetBuilding()
	rampD := navigation.GetRampDirection(m, x, y)
	if b != nil && b.Plan.BuildingType == building.BuildingTypeWall && rampD != building.DirectionNone {
		oldCost := b.Plan.ConstructionCost()
		roof := b.Plan.BaseShape[2][2].Roof
		roof.RoofType = building.RoofTypeRamp
		roof.RampD = rampD
		f.Building.BuildingComponents = b.ToBuildingUnits(2, 2, false)
		f.Building.BuildingComponents[len(f.Building.BuildingComponents)-1].SetConstruction(true)
		cost := artifacts.ArtifactsDiff(b.Plan.ConstructionCost(), oldCost)
		town.CreateIncrementalBuildingConstruction(b, cost, m)
		return true
	} else if b == nil && rampD != building.DirectionNone {
		bp := building.GetWallRampPlan(rampD)
		return m.AddBuildingConstruction(town, x, y, bp, building.DirectionNone)
	}
	return false
}

func (m *Map) CheckBuildingBaseField(pu *building.PlanUnits, bt building.BuildingType, f *navigation.Field, direction uint8) bool {
	if pu.Extension != nil && pu.Extension.T.OnWater {
		return f.Terrain.T == terrain.Water && f.Road == nil && f.Building.Empty()
	}
	if bt == building.BuildingTypeGate {
		if f.Terrain.T == terrain.Water {
			return (m.HasNeighborFieldInDirection(f.X, f.Y, terrain.Grass, (direction+1)%4) ||
				m.HasNeighborFieldInDirection(f.X, f.Y, terrain.Grass, (direction+3)%4))
		} else {
			return (!m.HasNonGateNeighborBuildingInDirection(f.X, f.Y, direction) &&
				!m.HasNonGateNeighborBuildingInDirection(f.X, f.Y, (direction+2)%4))
		}
	}
	if bt == building.BuildingTypeWall {
		return f.RoadCompatible()
	}
	return f.Buildable()
}

func (m *Map) GetBuildingBaseFields(x, y uint16, bp *building.BuildingPlan, direction uint8) []navigation.FieldWithContext {
	var fields []navigation.FieldWithContext
	for i := uint16(0); i < 5; i++ {
		for j := uint16(0); j < 5; j++ {
			bx := int(x+i) - 2
			by := int(y+j) - 2
			if bp.BaseShape[i][j] != nil {
				if bx >= 0 && by >= 0 && bx < int(m.SX) && by < int(m.SY) {
					f := &m.Fields[bx][by]
					if m.CheckBuildingBaseField(bp.BaseShape[i][j], bp.BuildingType, f, direction) {
						fields = append(fields, f)
					} else {
						return nil
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
				m.Fields[bx][by].Building.BuildingComponents = b.ToBuildingUnits(uint8(i), uint8(j), construction)
			}
		}
	}
}

func (m *Map) IsBuildingPossible(x, y uint16, bp *building.BuildingPlan, direction uint8) bool {
	return m.GetBuildingBaseFields(x, y, bp, direction) != nil
}

func (m *Map) AddBuilding(x, y uint16, bp *building.BuildingPlan, construction bool, direction uint8) *building.Building {
	if m.GetBuildingBaseFields(x, y, bp, direction) == nil {
		return nil
	}
	b := &building.Building{X: x, Y: y, Plan: *bp, Shape: building.GetShape(bp.BuildingType, x, y), Direction: direction}
	m.SetBuildingUnits(b, construction)
	return b
}

func (m *Map) ShortPath(start, dest navigation.Location, pathType navigation.PathType) *navigation.Path {
	if start == dest {
		return nil
	}
	p := FindShortPathBFS(m, start, dest, pathType)
	if p != nil {
		return &navigation.Path{P: p[1:]}
	}
	return nil
}

func (m *Map) FindDest(start navigation.Location, dest navigation.Destination, pathType navigation.PathType) *navigation.Field {
	p := FindShortPathBFS(m, start, dest, pathType)
	if p != nil {
		return p[len(p)-1].(*navigation.Field)
	}
	return nil
}

func (m *Map) HasNeighborFieldInDirection(x, y uint16, t *terrain.TerrainType, direction uint8) bool {
	d := navigation.DirectionOrthogonalXY[direction]
	if m.GetNField(x, d[0], y, d[1]) != nil && m.GetNField(x, d[0], y, d[1]).Terrain.T == t {
		return true
	}
	return false
}

func (m *Map) HasNonGateNeighborBuildingInDirection(x, y uint16, direction uint8) bool {
	d := navigation.DirectionOrthogonalXY[direction]
	f := m.GetNField(x, d[0], y, d[1])
	if f != nil && !f.Building.Empty() && f.Building.GetBuilding().Plan.BuildingType != building.BuildingTypeGate {
		return true
	}
	return false
}

func (m *Map) HasNeighborField(x, y uint16, t *terrain.TerrainType) bool {
	return (m.HasNeighborFieldInDirection(x, y, t, 0) ||
		m.HasNeighborFieldInDirection(x, y, t, 1) ||
		m.HasNeighborFieldInDirection(x, y, t, 2) ||
		m.HasNeighborFieldInDirection(x, y, t, 3))
}

func (m *Map) Shore(x, y uint16) bool {
	f := m.GetField(x, y)
	if f.Terrain.T != terrain.Water {
		return false
	}
	if m.HasNeighborField(x, y, terrain.Grass) {
		return true
	}
	return false
}
