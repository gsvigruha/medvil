package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type VehicleConstructionTask struct {
	TaskBase
	T        *VehicleConstruction
	F        *navigation.Field
	R        *artifacts.Resources
	Progress uint16
}

func (t *VehicleConstructionTask) Field() *navigation.Field {
	return t.F
}

func (t *VehicleConstructionTask) Complete(Calendar *time.CalendarType, tool bool) bool {
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
