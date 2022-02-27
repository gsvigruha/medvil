package building

type RoadType struct {
	Name  string
	Speed float64
}

var DirtRoadType = &RoadType{Name: "dirt_road", Speed: 1.5}
var CobbleRoadType = &RoadType{Name: "cobble_road", Speed: 2.0}

type Road struct {
	T *RoadType
}
