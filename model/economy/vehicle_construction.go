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
	Output                vehicles.Vehicle
}

var AllVehicleConstruction = [...]*VehicleConstruction{
	&VehicleConstruction{
		Name:   "boat",
		Time:   30 * 24,
		Power:  1000,
		Inputs: []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("log"), Quantity: 1}},
		Output: vehicles.Boat},
}
