package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
)

type Townhall interface {
	CreateTrader(v *vehicles.Vehicle, p Person)
	CreateExpedition(v *vehicles.Vehicle, p Person)
}

type CreateTraderTask struct {
	TaskBase
	Townhall Townhall
	PickupD  navigation.Destination
	Order    VehicleOrder
}

func (t *CreateTraderTask) Destination() navigation.Destination {
	return t.PickupD
}

func (t *CreateTraderTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	v := t.Order.PickupVehicle()
	t.Traveller.UseVehicle(v)
	t.Townhall.CreateTrader(v, t.Person)
	return true
}

func (t *CreateTraderTask) Blocked() bool {
	return !t.Order.IsBuilt()
}

func (t *CreateTraderTask) Name() string {
	return "create_trader"
}

func (t *CreateTraderTask) Tags() Tags {
	return MakeTags(FactoryPickupTaskTag(t.Order))
}

func (t *CreateTraderTask) Expired(Calendar *time.CalendarType) bool {
	return t.Order.IsExpired()
}

func (t *CreateTraderTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *CreateTraderTask) IsFieldCenter() bool {
	return false
}
