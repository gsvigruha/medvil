package building

import (
	"math/rand"
)

const WorkshopShapes = 15
const WallShapes = 8
const FarmShapes = 3

func GetShape(t BuildingType, x, y uint16) uint8 {
	if t == BuildingTypeWall {
		return uint8(rand.Intn(WallShapes))
	} else if t == BuildingTypeWorkshop || t == BuildingTypeFactory {
		return uint8(rand.Intn(WorkshopShapes))
	} else if t == BuildingTypeFarm {
		return uint8(rand.Intn(FarmShapes))
	}
	return 0
}
