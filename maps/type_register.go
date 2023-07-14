package maps

import (
	"encoding/json"
	"fmt"
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/economy"
	"medvil/model/materials"
	"medvil/model/navigation"
	"medvil/model/social"
	"medvil/model/terrain"
	"medvil/model/vehicles"
	"reflect"
)

func GetClassType(m json.RawMessage) reflect.Type {
	var mData map[string]json.RawMessage
	if err := json.Unmarshal(m, &mData); err != nil {
		fmt.Println("Error when loading struct: ", err)
	}
	var typeName string
	if err := json.Unmarshal(mData["$type"], &typeName); err != nil {
		fmt.Println("Error when loading struct: ", err)
	}
	switch typeName {
	case "Household":
		return reflect.TypeOf(social.Household{})
	case "Trader":
		return reflect.TypeOf(social.Trader{})
	case "Marketplace":
		return reflect.TypeOf(social.Marketplace{})
	case "BuildingUnit":
		return reflect.TypeOf(building.BuildingUnit{})
	case "RoofUnit":
		return reflect.TypeOf(building.RoofUnit{})
	case "ExtensionUnit":
		return reflect.TypeOf(building.ExtensionUnit{})
	case "EatTask":
		return reflect.TypeOf(economy.EatTask{})
	case "DrinkTask":
		return reflect.TypeOf(economy.DrinkTask{})
	case "HealTask":
		return reflect.TypeOf(economy.HealTask{})
	case "RelaxTask":
		return reflect.TypeOf(economy.RelaxTask{})
	case "GoHomeTask":
		return reflect.TypeOf(economy.GoHomeTask{})
	case "TransportTask":
		return reflect.TypeOf(economy.TransportTask{})
	case "ExchangeTask":
		return reflect.TypeOf(economy.ExchangeTask{})
	case "BuyTask":
		return reflect.TypeOf(economy.BuyTask{})
	case "SellTask":
		return reflect.TypeOf(economy.SellTask{})
	case "AgriculturalTask":
		return reflect.TypeOf(economy.AgriculturalTask{})
	case "BuildingTask":
		return reflect.TypeOf(economy.BuildingTask{})
	case "DemolishTask":
		return reflect.TypeOf(economy.DemolishTask{})
	case "ManufactureTask":
		return reflect.TypeOf(economy.ManufactureTask{})
	case "MiningTask":
		return reflect.TypeOf(economy.MiningTask{})
	case "FactoryPickupTask":
		return reflect.TypeOf(economy.FactoryPickupTask{})
	case "VehicleConstructionTask":
		return reflect.TypeOf(economy.VehicleConstructionTask{})
	case "TradeTask":
		return reflect.TypeOf(economy.TradeTask{})
	case "TerraformTask":
		return reflect.TypeOf(economy.TerraformTask{})
	case "Field":
		return reflect.TypeOf(navigation.Field{})
	case "Location":
		return reflect.TypeOf(navigation.Location{})
	case "BuildingDestination":
		return reflect.TypeOf(navigation.BuildingDestination{})
	}
	panic("Invalid type " + typeName)
}

func LoadStaticType(t reflect.Type, key string) reflect.Value {
	switch t.Elem().Name() {
	case "Artifact":
		return reflect.ValueOf(artifacts.GetArtifact(key))
	case "Material":
		return reflect.ValueOf(materials.GetMaterial(key))
	case "PlantType":
		return reflect.ValueOf(terrain.GetPlantType(key))
	case "TerrainType":
		return reflect.ValueOf(terrain.GetTerrainType(key))
	case "VehicleType":
		return reflect.ValueOf(vehicles.GetVehicleType(key))
	case "EquipmentType":
		return reflect.ValueOf(economy.GetEquipmentType(key))
	case "Manufacture":
		return reflect.ValueOf(economy.GetManufacture(key))
	case "RoadType":
		return reflect.ValueOf(building.GetRoadType(key))
	case "AnimalType":
		return reflect.ValueOf(terrain.GetAnimalType(key))
	case "BuildingExtensionType":
		return reflect.ValueOf(building.GetBuildingExtensionType(key))
	case "VehicleConstruction":
		return reflect.ValueOf(economy.GetVehicleConstruction(key))
	}
	panic("Invalid type " + t.Elem().Name())
}