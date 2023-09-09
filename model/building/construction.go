package building

import (
	"medvil/model/artifacts"
)

var ConstructionInputs = []*artifacts.Artifact{
	artifacts.GetArtifact("board"),
	artifacts.GetArtifact("cube"),
	artifacts.GetArtifact("thatch"),
	artifacts.GetArtifact("brick"),
	artifacts.GetArtifact("tile"),
}

type Construction struct {
	Building    *Building
	Road        *Road
	Statue      *Statue
	X           uint16
	Y           uint16
	Progress    uint16
	MaxProgress uint16
	Cost        []artifacts.Artifacts
	Storage     *artifacts.Resources
	T           BuildingType
	Expired     bool
}

func (c *Construction) IsComplete() bool {
	return c.Progress >= c.MaxProgress
}
