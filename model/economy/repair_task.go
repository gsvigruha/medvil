package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

const RepairTaskMaxProgress = 30 * 24

type Repairable interface {
	Repair()
	RepairCost() []artifacts.Artifacts
	Name() string
}

type RepairTask struct {
	TaskBase
	Repairable Repairable
	Field      *navigation.Field
	Resources  *artifacts.Resources
	Progress   uint16
	TaskTags   *Tags
}

func (t *RepairTask) Destination() navigation.Destination {
	return t.Field
}

func (t *RepairTask) Complete(m navigation.IMap, tool bool) bool {
	if t.Progress == 0 && !t.Blocked() {
		t.Resources.RemoveAll(t.Repairable.RepairCost())
		t.Progress = 1
	}
	if t.Progress > 0 {
		if t.Progress < RepairTaskMaxProgress {
			t.Progress++
		} else {
			t.Repairable.Repair()
			return true
		}
	}
	return false
}

func (t *RepairTask) Blocked() bool {
	return !t.Resources.HasAll(t.Repairable.RepairCost())
}

func (t *RepairTask) Name() string {
	return "repair"
}

func (t *RepairTask) Tags() Tags {
	if t.TaskTags == nil {
		tt := MakeTags(BuildingTaskTag(t.Field))
		t.TaskTags = &tt
	}
	return *t.TaskTags
}

func (t *RepairTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *RepairTask) Motion() uint8 {
	return navigation.MotionBuild
}

func (t *RepairTask) Description() string {
	return "Repair " + t.Repairable.Name()
}
