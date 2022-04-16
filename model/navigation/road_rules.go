package navigation

import (
	"medvil/model/building"
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

func SetRoadConnections(m IMap, f *Field) {
	for i := 0; i < 4; i++ {
		d := DirectionXY[i]
		of := m.GetField(uint16(int(f.X)+d[0]), uint16(int(f.Y)+d[1]))
		if of != nil && of.Road != nil {
			f.Road.EdgeConnections[i] = true
			of.Road.EdgeConnections[(i+2)%4] = true
		}
	}
	for i := 0; i < 4; i++ {
		d := DirectionDiagonalXY[i]
		of := m.GetField(uint16(int(f.X)+d[0]), uint16(int(f.Y)+d[1]))
		if of != nil && of.Road != nil {
			f.Road.CornerConnections[i] = true
			of.Road.CornerConnections[(i+2)%4] = true
		}
	}
}
