package building

import (
	"medvil/model/materials"
)

type BuildingExtensionType struct {
	Name        string
	OnWater     bool
	InUnit      bool
	Description string
}

var WaterMillWheel = &BuildingExtensionType{Name: "water_mill_wheel", OnWater: true, InUnit: false, Description: "Waterwheels are needed for milling."}
var Forge = &BuildingExtensionType{Name: "forge", OnWater: false, InUnit: false, Description: "Forges are used to work metals."}
var Kiln = &BuildingExtensionType{Name: "kiln", OnWater: false, InUnit: false, Description: "Kilns are needed to burn clay."}
var Cooker = &BuildingExtensionType{Name: "cooker", OnWater: false, InUnit: true, Description: "Cookers are used to make food and chemicals."}
var Workshop = &BuildingExtensionType{Name: "workshop", OnWater: false, InUnit: true, Description: "Workshops are used to work raw materials."}
var Deck = &BuildingExtensionType{Name: "deck", OnWater: true, InUnit: false}
var NonExtension *BuildingExtensionType = nil

var BuildingExtensionTypes = [...]*BuildingExtensionType{
	Workshop,
	Cooker,
	WaterMillWheel,
	Kiln,
	Forge,
	Deck,
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

func ForgeBuildingRoof(b *Building, rm *materials.Material, wm *materials.Material, construction bool) *RoofUnit {
	return &RoofUnit{
		BuildingComponentBase: BuildingComponentBase{B: b, Construction: construction},
		Roof:                  Roof{M: rm, RoofType: RoofTypeSplit},
		WallM:                 wm,
		Connected:             [4]bool{false, false, false, false},
	}
}
