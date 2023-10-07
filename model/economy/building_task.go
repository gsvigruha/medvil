package economy

import (
	"medvil/model/building"
	"medvil/model/navigation"
	"medvil/model/time"
)

const BuildingTaskMaxProgress = 30 * 24

type BuildingTask struct {
	TaskBase
	D        navigation.Destination
	C        *building.Construction
	Started  bool
	Progress uint16
	TaskTags *Tags
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

func (t *BuildingTask) Tags() Tags {
	if t.TaskTags == nil {
		tt := MakeTags(BuildingTaskTag(t.D))
		t.TaskTags = &tt
	}
	return *t.TaskTags
}

func BuildingTaskTag(dest navigation.Destination) Tag {
	if f, ok := dest.(*navigation.Field); ok {
		return SingleTag(f.X, f.Y)
	}
	if l, ok := dest.(*navigation.Location); ok {
		return SingleTag(l.X, l.Y, uint16(l.Z))
	}
	return EmptyTag
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
