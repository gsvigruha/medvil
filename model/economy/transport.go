package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type TransportTask struct {
	TaskBase
	PickupD          navigation.Destination
	DropoffD         navigation.Destination
	PickupR          *artifacts.Resources
	DropoffR         *artifacts.Resources
	A                *artifacts.Artifact
	TargetQuantity   uint16
	CompleteQuantity bool
	ActualQuantity   uint16
	Dropoff          bool
	TaskTags         *Tags
}

func (t *TransportTask) Destination() navigation.Destination {
	if t.Dropoff {
		return t.DropoffD
	} else {
		return t.PickupD
	}
}

func (t *TransportTask) Complete(m navigation.IMap, tool bool) bool {
	if t.Dropoff {
		t.DropoffR.Add(t.A, t.ActualQuantity)
		if t.TargetQuantity == 0 || !t.CompleteQuantity {
			return true
		}
		t.Dropoff = false
		return false
	} else {
		t.ActualQuantity = t.PickupR.Remove(t.A, t.TargetQuantity)
		t.TargetQuantity -= t.ActualQuantity
		t.Dropoff = true
	}
	return false
}

func (t *TransportTask) Blocked() bool {
	if t.DropoffR.UsedVolumeCapacity() > 1.0 {
		return true
	}
	if !t.Dropoff {
		return t.PickupR.Get(t.A) < t.TargetQuantity
	}
	return false
}

func (t *TransportTask) Name() string {
	return "transport"
}

func (t *TransportTask) Tags() Tags {
	if t.TaskTags == nil {
		tt := MakeTags(TransportTaskTag(t.PickupD, t.A))
		t.TaskTags = &tt
	}
	return *t.TaskTags
}

func TransportTaskTag(dest navigation.Destination, a *artifacts.Artifact) Tag {
	if f, ok := dest.(*navigation.Field); ok {
		return SingleTag(a.Idx, TagField, f.X, f.Y)
	}
	if l, ok := dest.(*navigation.Location); ok {
		return SingleTag(a.Idx, TagLocation, l.X, l.Y, uint16(l.Z))
	}
	if _, ok := dest.(*navigation.TravellerDestination); ok {
		return SingleTag(a.Idx, TagTraveller)
	}
	if b, ok := dest.(*navigation.BuildingDestination); ok {
		return SingleTag(a.Idx, TagBuilding, b.B.X, b.B.Y)
	}
	return SingleTag(a.Idx)
}

func (t *TransportTask) Expired(Calendar *time.CalendarType) bool {
	// If a water collection task is paused better to drop it and create a new one
	if t.A == Water && t.Paused {
		return true
	}
	if !t.CompleteQuantity && t.Paused {
		return true
	}
	if t.PickupR.Deleted || t.DropoffR.Deleted {
		return true
	}
	return false
}

func (t *TransportTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *TransportTask) Description() string {
	return "Transport goods"
}
