package economy

import (
	"medvil/model/artifacts"
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/time"
)

const RepairTaskMaxProgress = 30 * 24

type RepairTask struct {
	TaskBase
	B        *building.Building
	F        *navigation.Field
	R        *artifacts.Resources
	Progress uint16
}

func (t *RepairTask) Destination() navigation.Destination {
	return t.F
}

func (t *RepairTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.Progress == 0 && !t.Blocked() {
		t.R.RemoveAll(t.B.Plan.RepairCost())
		t.Progress = 1
	}
	if t.Progress > 0 {
		if t.Progress < RepairTaskMaxProgress {
			t.Progress++
		} else {
			t.B.Broken = false
			return true
		}
	}
	return false
}

func (t *RepairTask) Blocked() bool {
	return t.R.Has(t.B.Plan.RepairCost())
}

func (t *RepairTask) Name() string {
	return "repair"
}

func (t *RepairTask) Tag() string {
	return ""
}

func (t *RepairTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *RepairTask) Motion() uint8 {
	return navigation.MotionBuild
}
