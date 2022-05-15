package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type VehicleOrder interface {
	CompleteBuild()
}

type VehicleConstructionTask struct {
	TaskBase
	T        *VehicleConstruction
	O        VehicleOrder
	F        *navigation.Field
	R        *artifacts.Resources
	Progress uint16
}

func (t *VehicleConstructionTask) Field() *navigation.Field {
	return t.F
}

func (t *VehicleConstructionTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.Progress < t.T.Time {
		t.Progress++
		if tool {
			t.Progress++
		}
	}
	if t.Progress >= t.T.Time {
		t.O.CompleteBuild()
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

func (t *VehicleConstructionTask) Tag() string {
	return ""
}

func (t *VehicleConstructionTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *VehicleConstructionTask) Motion() uint8 {
	return navigation.MotionStand
}
