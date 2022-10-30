package economy

import (
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/time"
)

const DemolishTaskMaxProgress = 30 * 24

type DemolishTask struct {
	TaskBase
	Building *building.Building
	F        *navigation.Field
	progress uint16
}

func (t *DemolishTask) Field() *navigation.Field {
	return t.F
}

func (t *DemolishTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.progress < BuildingTaskMaxProgress {
		t.progress++
	} else {
		return true
	}
	return false
}

func (t *DemolishTask) Blocked() bool {
	return false
}

func (t *DemolishTask) Name() string {
	return "demolish"
}

func (t *DemolishTask) Tag() string {
	return ""
}

func (t *DemolishTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *DemolishTask) Motion() uint8 {
	return navigation.MotionBuild
}
