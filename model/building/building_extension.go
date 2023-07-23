package building

import (
	"medvil/model/materials"
)

type BuildingExtensionType struct {
	Name    string
	OnWater bool
	InUnit  bool
}

var WaterMillWheel = &BuildingExtensionType{Name: "water_mill_wheel", OnWater: true, InUnit: false}
var Forge = &BuildingExtensionType{Name: "forge", OnWater: false, InUnit: false}
var Kiln = &BuildingExtensionType{Name: "kiln", OnWater: false, InUnit: false}
var Cooker = &BuildingExtensionType{Name: "cooker", OnWater: false, InUnit: true}
var Workshop = &BuildingExtensionType{Name: "workshop", OnWater: false, InUnit: true}
var Deck = &BuildingExtensionType{Name: "deck", OnWater: true, InUnit: false}
var NonExtension *BuildingExtensionType = nil

var BuildingExtensionTypes = [...]*BuildingExtensionType{
	WaterMillWheel,
	Forge,
	Kiln,
	Deck,
	Cooker,
	Workshop,
}

func GetBuildingExtensionType(name string) *BuildingExtensionType {
	for _, t := range BuildingExtensionTypes {
		if t.Name == name {
			return t
		}
	}
	return nil
}

type BuildingExtension struct {
	T *BuildingExtensionType
}

func GetExtensionDirection(t *BuildingExtensionType, x, y uint8, bp BuildingPlan) uint8 {
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
	case Deck:
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
	}
	return 255
}

func ForgeBuildingUnit(b *Building, m *materials.Material, construction bool) *BuildingUnit {
	w := &BuildingWall{M: m, Windows: WindowTypeNone, Door: false}
	return &BuildingUnit{
		BuildingComponentBase: BuildingComponentBase{B: b, Construction: construction},
		Walls:                 []*BuildingWall{w, w, w, w},
	}
}

func ForgeBuildingRoof(b *Building, m *materials.Material, construction bool) *RoofUnit {
	return &RoofUnit{
		BuildingComponentBase: BuildingComponentBase{B: b, Construction: construction},
		Roof:                  Roof{M: m, RoofType: RoofTypeSplit},
		Connected:             [4]bool{false, false, false, false},
	}
}
