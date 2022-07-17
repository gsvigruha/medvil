package building

import (
	"medvil/model/materials"
)

var WallFloor = Floor{M: materials.GetMaterial("stone")}
var WallRoof = &Roof{M: materials.GetMaterial("stone"), RoofType: RoofTypeFlat}
var TowerRoof = &Roof{M: materials.GetMaterial("tile"), RoofType: RoofTypeSplit}

var Wall1 = &PlanUnits{
	Floors: []Floor{WallFloor},
	Roof:   WallRoof,
}

var Wall2 = &PlanUnits{
	Floors: []Floor{WallFloor, WallFloor},
	Roof:   WallRoof,
}

var Tower1 = &PlanUnits{
	Floors: []Floor{WallFloor, WallFloor},
	Roof:   TowerRoof,
}

var Tower2 = &PlanUnits{
	Floors: []Floor{WallFloor, WallFloor, WallFloor},
	Roof:   TowerRoof,
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

var Tower1Type = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Tower1, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeTower,
}

var Tower2Type = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Tower2, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeTower,
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

var Gate1 = &PlanUnits{
	Floors: []Floor{WallFloor},
	Roof:   WallRoof,
}
var SmallGate = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Gate1, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeGate,
}
