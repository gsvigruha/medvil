package vehicles

import (
	"medvil/model/building"
	"medvil/model/navigation"
)

type VehicleType struct {
	Name                  string
	Water                 bool
	Land                  bool
	IndoorStorage         bool
	Trader                bool
	MaxVolume             uint16
	BuildingCheckFn       func(navigation.Field) bool
	BuildingExtensionType *building.BuildingExtensionType
}

var Boat = &VehicleType{
	Name:                  "boat",
	Water:                 true,
	Land:                  false,
	IndoorStorage:         false,
	Trader:                false,
	MaxVolume:             75,
	BuildingCheckFn:       navigation.Field.Sailable,
	BuildingExtensionType: building.Deck,
}
var Cart = &VehicleType{
	Name:                  "cart",
	Water:                 false,
	Land:                  true,
	IndoorStorage:         true,
	Trader:                false,
	MaxVolume:             50,
	BuildingCheckFn:       navigation.Field.BuildingNonExtension,
	BuildingExtensionType: building.NonExtension,
}
var TradingBoat = &VehicleType{
	Name:                  "trading_boat",
	Water:                 true,
	Land:                  false,
	IndoorStorage:         false,
	Trader:                true,
	MaxVolume:             75,
	BuildingCheckFn:       navigation.Field.Sailable,
	BuildingExtensionType: building.Deck,
}
var TradingCart = &VehicleType{
	Name:                  "trading_cart",
	Water:                 false,
	Land:                  true,
	IndoorStorage:         true,
	Trader:                true,
	MaxVolume:             50,
	BuildingCheckFn:       navigation.Field.BuildingNonExtension,
	BuildingExtensionType: building.NonExtension,
}

var VehicleTypes = [...]*VehicleType{
	Boat,
	Cart,
	TradingBoat,
	TradingCart,
}

func GetVehicleType(name string) *VehicleType {
	for _, t := range VehicleTypes {
		if t.Name == name {
			return t
		}
	}
	return nil
}

type Vehicle struct {
	T         *VehicleType
	Traveller *navigation.Traveller
	InUse     bool
	Broken    bool
}

func (v *Vehicle) PathType() navigation.PathType {
	if v.T.Water {
		return navigation.PathTypeBoat
	}
	return navigation.PathTypePedestrian
}

func (v *Vehicle) GetTraveller() *navigation.Traveller {
	return v.Traveller
}

func (v *Vehicle) SetInUse(inUse bool) {
	v.InUse = inUse
	if v.Traveller != nil && v.T.IndoorStorage {
		v.Traveller.Visible = inUse
	}
}

func (v *Vehicle) SetHome(home bool) {
	if v.Traveller != nil && v.T.IndoorStorage {
		v.Traveller.Visible = !home
	}
}

func (v *Vehicle) Water() bool {
	return v.T.Water
}

func (v *Vehicle) Break() {
	v.Broken = true
}
