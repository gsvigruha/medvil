package building

import (
	"medvil/model/artifacts"
)

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
