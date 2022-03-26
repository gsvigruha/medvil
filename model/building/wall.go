package building

import (
	"medvil/model/materials"
)

var WallFloor = Floor{M: materials.GetMaterial("stone")}
var WallRoof = &Roof{M: materials.GetMaterial("stone"), RoofType: RoofTypeFlat}

var Wall1 = &PlanUnits{
	Floors: []Floor{WallFloor},
	Roof:   WallRoof,
}

var Wall2 = &PlanUnits{
	Floors: []Floor{WallFloor, WallFloor},
	Roof:   WallRoof,
}

var Wall3 = &PlanUnits{
	Floors: []Floor{WallFloor, WallFloor, WallFloor},
	Roof:   WallRoof,
}

var StoneWall1Type = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Wall1, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeWall,
}

var StoneWall2Type = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Wall2, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeWall,
}

var StoneWall3Type = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Wall3, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeWall,
}

var StoneWallRampType = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Wall3, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeWall,
}

func GetWallRampPlan(rampD uint8) *BuildingPlan {
	ramp := &PlanUnits{Roof: &Roof{M: materials.GetMaterial("stone"), RoofType: RoofTypeRamp, RampD: rampD}}
	return &BuildingPlan{
		BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
			{nil, nil, nil, nil, nil},
			{nil, nil, nil, nil, nil},
			{nil, nil, ramp, nil, nil},
			{nil, nil, nil, nil, nil},
			{nil, nil, nil, nil, nil},
		},
		BuildingType: BuildingTypeWall,
	}
}
