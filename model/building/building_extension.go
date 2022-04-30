package building

type BuildingExtensionType uint8

const WaterMillWheel BuildingExtensionType = 1
const Forge BuildingExtensionType = 2

type BuildingExtension struct {
	T BuildingExtensionType
}

func GetExtensionDirection(t BuildingExtensionType, x, y uint8, bp BuildingPlan) uint8 {
	switch t {
	case WaterMillWheel:
		// TODO: migrate direction and use it here
		if bp.HasUnit(x, y-1, 0) {
			return 0
		} else if bp.HasUnit(x+1, y, 0) {
			return 1
		} else if bp.HasUnit(x, y+1, 0) {
			return 2
		} else if bp.HasUnit(x-1, y, 0) {
			return 3
		}
	case Forge:
		// TODO: migrate direction and use it here
		if bp.HasUnit(x, y-1, 0) {
			return 3
		} else if bp.HasUnit(x+1, y, 0) {
			return 0
		} else if bp.HasUnit(x, y+1, 0) {
			return 1
		} else if bp.HasUnit(x-1, y, 0) {
			return 2
		}
	}
	return 255
}
