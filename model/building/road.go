package building

import (
	"medvil/model/artifacts"
)

type RoadType struct {
	Name  string
	Speed float64
	Cost  []artifacts.Artifacts
}

var DirtRoadType = &RoadType{Name: "dirt_road", Speed: 1.5, Cost: []artifacts.Artifacts{}}
var CobbleRoadType = &RoadType{Name: "cobble_road", Speed: 2.0, Cost: []artifacts.Artifacts{artifacts.Artifacts{artifacts.GetArtifact("cube"), 1}}}

type Road struct {
	T            *RoadType
	Construction bool
}
