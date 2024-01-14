package building

import (
	"medvil/model/artifacts"
)

type RoadType struct {
	Name        string
	DisplayName string
	Speed       float64
	Bridge      bool
	Cost        []artifacts.Artifacts
}

var DirtRoadType = &RoadType{
	Name: "dirt_road", DisplayName: "Dirt road", Speed: 1.5, Bridge: false, Cost: []artifacts.Artifacts{},
}
var CobbleRoadType = &RoadType{
	Name: "cobble_road", DisplayName: "Cobble stone road", Speed: 2.0, Bridge: false,
	Cost: []artifacts.Artifacts{artifacts.Artifacts{artifacts.GetArtifact("cube"), 1}},
}
var BridgeRoadType = &RoadType{
	Name: "bridge", DisplayName: "Bridge", Speed: 1.5, Bridge: true,
	Cost: []artifacts.Artifacts{artifacts.Artifacts{artifacts.GetArtifact("board"), 3}, artifacts.Artifacts{artifacts.GetArtifact("paper"), 2}},
}

var RoadTypes = [...]*RoadType{
	DirtRoadType,
	CobbleRoadType,
	BridgeRoadType,
}

func GetRoadType(name string) *RoadType {
	for _, t := range RoadTypes {
		if t.Name == name {
			return t
		}
	}
	return nil
}

type Road struct {
	T                 *RoadType
	Construction      bool
	Broken            bool
	EdgeConnections   [4]bool
	CornerConnections [4]bool
}

func (r *Road) Repair() {
	r.Broken = false
}

func (r *Road) RepairCost() []artifacts.Artifacts {
	return r.T.Cost
}

func (r *Road) Name() string {
	return "road"
}
