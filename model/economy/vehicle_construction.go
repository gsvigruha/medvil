package economy

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/vehicles"
)

type VehicleConstruction struct {
	Name                  string
	Time                  uint16
	Power                 uint16
	BuildingExtensionType *building.BuildingExtensionType
	Inputs                []artifacts.Artifacts
	Output                vehicles.VehicleType
}

var AllVehicleConstruction = [...]*VehicleConstruction{
	&VehicleConstruction{
		Name:                  "boat",
		Time:                  30 * 24,
		Power:                 1000,
		BuildingExtensionType: building.Deck,
		Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 3}},
		Output:                vehicles.Boat},
}

func GetVehicleConstructions(be *building.BuildingExtension) []*VehicleConstruction {
	result := make([]*VehicleConstruction, 0, len(AllVehicleConstruction))
	for _, m := range AllVehicleConstruction {
		if (be == nil && m.BuildingExtensionType == building.BuildingExtensionTypeNone) || (be != nil && be.T == m.BuildingExtensionType) {
			result = append(result, m)
		}
	}
	return result
}
