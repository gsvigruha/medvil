package vehicles

type VehicleType struct {
	Name  string
	Water bool
	Land  bool
}

var Boat = VehicleType{Name: "boat", Water: true, Land: false}
