package building

import (
	"medvil/model/artifacts"
)

var ConstructionInputs = []*artifacts.Artifact{
	artifacts.GetArtifact("board"),
	artifacts.GetArtifact("cube"),
}

type Construction struct {
	Building    *Building
	Road        *Road
	Progress    uint16
	MaxProgress uint16
	Cost        []artifacts.Artifacts
	Storage     *artifacts.Resources
	T           BuildingType
}

func (c *Construction) IsComplete() bool {
	return c.Progress == c.MaxProgress
}
