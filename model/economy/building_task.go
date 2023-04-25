package economy

import (
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/time"
	"strconv"
)

const BuildingTaskMaxProgress = 30 * 24

type BuildingTask struct {
	TaskBase
	D        navigation.Destination
	C        *building.Construction
	Started  bool
	Progress uint16
}

func (t *BuildingTask) Destination() navigation.Destination {
	return t.D
}

func (t *BuildingTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if !t.Started && !t.Blocked() {
		if len(t.C.Cost) > 0 {
			a := t.C.Storage.GetArtifacts()[0]
			t.C.Storage.Remove(a, 1)
		}
		t.Started = true
	}
	if t.Started {
		if t.Progress < BuildingTaskMaxProgress {
			t.Progress++
		} else {
			t.C.Progress++
			return true
		}
	}
	return false
}

func (t *BuildingTask) Blocked() bool {
	return !t.Started && t.C.Storage.IsEmpty() && len(t.C.Cost) > 0
}

func (t *BuildingTask) Name() string {
	return "building"
}

func (t *BuildingTask) Tag() string {
	return BuildingTaskTag(t.D)
}

func BuildingTaskTag(dest navigation.Destination) string {
	if f, ok := dest.(*navigation.Field); ok {
		return strconv.Itoa(int(f.X)) + "#" + strconv.Itoa(int(f.Y))
	}
	if l, ok := dest.(navigation.Location); ok {
		return strconv.Itoa(int(l.X)) + "#" + strconv.Itoa(int(l.Y)) + "#" + strconv.Itoa(int(l.Z))
	}
	return ""
}

func (t *BuildingTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *BuildingTask) Motion() uint8 {
	if t.C.Road != nil {
		return navigation.MotionFieldWork
	} else if t.C.Building != nil {
		return navigation.MotionBuild
	}
	return navigation.MotionStand
}
