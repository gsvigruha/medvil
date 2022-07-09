package navigation

import (
	"medvil/model/building"
	"medvil/model/terrain"
)

func RampPossible(m IMap, x, y, x2, y2 uint16) bool {
	f := m.GetField(x, y)
	f2 := m.GetField(x2, y2)
	var d = 1
	if f.Building.Empty() {
		d = 2
	}
	if f2 != nil && !f2.Building.Empty() &&
		f2.Building.GetBuilding().Plan.BuildingType == building.BuildingTypeWall &&
		len(f.Building.BuildingComponents)+d == len(f2.Building.BuildingComponents) {
		return true
	}
	return false
}

func GetRampDirection(m IMap, x, y uint16) uint8 {
	if RampPossible(m, x, y, x, y-1) {
		return building.DirectionN
	} else if RampPossible(m, x, y, x+1, y) {
		return building.DirectionE
	} else if RampPossible(m, x, y, x, y+1) {
		return building.DirectionS
	} else if RampPossible(m, x, y, x-1, y) {
		return building.DirectionW
	}
	return DirectionNone
}

func SetRoadConnectionsForNeighbors(m IMap, f *Field) {
	for i := 0; i < 8; i++ {
		d := DirectionAllXY[i]
		of := m.GetField(uint16(int(f.X)+d[0]), uint16(int(f.Y)+d[1]))
		if of != nil && of.Road != nil {
			SetRoadConnections(m, of)
		}
	}
}

func SetRoadConnections(m IMap, f *Field) {
	for i := 0; i < 4; i++ {
		d := DirectionOrthogonalXY[i]
		of := m.GetField(uint16(int(f.X)+d[0]), uint16(int(f.Y)+d[1]))
		if of != nil {
			if of.Road != nil {
				f.Road.EdgeConnections[i] = true
				of.Road.EdgeConnections[(i+2)%4] = true
			}
			if f.Road.T.Bridge {
				f.Road.EdgeConnections[i] = (of.Terrain.T != terrain.Water || of.Road != nil)
			}
			if !of.Building.Empty() && f.Terrain.T != terrain.Water {
				f.Road.EdgeConnections[i] = true
				b := of.Building.GetBuilding()
				if b.Plan.BuildingType != building.BuildingTypeWall && of.X == b.X && of.Y == b.Y {
					if unit, ok := of.Building.BuildingComponents[0].(*building.BuildingUnit); ok {
						if !unit.HasDoor() {
							unit.Walls[(i+2)%4].Door = true
						}
					}
				}
			}
		}
	}
	for i := 0; i < 4; i++ {
		d := DirectionDiagonalXY[i]
		of := m.GetField(uint16(int(f.X)+d[0]), uint16(int(f.Y)+d[1]))
		if of != nil && of.Road != nil {
			f.Road.CornerConnections[i] = true
			of.Road.CornerConnections[(i+2)%4] = true
		}
		if !of.Building.Empty() {
			f.Road.CornerConnections[i] = true
		}
	}
}

func SetBuildingDeck(m IMap, f *Field, of *Field) {
	b := of.Building.GetBuilding()
	if b != nil && !of.Building.IsBuildingExtension() && f.Building.GetBuilding() == nil {
		i := f.X + 2 - b.X
		j := f.Y + 2 - b.Y
		if i < 5 && j < 5 {
			b.Plan.BaseShape[i][j] = &building.PlanUnits{Extension: &building.BuildingExtension{T: building.Deck}}
			f.Building.BuildingComponents = b.ToBuildingUnits(uint8(i), uint8(j), false)
		}
	}
}

func SetBuildingDeckForNeighbors(m IMap, f *Field) {
	for i := 0; i < 4; i++ {
		d := DirectionOrthogonalXY[i]
		of := m.GetField(uint16(int(f.X)+d[0]), uint16(int(f.Y)+d[1]))
		if f != nil && f.Terrain.T == terrain.Water && of != nil && of.Building.GetBuilding() != nil {
			SetBuildingDeck(m, f, of)
		}
		if of != nil && of.Terrain.T == terrain.Water && f != nil && f.Building.GetBuilding() != nil {
			SetBuildingDeck(m, of, f)
		}
	}
}
