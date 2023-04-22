package building

import (
	"math/rand"
)

const NumShapes = 10

const ShapeAreaSize = 3

func GetShape(t BuildingType, x, y uint16) uint8 {
	if t == BuildingTypeWall {
		return WallShape(x, y)
	} else if t == BuildingTypeWorkshop {
		return uint8(rand.Intn(NumShapes))
	}
	return 0
}

func WallShape(x, y uint16) uint8 {
	return uint8((13*(x/ShapeAreaSize) + 7*(y/ShapeAreaSize)) % NumShapes)
}
