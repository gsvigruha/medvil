package navigation

import (
	"math"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func getElevation(f Field, dir uint8) int {
	if dir == 0 {
		return int(f.NE)
	} else if dir == 1 {
		return int(f.SE)
	} else if dir == 2 {
		return int(f.SW)
	} else /* if dir == 3 */ {
		return int(f.NW)
	}
}

func setElevation(f *Field, dir uint8, newH uint8) {
	if dir == 0 {
		f.NE = newH
	} else if dir == 1 {
		f.SE = newH
	} else if dir == 2 {
		f.SW = newH
	} else /* if dir == 3 */ {
		f.NW = newH
	}
}

func checkCorner(f Field, dir uint8, newH int, m IMap) bool {
	c := getElevation(f, dir)
	if c == newH {
		return true
	}
	if abs(newH-getElevation(f, (dir+1)%4)) > MaxFieldCornerDiff || abs(newH-getElevation(f, (dir+3)%4)) > MaxFieldCornerDiff {
		return false
	}
	{
		d1 := DirectionOrthogonalXY[dir]
		f1 := m.GetField(uint16(int(f.X)+d1[0]), uint16(int(f.Y)+d1[1]))
		if f1 != nil && (!f1.Empty() || abs(newH-getElevation(*f1, (dir+1)%4)) > MaxFieldCornerDiff) {
			return false
		}
	}
	{
		d2 := DirectionOrthogonalXY[(dir+1)%4]
		f2 := m.GetField(uint16(int(f.X)+d2[0]), uint16(int(f.Y)+d2[1]))
		if f2 != nil && (!f2.Empty() || abs(newH-getElevation(*f2, (dir+3)%4)) > MaxFieldCornerDiff) {
			return false
		}
	}
	{
		d3 := DirectionDiagonalXY[dir]
		f3 := m.GetField(uint16(int(f.X)+d3[0]), uint16(int(f.Y)+d3[1]))
		if f3 != nil && !f3.Empty() {
			return false
		}
	}
	return true
}

func averageHeight(f Field) int {
	return int(math.Round((float64(f.NE)+float64(f.NW)+float64(f.SE)+float64(f.SW))/4.0 - 0.01))
}

func FieldCanBeLeveledForBuilding(f Field, m IMap) bool {
	if !f.Empty() {
		return false
	}
	if !f.Terrain.T.Buildable {
		return false
	}
	avgH := averageHeight(f)
	return (checkCorner(f, 0, avgH, m) &&
		checkCorner(f, 1, avgH, m) &&
		checkCorner(f, 2, avgH, m) &&
		checkCorner(f, 3, avgH, m))
}

func setElevationForCorner(f *Field, dir uint8, newH uint8, m IMap) {
	setElevation(f, dir, newH)
	{
		d1 := DirectionDiagonalXY[dir]
		f1 := m.GetField(uint16(int(f.X)+d1[0]), uint16(int(f.Y)+d1[1]))
		if f1 != nil {
			setElevation(f1, (dir+2)%4, newH)
		}
	}
	{
		d2 := DirectionOrthogonalXY[dir]
		f2 := m.GetField(uint16(int(f.X)+d2[0]), uint16(int(f.Y)+d2[1]))
		if f2 != nil {
			setElevation(f2, (dir+1)%4, newH)
		}
	}
	{
		d3 := DirectionOrthogonalXY[(dir+1)%4]
		f3 := m.GetField(uint16(int(f.X)+d3[0]), uint16(int(f.Y)+d3[1]))
		if f3 != nil {
			setElevation(f3, (dir+3)%4, newH)
		}
	}
}

func LevelFieldForBuilding(f *Field, m IMap) bool {
	if FieldCanBeLeveledForBuilding(*f, m) {
		avgH := uint8(averageHeight(*f))
		setElevationForCorner(f, 0, avgH, m)
		setElevationForCorner(f, 1, avgH, m)
		setElevationForCorner(f, 2, avgH, m)
		setElevationForCorner(f, 3, avgH, m)
		return true
	}
	return false
}

func checkEdge(f Field, dir uint8, m IMap) bool {
	return true
}

func FieldCanBeLeveledForRoad(f Field, m IMap) bool {
	if !f.Empty() {
		return false
	}
	if !f.Terrain.T.Buildable {
		return false
	}
	return (checkEdge(f, 0, m) &&
		checkEdge(f, 1, m) &&
		checkEdge(f, 2, m) &&
		checkEdge(f, 3, m))
}

func LevelFieldForRoad(f *Field, m IMap) bool {
	return true
}
