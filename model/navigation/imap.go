package navigation

import (
	"medvil/model/building"
)

type IMap interface {
	GetField(x uint16, y uint16) *Field
	ShortPath(start, dest Location, travellerType uint8) *Path
	FindDest(start Location, dest Destination, travellerType uint8) *Field
	SetBuildingUnits(b *building.Building, construction bool)
	Shore(x, y uint16) bool
	GetBuildingBaseFields(x, y uint16, bp *building.BuildingPlan) []FieldWithContext
}
