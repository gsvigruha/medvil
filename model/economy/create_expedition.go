package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

type CreateExpeditionTask struct {
	TaskBase
	Townhall Townhall
	PickupD  navigation.Destination
	Order    VehicleOrder
}

func (t *CreateExpeditionTask) Destination() navigation.Destination {
	return t.PickupD
}

func (t *CreateExpeditionTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	v := t.Order.PickupVehicle()
	t.Townhall.CreateExpedition(v, t.Person)
	return true
}

func (t *CreateExpeditionTask) Blocked() bool {
	return !t.Order.IsBuilt()
}

func (t *CreateExpeditionTask) Name() string {
	return "create_expedition"
}

func (t *CreateExpeditionTask) Tags() Tags {
	return MakeTags(FactoryPickupTaskTag(t.Order))
}

func (t *CreateExpeditionTask) Expired(Calendar *time.CalendarType) bool {
	return t.Order.IsExpired()
}

func (t *CreateExpeditionTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *CreateExpeditionTask) IsFieldCenter() bool {
	return false
}
