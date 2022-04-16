package economy

import (
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/time"
)

const BuildingTaskMaxProgress = 30 * 24

type BuildingTask struct {
	F        *navigation.Field
	C        *building.Construction
	started  bool
	progress uint16
}

func (t *BuildingTask) Field() *navigation.Field {
	return t.F
}

func (t *BuildingTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if !t.started && !t.Blocked() {
		if len(t.C.Cost) > 0 {
			a := t.C.Storage.GetArtifacts()[0]
			t.C.Storage.Remove(a, 1)
		}
		t.started = true
	}
	if t.started {
		if t.progress < BuildingTaskMaxProgress {
			t.progress++
		} else {
			t.C.Progress++
			return true
		}
	}
	return false
}

func (t *BuildingTask) Blocked() bool {
	return !t.started && t.C.Storage.IsEmpty() && len(t.C.Cost) > 0
}

func (t *BuildingTask) Name() string {
	return "building"
}

func (t *BuildingTask) Tag() string {
	return BuildingTaskTag()
}

func BuildingTaskTag() string {
	return ""
}

func (t *BuildingTask) Expired(Calendar *time.CalendarType) bool {
	return false
}
