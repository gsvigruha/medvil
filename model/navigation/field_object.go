package navigation

import (
	"medvil/model/building"
)

type FieldObject interface {
	Walkable() bool
	LiftN() int8
	LiftE() int8
	LiftS() int8
	LiftW() int8
}

type FieldBuildingObjects struct {
	BuildingUnits []building.BuildingUnit
	RoofUnit      *building.RoofUnit
}

func (o FieldBuildingObjects) Empty() bool {
	return o.RoofUnit == nil && len(o.BuildingUnits) == 0
}
