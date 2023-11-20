package building

import (
	"medvil/model/materials"
)

type BuildingType uint8

const BuildingTypeFarm = 1
const BuildingTypeWorkshop = 2
const BuildingTypeMine = 3
const BuildingTypeFactory = 4

const BuildingTypeRoad = 10
const BuildingTypeCanal = 11
const BuildingTypeAqueduct = 12
const BuildingTypeBridge = 13
const BuildingTypeWall = 14
const BuildingTypeGate = 15
const BuildingTypeTower = 16

const BuildingTypeTownhall = 20
const BuildingTypeMarket = 21

const BuildingTypeStatue = 31

func FloorMaterials(bt BuildingType) []*materials.Material {
	switch bt {
	case BuildingTypeFarm:
		return []*materials.Material{
			materials.GetMaterial("whitewash"),
		}
	case BuildingTypeWorkshop:
		return []*materials.Material{
			materials.GetMaterial("sandstone"),
			materials.GetMaterial("brick"),
		}
	case BuildingTypeMine:
		return []*materials.Material{
			materials.GetMaterial("stone"),
			materials.GetMaterial("wood"),
		}
	case BuildingTypeFactory:
		return []*materials.Material{
			materials.GetMaterial("sandstone"),
			materials.GetMaterial("brick"),
		}
	case BuildingTypeTower:
		return []*materials.Material{
			materials.GetMaterial("stone"),
		}
	case BuildingTypeTownhall:
		return []*materials.Material{
			materials.GetMaterial("marble"),
		}
	case BuildingTypeMarket:
		return []*materials.Material{
			materials.GetMaterial("wood"),
		}
	}
	return nil
}

func RoofMaterials(bt BuildingType) []*materials.Material {
	switch bt {
	case BuildingTypeFarm:
		return []*materials.Material{
			materials.GetMaterial("reed"),
			materials.GetMaterial("tile"),
		}
	case BuildingTypeWorkshop:
		return []*materials.Material{
			materials.GetMaterial("tile"),
		}
	case BuildingTypeMine:
		return []*materials.Material{
			materials.GetMaterial("reed"),
			materials.GetMaterial("tile"),
		}
	case BuildingTypeFactory:
		return []*materials.Material{}
	case BuildingTypeTower:
		return []*materials.Material{
			materials.GetMaterial("tile"),
		}
	case BuildingTypeTownhall:
		return []*materials.Material{
			materials.GetMaterial("copper"),
		}
	case BuildingTypeMarket:
		return []*materials.Material{
			materials.GetMaterial("textile"),
		}
	}
	return nil
}

func ExtensionTypes(bt BuildingType) []*BuildingExtensionType {
	switch bt {
	case BuildingTypeWorkshop:
		return []*BuildingExtensionType{
			WaterMillWheel,
			Forge,
			Kiln,
			Cooker,
			Workshop,
		}
	}
	return nil
}

func NeedsRoof(bt BuildingType) bool {
	switch bt {
	case BuildingTypeTower, BuildingTypeMarket, BuildingTypeMine, BuildingTypeFarm:
		return true
	default:
		return false
	}
}

func MinNumFloors(bt BuildingType) int {
	switch bt {
	case BuildingTypeFactory, BuildingTypeTower:
		return 2
	default:
		return 1
	}
}

func MaxNumFloors(bt BuildingType) int {
	switch bt {
	case BuildingTypeFarm, BuildingTypeMine, BuildingTypeFactory:
		return 2
	case BuildingTypeWorkshop, BuildingTypeTownhall, BuildingTypeTower:
		return 3
	case BuildingTypeMarket:
		return 1
	default:
		return 2
	}
}
