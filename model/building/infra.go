package building

import (
	"medvil/model/artifacts"
)

type InfraType struct {
	Name string
	Cost []artifacts.Artifacts
	BT   BuildingType
}

var CanalType = &InfraType{Name: "canal", Cost: []artifacts.Artifacts{artifacts.Artifacts{artifacts.GetArtifact("cube"), 1}}, BT: BuildingTypeCanal}
