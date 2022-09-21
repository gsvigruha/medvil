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
	Output                *vehicles.VehicleType
}

var BoatConstruction = &VehicleConstruction{
	Name:                  "boat",
	Time:                  30 * 24,
	Power:                 1000,
	BuildingExtensionType: building.Deck,
	Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 3}},
	Output:                vehicles.Boat,
}

var CartConstruction = &VehicleConstruction{
	Name:                  "cart",
	Time:                  30 * 24,
	Power:                 1000,
	BuildingExtensionType: nil,
	Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 3}},
	Output:                vehicles.Cart,
}

var TradingBoatConstruction = &VehicleConstruction{
	Name:                  "trading boat",
	Time:                  30 * 24,
	Power:                 1000,
	BuildingExtensionType: building.Deck,
	Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 3}},
	Output:                vehicles.TradingBoat,
}

var TradingCartConstruction = &VehicleConstruction{
	Name:                  "trading cart",
	Time:                  30 * 24,
	Power:                 1000,
	BuildingExtensionType: nil,
	Inputs:                []artifacts.Artifacts{artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 3}},
	Output:                vehicles.TradingCart,
}

var AllVehicleConstruction = [...]*VehicleConstruction{
	BoatConstruction,
	CartConstruction,
	TradingBoatConstruction,
	TradingCartConstruction,
}

func ConstructionCompatible(vc *VehicleConstruction, be *building.BuildingExtension) bool {
	return vc.BuildingExtensionType == nil || (be != nil && be.T == vc.BuildingExtensionType)
}

func GetVehicleConstructions(be *building.BuildingExtension) []*VehicleConstruction {
	result := make([]*VehicleConstruction, 0, len(AllVehicleConstruction))
	for _, m := range AllVehicleConstruction {
		if ConstructionCompatible(m, be) {
			result = append(result, m)
		}
	}
	return result
}
