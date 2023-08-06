package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
	"strconv"
)

type TransportTask struct {
	TaskBase
	PickupD          navigation.Destination
	DropoffD         navigation.Destination
	PickupR          *artifacts.Resources
	DropoffR         *artifacts.Resources
	A                *artifacts.Artifact
	Quantity         uint16
	CompleteQuantity bool
	Q                uint16
	Dropoff          bool
}

func (t *TransportTask) Destination() navigation.Destination {
	if t.Dropoff {
		return t.DropoffD
	} else {
		return t.PickupD
	}
}

func (t *TransportTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.Dropoff {
		t.DropoffR.Add(t.A, t.Q)
		if t.Quantity == 0 || !t.CompleteQuantity {
			return true
		}
		t.Dropoff = false
		return false
	} else {
		t.Q = t.PickupR.Remove(t.A, t.Quantity)
		t.Quantity -= t.Q
		t.Dropoff = true
	}
	return false
}

func (t *TransportTask) Blocked() bool {
	if t.DropoffR.UsedVolumeCapacity() > 1.0 {
		return true
	}
	if !t.Dropoff {
		return t.PickupR.Get(t.A) < t.Quantity
	}
	return false
}

func (t *TransportTask) Name() string {
	return "transport"
}

func (t *TransportTask) Tag() string {
	return TransportTaskTag(t.PickupD, t.A)
}

func TransportTaskTag(dest navigation.Destination, a *artifacts.Artifact) string {
	if f, ok := dest.(*navigation.Field); ok {
		return strconv.Itoa(int(f.X)) + "#" + strconv.Itoa(int(f.Y)) + "#" + a.Name
	}
	if l, ok := dest.(*navigation.Location); ok {
		return strconv.Itoa(int(l.X)) + "#" + strconv.Itoa(int(l.Y)) + "#" + strconv.Itoa(int(l.Z)) + "#" + a.Name
	}
	return a.Name
}

func (t *TransportTask) Expired(Calendar *time.CalendarType) bool {
	// If a water collection task is paused better to drop it and create a new one
	return t.A == Water && t.Paused
}

func (t *TransportTask) Motion() uint8 {
	return navigation.MotionStand
}
