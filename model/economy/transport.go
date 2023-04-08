package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
	"strconv"
)

type TransportTask struct {
	TaskBase
	PickupF          *navigation.Field
	DropoffF         *navigation.Field
	PickupR          *artifacts.Resources
	DropoffR         *artifacts.Resources
	A                *artifacts.Artifact
	Quantity         uint16
	CompleteQuantity bool
	q                uint16
	dropoff          bool
}

func (t *TransportTask) Destination() navigation.Destination {
	if t.dropoff {
		return t.DropoffF
	} else {
		return t.PickupF
	}
}

func (t *TransportTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.dropoff {
		t.DropoffR.Add(t.A, t.q)
		return t.Quantity == 0 || !t.CompleteQuantity
	} else {
		t.q = t.PickupR.Remove(t.A, t.Quantity)
		t.Quantity -= t.q
		t.dropoff = true
	}
	return false
}

func (t *TransportTask) Blocked() bool {
	if t.DropoffR.UsedVolumeCapacity() > 1.0 {
		return true
	}
	if !t.dropoff {
		return t.PickupR.Get(t.A) < t.Quantity
	}
	return false
}

func (t *TransportTask) Name() string {
	return "transport"
}

func (t *TransportTask) Tag() string {
	return TransportTaskTag(t.PickupF, t.A)
}

func TransportTaskTag(f *navigation.Field, a *artifacts.Artifact) string {
	return strconv.Itoa(int(f.X)) + "#" + strconv.Itoa(int(f.Y)) + "#" + a.Name
}

func (t *TransportTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *TransportTask) Motion() uint8 {
	return navigation.MotionStand
}
