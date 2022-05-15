package social

import (
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
	//"medvil/model/vehicles"
)

type VehicleOrder struct {
	T *economy.VehicleConstruction
	N int
}

type Factory struct {
	Household Household
	Orders    map[*economy.VehicleConstruction]*VehicleOrder
}

func (f *Factory) Init() {
	f.Orders = make(map[*economy.VehicleConstruction]*VehicleOrder)
	for _, vc := range economy.GetVehicleConstructions(f.Household.Building.Plan.GetExtension()) {
		if _, ok := f.Orders[vc]; !ok {
			f.Orders[vc] = &VehicleOrder{T: vc, N: 0}
		}
	}
}

func (f *Factory) ElapseTime(Calendar *time.CalendarType, m navigation.IMap) {
	f.Household.ElapseTime(Calendar, m)
	
}
