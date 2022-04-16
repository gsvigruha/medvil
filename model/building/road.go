package building

import (
	"medvil/model/artifacts"
)

type RoadType struct {
	Name   string
	Speed  float64
	Bridge bool
	Cost   []artifacts.Artifacts
}

var DirtRoadType = &RoadType{
	Name: "dirt_road", Speed: 1.5, Bridge: false, Cost: []artifacts.Artifacts{},
}
var CobbleRoadType = &RoadType{
	Name: "cobble_road", Speed: 2.0, Bridge: false,
	Cost: []artifacts.Artifacts{artifacts.Artifacts{artifacts.GetArtifact("cube"), 1}},
}
var BridgeRoadType = &RoadType{
	Name: "bridge", Speed: 1.5, Bridge: true,
	Cost: []artifacts.Artifacts{artifacts.Artifacts{artifacts.GetArtifact("board"), 3}},
}

type Road struct {
	T                 *RoadType
	Construction      bool
	EdgeConnections   [4]bool
	CornerConnections [4]bool
}
