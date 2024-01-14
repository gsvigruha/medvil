package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
	"medvil/model/vehicles"
)

type VehicleOrder interface {
	CompleteBuild(*navigation.Field)
	IsBuilt() bool
	IsExpired() bool
	Expire()
	PickupVehicle() *vehicles.Vehicle
	Name() string
	Idx() uint16
}

type VehicleConstructionTask struct {
	TaskBase
	T        *VehicleConstruction
	O        VehicleOrder
	F        *navigation.Field
	R        *artifacts.Resources
	Progress uint16
}

func (t *VehicleConstructionTask) Destination() navigation.Destination {
	return t.F
}

func (t *VehicleConstructionTask) Complete(m navigation.IMap, tool bool) bool {
	if t.Progress < t.T.Time {
		t.Progress++
		if tool {
			t.Progress++
		}
	}
	if t.Progress >= t.T.Time {
		t.O.CompleteBuild(t.F)
		return true
	}
	return false
}

func (t *VehicleConstructionTask) Blocked() bool {
	return false
}

func (t *VehicleConstructionTask) Name() string {
	return t.T.Name
}

func (t *VehicleConstructionTask) Tags() Tags {
	return EmptyTags
}

func (t *VehicleConstructionTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *VehicleConstructionTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *VehicleConstructionTask) Description() string {
	return "Construct vehicles"
}
