package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type TransportTask struct {
	PickupL  navigation.Location
	DropoffL navigation.Location
	PickupR  *artifacts.Resources
	DropoffR *artifacts.Resources
	A        *artifacts.Artifact
	Quantity uint16
	dropoff  bool
}

func (t *TransportTask) Location() navigation.Location {
	if t.dropoff {
		return t.DropoffL
	} else {
		return t.PickupL
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
	if t.dropoff {
		return false
	} else {
		return t.PickupR.Get(t.A) < t.Quantity
	}
}

func (t *TransportTask) Name() string {
	return "transport"
}
