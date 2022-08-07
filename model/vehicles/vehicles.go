package vehicles

import (
	"medvil/model/navigation"
)

type VehicleType struct {
	Name  string
	Water bool
	Land  bool
}

var Boat = &VehicleType{Name: "boat", Water: true, Land: false}
var Cart = &VehicleType{Name: "cart", Water: false, Land: true}

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
}
