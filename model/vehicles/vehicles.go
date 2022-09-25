package vehicles

import (
	"medvil/model/navigation"
)

type VehicleType struct {
	Name            string
	Water           bool
	Land            bool
	IndoorStorage   bool
	Trader          bool
	MaxVolume       uint16
	BuildingCheckFn func(navigation.Field) bool
}

var Boat = &VehicleType{
	Name:            "boat",
	Water:           true,
	Land:            false,
	IndoorStorage:   false,
	Trader:          false,
	MaxVolume:       75,
	BuildingCheckFn: navigation.Field.Sailable,
}
var Cart = &VehicleType{
	Name:            "cart",
	Water:           false,
	Land:            true,
	IndoorStorage:   true,
	Trader:          false,
	MaxVolume:       50,
	BuildingCheckFn: navigation.Field.BuildingNonExtension,
}
var TradingBoat = &VehicleType{
	Name:            "trading_boat",
	Water:           true,
	Land:            false,
	IndoorStorage:   false,
	Trader:          true,
	MaxVolume:       75,
	BuildingCheckFn: navigation.Field.Sailable,
}
var TradingCart = &VehicleType{
	Name:            "trading_cart",
	Water:           false,
	Land:            true,
	IndoorStorage:   true,
	Trader:          true,
	MaxVolume:       50,
	BuildingCheckFn: navigation.Field.BuildingNonExtension,
}

type Vehicle struct {
	T         *VehicleType
	Traveller *navigation.Traveller
	InUse     bool
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
