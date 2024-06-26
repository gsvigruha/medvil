package social

import (
	"medvil/model/building"
	"medvil/model/navigation"
)

type Supplier interface {
	GetHome() Home
	ReassignFirstPerson(dstH Home, assingTask bool, m navigation.IMap)
	FieldWithinDistance(*navigation.Field) bool
	CreateBuildingConstruction(b *building.Building, m navigation.IMap)
	AddConstruction(c *building.Construction)
	BuildMarketplaceEnabled() bool
	BuildHousesEnabled() bool
	CreateLevelingTask(f *navigation.Field, taskType uint8, m navigation.IMap)
	GetName() string
	CloseToTown(town *Town, m navigation.IMap) bool
	HasHousehold(bt building.BuildingType) bool
}
