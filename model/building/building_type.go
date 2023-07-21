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
			materials.GetMaterial("stone"),
			materials.GetMaterial("wood"),
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
			Cooker,
		}
	}
	return nil
}

func MaxNumFloors(bt BuildingType) int {
	switch bt {
	case BuildingTypeFarm, BuildingTypeMine:
		return 2
	case BuildingTypeWorkshop, BuildingTypeFactory, BuildingTypeTownhall:
		return 3
	case BuildingTypeMarket:
		return 1
	default:
		return 2
	}
}
