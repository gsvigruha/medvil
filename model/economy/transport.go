package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type TransportTask struct {
	PickupF  *navigation.Field
	DropoffF *navigation.Field
	PickupR  *artifacts.Resources
	DropoffR *artifacts.Resources
	A        *artifacts.Artifact
	Quantity uint16
	dropoff  bool
}

func (t *TransportTask) Field() *navigation.Field {
	if t.dropoff {
		return t.DropoffF
	} else {
		return t.PickupF
	}
}

func (t *TransportTask) Complete(Calendar *time.CalendarType) bool {
	if t.dropoff {
		t.DropoffR.Add(t.A, t.Quantity)
		return true
	} else {
		t.Quantity = t.PickupR.Remove(t.A, t.Quantity)
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
