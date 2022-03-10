package building

import (
	"medvil/model/materials"
)

type BuildingType uint8

const BuildingTypeFarm = 1
const BuildingTypeWorkshop = 2
const BuildingTypeMine = 3

const BuildingTypeRoad = 10

func FloorMaterials(bt BuildingType) []*materials.Material {
	switch bt {
	case BuildingTypeFarm:
		return []*materials.Material{
			materials.GetMaterial("whitewash"),
		}
	case BuildingTypeWorkshop:
		return []*materials.Material{
			materials.GetMaterial("stone"),
			materials.GetMaterial("sandstone"),
			materials.GetMaterial("brick"),
		}
	case BuildingTypeMine:
		return []*materials.Material{
			materials.GetMaterial("stone"),
			materials.GetMaterial("wood"),
		}
	}
	return nil
}

func RoofMaterials(bt BuildingType) []*materials.Material {
	switch bt {
	case BuildingTypeFarm:
		return []*materials.Material{
			materials.GetMaterial("hay"),
			materials.GetMaterial("tile"),
		}
	case BuildingTypeWorkshop:
		return []*materials.Material{
			materials.GetMaterial("tile"),
		}
	case BuildingTypeMine:
		return []*materials.Material{
			materials.GetMaterial("hay"),
			materials.GetMaterial("tile"),
		}
	}
	return nil
}
