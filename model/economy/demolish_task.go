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
	Road     *building.Road
	F        *navigation.Field
	Town     ITown
	M        navigation.IMap
	Progress uint16
}

func (t *DemolishTask) Destination() navigation.Destination {
	return t.F
}

func (t *DemolishTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.Progress < DemolishTaskMaxProgress {
		t.Progress++
	} else {
		if t.Building != nil {
			t.Town.DestroyBuilding(t.Building, t.M)
		}
		if t.Road != nil {
			t.Town.DestroyRoad(t.Road, t.M)
		}
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
