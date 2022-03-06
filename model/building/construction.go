package building

import (
	"medvil/model/artifacts"
)

type Construction interface {
	GetCost() []artifacts.Artifacts
	GetStorage() *artifacts.Resources
	SetMaxProgress(uint16)
	IsComplete() bool
	GetT() BuildingType
	GetBuilding() *Building
	GetRoad() *Road
	IncProgress()
	GetProgress() uint16
	GetMaxProgress() uint16
}

type BuildingConstruction struct {
	Building    *Building
	Progress    uint16
	MaxProgress uint16
	Cost        []artifacts.Artifacts
	Storage     *artifacts.Resources
	T           BuildingType
}

func (c *BuildingConstruction) IsComplete() bool {
	return c.Progress == c.MaxProgress
}

func (c *BuildingConstruction) SetMaxProgress(p uint16) {
	c.MaxProgress = p
}

func (c *BuildingConstruction) GetCost() []artifacts.Artifacts {
	return c.Cost
}

func (c *BuildingConstruction) GetStorage() *artifacts.Resources {
	return c.Storage
}

func (c *BuildingConstruction) GetT() BuildingType {
	return c.T
}

func (c *BuildingConstruction) GetBuilding() *Building {
	return c.Building
}

func (c *BuildingConstruction) GetRoad() *Road {
	return nil
}

func (c *BuildingConstruction) IncProgress() {
	c.Progress++
}

func (c *BuildingConstruction) GetProgress() uint16 {
	return c.Progress
}

func (c *BuildingConstruction) GetMaxProgress() uint16 {
	return c.MaxProgress
}

type RoadConstruction struct {
	Road        *Road
	Progress    uint16
	MaxProgress uint16
	Cost        []artifacts.Artifacts
	Storage     *artifacts.Resources
}

func (c *RoadConstruction) IsComplete() bool {
	return c.Progress == c.MaxProgress
}

func (c *RoadConstruction) SetMaxProgress(p uint16) {
	c.MaxProgress = p
}

func (c *RoadConstruction) GetCost() []artifacts.Artifacts {
	return c.Cost
}

func (c *RoadConstruction) GetStorage() *artifacts.Resources {
	return c.Storage
}

func (c *RoadConstruction) GetT() BuildingType {
	return BuildingTypeRoad
}

func (c *RoadConstruction) GetBuilding() *Building {
	return nil
}

func (c *RoadConstruction) GetRoad() *Road {
	return c.Road
}

func (c *RoadConstruction) IncProgress() {
	c.Progress++
}

func (c *RoadConstruction) GetProgress() uint16 {
	return c.Progress
}

func (c *RoadConstruction) GetMaxProgress() uint16 {
	return c.MaxProgress
}
