package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
)

type Household interface {
	AddVehicle(*vehicles.Vehicle)
	AllocateVehicle(waterOk bool) *vehicles.Vehicle
	GetExchange() Exchange
}

type FactoryPickupTask struct {
	TaskBase
	PickupD  navigation.Destination
	DropoffD navigation.Destination
	Order    VehicleOrder
	Dropoff  bool
}

func (t *FactoryPickupTask) Destination() navigation.Destination {
	if t.Dropoff {
		return t.DropoffD
	} else {
		return t.PickupD
	}
}

func (t *FactoryPickupTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.Dropoff {
		t.Traveller.ExitVehicle()
		return true
	} else {
		v := t.Order.PickupVehicle()
		t.Household.AddVehicle(v)
		t.Traveller.UseVehicle(v)
		if f, ok := t.DropoffD.(*navigation.Field); ok {
			t.Traveller.Vehicle.SetParking(f)
		}
		t.Dropoff = true
	}
	return false
}

func (t *FactoryPickupTask) Blocked() bool {
	if !t.Dropoff {
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
	if t.Paused {
		t.Order.Expire()
	}
	return t.Order.IsExpired()
}

func (t *FactoryPickupTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *FactoryPickupTask) IsFieldCenter() bool {
	if t.Dropoff {
		return t.FieldCenter
	} else {
		return false
	}
}
