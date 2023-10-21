package building

import (
	"math/rand"
)

const NumShapes = 15

func GetShape(t BuildingType, x, y uint16) uint8 {
	if t == BuildingTypeWall {
		return WallShape(x, y)
	} else if t == BuildingTypeWorkshop || t == BuildingTypeFactory {
		return uint8(rand.Intn(NumShapes))
	}
	return 0
}

func WallShape(x, y uint16) uint8 {
	return uint8(rand.Intn(NumShapes))
}
