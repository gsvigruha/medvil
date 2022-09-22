package vehicles

import (
	"medvil/model/navigation"
)

type VehicleType struct {
	Name            string
	Water           bool
	Land            bool
	IndoorStorage   bool
	MaxVolume       uint16
	BuildingCheckFn func(navigation.Field) bool
}

var Boat = &VehicleType{Name: "boat", Water: true, Land: false, IndoorStorage: false, MaxVolume: 75, BuildingCheckFn: navigation.Field.Sailable}
var Cart = &VehicleType{Name: "cart", Water: false, Land: true, IndoorStorage: true, MaxVolume: 50, BuildingCheckFn: navigation.Field.BuildingNonExtension}
var TradingBoat = &VehicleType{Name: "trading_boat", Water: true, Land: false, IndoorStorage: false, MaxVolume: 75, BuildingCheckFn: navigation.Field.Sailable}
var TradingCart = &VehicleType{Name: "trading_cart", Water: false, Land: true, IndoorStorage: true, MaxVolume: 50, BuildingCheckFn: navigation.Field.BuildingNonExtension}

type Vehicle struct {
	T         *VehicleType
	Traveller *navigation.Traveller
	InUse     bool
}

func (v *Vehicle) TravellerType() uint8 {
	if v.T == Boat {
		return navigation.TravellerTypeBoat
	}
	return navigation.TravellerTypePedestrian
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
