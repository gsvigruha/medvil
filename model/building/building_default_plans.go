package building

import (
	"medvil/model/materials"
)

var WhitewashFloor = Floor{M: materials.GetMaterial("whitewash")}
var ReedRoof = &Roof{M: materials.GetMaterial("reed"), RoofType: RoofTypeSplit}

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

func DefaultPlans(bt BuildingType) []*BuildingPlan {
	switch bt {
	case BuildingTypeFarm:
		return []*BuildingPlan{Farm1Plan, Farm2Plan, Farm3Plan}
	}
	return nil
}
