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
	BuildingComponents []building.BuildingComponent
}

func (o FieldBuildingObjects) Empty() bool {
	return len(o.BuildingComponents) == 0
}

func (o FieldBuildingObjects) GetBuilding() *building.Building {
	if o.Empty() {
		return nil
	}
	return o.BuildingComponents[0].Building()
}
