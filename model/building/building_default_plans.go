package building

import (
	"medvil/model/materials"
)

var WhitewashFloor = Floor{M: materials.GetMaterial("whitewash")}
var StoneFloor = Floor{M: materials.GetMaterial("stone")}
var WoodFloor = Floor{M: materials.GetMaterial("wood")}
var BrickFloor = Floor{M: materials.GetMaterial("brick")}
var SandstoneFloor = Floor{M: materials.GetMaterial("sandstone")}

var ReedRoof = &Roof{M: materials.GetMaterial("reed"), RoofType: RoofTypeSplit}
var TileRoof = &Roof{M: materials.GetMaterial("tile"), RoofType: RoofTypeSplit}
var BrickRoof = &Roof{M: materials.GetMaterial("brick"), RoofType: RoofTypeFlat}
var SandstoneRoof = &Roof{M: materials.GetMaterial("sandstone"), RoofType: RoofTypeFlat}

// Farms

var Farm1 = &PlanUnits{
	Floors: []Floor{WhitewashFloor},
	Roof:   ReedRoof,
}

var Farm2 = &PlanUnits{
	Floors: []Floor{WhitewashFloor, WhitewashFloor},
	Roof:   ReedRoof,
}

var Farm1Plan = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Farm2, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeFarm,
}

var Farm2Plan = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, Farm1, Farm2, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeFarm,
}

var Farm3Plan = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, Farm1, Farm2, nil, nil},
		{nil, nil, Farm1, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeFarm,
}

// Mines

var Mine1 = &PlanUnits{
	Floors: []Floor{WoodFloor, WoodFloor},
	Roof:   ReedRoof,
}

var Mine2 = &PlanUnits{
	Floors: []Floor{StoneFloor, StoneFloor},
	Roof:   TileRoof,
}

var Mine1Plan = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Mine1, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeMine,
}

var Mine2Plan = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Mine2, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeMine,
}

// Workshops

// Factories

var Factory1 = &PlanUnits{
	Floors: []Floor{BrickFloor, BrickFloor},
	Roof:   BrickRoof,
}

var Factory1Plan = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Factory1, Factory1, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeFactory,
}

var Factory2 = &PlanUnits{
	Floors: []Floor{SandstoneFloor, SandstoneFloor},
	Roof:   SandstoneRoof,
}

var Factory2Plan = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Factory2, Factory2, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeFactory,
}

// Towers

var Tower1 = &PlanUnits{
	Floors: []Floor{StoneFloor, StoneFloor},
	Roof:   TileRoof,
}

var Tower1Plan = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Tower1, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeTower,
}

var Tower2 = &PlanUnits{
	Floors: []Floor{StoneFloor, StoneFloor, StoneFloor},
	Roof:   TileRoof,
}

var Tower2Plan = &BuildingPlan{
	BaseShape: [BuildingBaseMaxSize][BuildingBaseMaxSize]*PlanUnits{
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, Tower2, nil, nil},
		{nil, nil, nil, nil, nil},
		{nil, nil, nil, nil, nil},
	},
	BuildingType: BuildingTypeTower,
}

func DefaultPlans(bt BuildingType) []*BuildingPlan {
	switch bt {
	case BuildingTypeFarm:
		return []*BuildingPlan{Farm1Plan, Farm2Plan, Farm3Plan}
	case BuildingTypeMine:
		return []*BuildingPlan{Mine1Plan, Mine2Plan}
	case BuildingTypeFactory:
		return []*BuildingPlan{Factory1Plan, Factory2Plan}
	case BuildingTypeTower:
		return []*BuildingPlan{Tower1Plan, Tower2Plan}
	}
	return nil
}
