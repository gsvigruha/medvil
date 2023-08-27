package navigation

import (
	"math"
	"medvil/model/terrain"
)

const MaxTerraformFieldCornerDiff = 2

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
	if abs(newH-getElevation(f, (dir+1)%4)) > MaxTerraformFieldCornerDiff || abs(newH-getElevation(f, (dir+3)%4)) > MaxTerraformFieldCornerDiff {
		return false
	}
	{
		d1 := DirectionOrthogonalXY[dir]
		f1 := m.GetNField(f.X, d1[0], f.Y, d1[1])
		if f1 != nil && (!f1.Empty() || abs(newH-getElevation(*f1, (dir+1)%4)) > MaxTerraformFieldCornerDiff) {
			return false
		}
	}
	{
		d2 := DirectionOrthogonalXY[(dir+1)%4]
		f2 := m.GetNField(f.X, d2[0], f.Y, d2[1])
		if f2 != nil && (!f2.Empty() || abs(newH-getElevation(*f2, (dir+3)%4)) > MaxTerraformFieldCornerDiff) {
			return false
		}
	}
	{
		d3 := DirectionDiagonalXY[dir]
		f3 := m.GetNField(f.X, d3[0], f.Y, d3[1])
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
	if f.NE == f.NW && f.NE == f.SE && f.NE == f.SW {
		// No need to level, already suitable for buildings
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
		f1 := m.GetNField(f.X, d1[0], f.Y, d1[1])
		if f1 != nil {
			setElevation(f1, (dir+2)%4, newH)
		}
	}
	{
		d2 := DirectionOrthogonalXY[dir]
		f2 := m.GetNField(f.X, d2[0], f.Y, d2[1])
		if f2 != nil {
			setElevation(f2, (dir+1)%4, newH)
		}
	}
	{
		d3 := DirectionOrthogonalXY[(dir+1)%4]
		f3 := m.GetNField(f.X, d3[0], f.Y, d3[1])
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
	e1 := getElevation(f, dir)
	e2 := getElevation(f, (dir+1)%4)
	if e1 != e2 {
		// Cannot level, edge being checked is uneven
		return false
	}
	e3 := getElevation(f, (dir+2)%4)
	e4 := getElevation(f, (dir+3)%4)
	if e3 == e4 {
		// Unnecessary to level, the field is already good for road building
		return false
	}
	newH := int(math.Round((float64(e3)+float64(e4))/2.0 - 0.01))
	return checkCorner(f, (dir+2)%4, newH, m) && checkCorner(f, (dir+3)%4, newH, m)
}

func FieldCanBeLeveledForRoad(f Field, m IMap) bool {
	if !f.Empty() {
		return false
	}
	if !f.Terrain.T.Buildable {
		return false
	}
	for dir := uint8(0); dir < 4; dir++ {
		if checkEdge(f, dir, m) {
			return true
		}
	}
	return false
}

func LevelFieldForRoad(f *Field, m IMap) bool {
	if FieldCanBeLeveledForRoad(*f, m) {
		for dir := uint8(0); dir < 4; dir++ {
			if checkEdge(*f, dir, m) {
				e3 := getElevation(*f, (dir+2)%4)
				e4 := getElevation(*f, (dir+3)%4)
				newH := uint8(math.Round((float64(e3)+float64(e4))/2.0 - 0.01))
				setElevationForCorner(f, (dir+2)%4, uint8(newH), m)
				setElevationForCorner(f, (dir+3)%4, uint8(newH), m)
				return true
			}
		}
	}
	return false
}

func GetSurroundingType(f *Field, of1 *Field, of2 *Field, of3 *Field) uint8 {
	if f.Terrain.T == terrain.Grass && of1.Terrain.T == terrain.Water && of2.Terrain.T == terrain.Water && of3.Terrain.T == terrain.Water {
		return SurroundingWater
	} else if f.Terrain.T == terrain.Water && of1.Terrain.T == terrain.Grass && of2.Terrain.T == terrain.Grass && of3.Terrain.T == terrain.Grass {
		if of1.Flat() && of2.Flat() && of3.Flat() {
			return SurroundingGrass
		} else {
			return SurroundingDarkSlope
		}
	}
	if f.Terrain.T == terrain.Grass && of1.Terrain.T == terrain.Grass && of2.Terrain.T == terrain.Grass && of3.Terrain.T == terrain.Grass {
		if f.Flat() && !of1.Flat() && !of2.Flat() && !of3.Flat() {
			return SurroundingDarkSlope
		}
	}
	return SurroundingSame
}

func SetSurroundingTypes(m IMap, f *Field) {
	i := int(f.X)
	j := int(f.Y)
	sx, sy := m.Size()
	for k, _ := range DirectionOrthogonalXY {
		di1 := DirectionOrthogonalXY[k][0]
		dj1 := DirectionOrthogonalXY[k][1]
		di2 := DirectionOrthogonalXY[(k+1)%4][0]
		dj2 := DirectionOrthogonalXY[(k+1)%4][1]
		di3 := DirectionDiagonalXY[k][0]
		dj3 := DirectionDiagonalXY[k][1]
		if i > 0 && j > 0 && i < int(sx)-1 && j < int(sy)-1 {
			f.Surroundings[k] = GetSurroundingType(f, m.GetField(uint16(i+di1), uint16(j+dj1)), m.GetField(uint16(i+di2), uint16(j+dj2)), m.GetField(uint16(i+di3), uint16(j+dj3)))
		}
	}
}

func SetSurroundingTypesForNeighbors(m IMap, f *Field) {
	SetSurroundingTypes(m, f)
	for i := 0; i < 8; i++ {
		d := DirectionAllXY[i]
		of := m.GetNField(f.X, d[0], f.Y, d[1])
		if of != nil {
			SetSurroundingTypes(m, of)
		}
	}
}
