package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
)

type Household interface {
	AddVehicle(*vehicles.Vehicle)
	GetVehicle() *vehicles.Vehicle
}

type FactoryPickupTask struct {
	TaskBase
	PickupD  navigation.Destination
	DropoffD navigation.Destination
	Order    VehicleOrder
	dropoff  bool
}

func (t *FactoryPickupTask) Destination() navigation.Destination {
	if t.dropoff {
		return t.DropoffD
	} else {
		return t.PickupD
	}
}

func (t *FactoryPickupTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.dropoff {
		t.Traveller.ExitVehicle()
		return true
	} else {
		v := t.Order.PickupVehicle()
		t.Household.AddVehicle(v)
		t.Traveller.UseVehicle(v)
		t.dropoff = true
	}
	return false
}

func (t *FactoryPickupTask) Blocked() bool {
	if !t.dropoff {
		return !t.Order.IsBuilt()
	}
	return false
}

func (t *FactoryPickupTask) Name() string {
	return "factory_pickup"
}

func (t *FactoryPickupTask) Tag() string {
	return FactoryPickupTaskTag(t.Order)
}

func FactoryPickupTaskTag(o VehicleOrder) string {
	return o.Name()
}

func (t *FactoryPickupTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *FactoryPickupTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *FactoryPickupTask) IsFieldCenter() bool {
	if t.dropoff {
		return t.FieldCenter
	} else {
		return false
	}
}
