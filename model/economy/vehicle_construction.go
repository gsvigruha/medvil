package economy

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/vehicles"
)

type VehicleConstruction struct {
	Idx                   uint16
	Name                  string
	Time                  uint16
	Power                 uint16
	BuildingExtensionType *building.BuildingExtensionType
	Inputs                []artifacts.Artifacts
	Output                *vehicles.VehicleType
}

var BoatConstruction = &VehicleConstruction{
	Idx:                   1,
	Name:                  "boat",
	Time:                  30 * 24,
	Power:                 1000,
	BuildingExtensionType: building.Deck,
	Inputs: []artifacts.Artifacts{
		artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 3},
		artifacts.Artifacts{A: artifacts.GetArtifact("leather"), Quantity: 2},
	},
	Output: vehicles.Boat,
}

var CartConstruction = &VehicleConstruction{
	Idx:                   2,
	Name:                  "cart",
	Time:                  30 * 24,
	Power:                 1000,
	BuildingExtensionType: nil,
	Inputs: []artifacts.Artifacts{
		artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 2},
		artifacts.Artifacts{A: artifacts.GetArtifact("iron_bar"), Quantity: 1},
		artifacts.Artifacts{A: artifacts.GetArtifact("leather"), Quantity: 2},
	},
	Output: vehicles.Cart,
}

var TradingBoatConstruction = &VehicleConstruction{
	Idx:                   3,
	Name:                  "trading_boat",
	Time:                  30 * 24,
	Power:                 1000,
	BuildingExtensionType: building.Deck,
	Inputs: []artifacts.Artifacts{
		artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 3},
		artifacts.Artifacts{A: artifacts.GetArtifact("leather"), Quantity: 2},
		artifacts.Artifacts{A: artifacts.GetArtifact("textile"), Quantity: 2},
	},
	Output: vehicles.TradingBoat,
}

var TradingCartConstruction = &VehicleConstruction{
	Idx:                   4,
	Name:                  "trading_cart",
	Time:                  30 * 24,
	Power:                 1000,
	BuildingExtensionType: nil,
	Inputs: []artifacts.Artifacts{
		artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 2},
		artifacts.Artifacts{A: artifacts.GetArtifact("iron_bar"), Quantity: 1},
		artifacts.Artifacts{A: artifacts.GetArtifact("leather"), Quantity: 2},
		artifacts.Artifacts{A: artifacts.GetArtifact("textile"), Quantity: 2},
	},
	Output: vehicles.TradingCart,
}

var ExpeditionBoatConstruction = &VehicleConstruction{
	Idx:                   5,
	Name:                  "expedition_boat",
	Time:                  30 * 24 * 3,
	Power:                 1000,
	BuildingExtensionType: building.Deck,
	Inputs: []artifacts.Artifacts{
		artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 8},
		artifacts.Artifacts{A: artifacts.GetArtifact("leather"), Quantity: 5},
		artifacts.Artifacts{A: artifacts.GetArtifact("textile"), Quantity: 10},
	},
	Output: vehicles.ExpeditionBoat,
}

var ExpeditionCartConstruction = &VehicleConstruction{
	Idx:                   6,
	Name:                  "expedition_cart",
	Time:                  30 * 24 * 3,
	Power:                 1000,
	BuildingExtensionType: nil,
	Inputs: []artifacts.Artifacts{
		artifacts.Artifacts{A: artifacts.GetArtifact("board"), Quantity: 6},
		artifacts.Artifacts{A: artifacts.GetArtifact("iron_bar"), Quantity: 2},
		artifacts.Artifacts{A: artifacts.GetArtifact("leather"), Quantity: 5},
		artifacts.Artifacts{A: artifacts.GetArtifact("textile"), Quantity: 4},
	},
	Output: vehicles.ExpeditionCart,
}

var AllVehicleConstruction = [...]*VehicleConstruction{
	BoatConstruction,
	CartConstruction,
	TradingBoatConstruction,
	TradingCartConstruction,
	ExpeditionBoatConstruction,
	ExpeditionCartConstruction,
}

func GetVehicleConstruction(name string) *VehicleConstruction {
	for _, t := range AllVehicleConstruction {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func ConstructionCompatible(vc *VehicleConstruction, extensions []*building.BuildingExtension) bool {
	if vc.BuildingExtensionType == nil {
		return true
	} else {
		for _, e := range extensions {
			if e != nil && e.T == vc.BuildingExtensionType {
				return true
			}
		}
	}
	return false
}

func GetVehicleConstructions(extensions []*building.BuildingExtension) []*VehicleConstruction {
	result := make([]*VehicleConstruction, 0, len(AllVehicleConstruction))
	for _, m := range AllVehicleConstruction {
		if ConstructionCompatible(m, extensions) {
			result = append(result, m)
		}
	}
	return result
}
