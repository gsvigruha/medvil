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

type Vehicle struct {
	T         *VehicleType
	Traveller *navigation.Traveller
}
