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
	Progress uint16
}

func (t *DemolishTask) Destination() navigation.Destination {
	return t.F
}

func (t *DemolishTask) Complete(m navigation.IMap, tool bool) bool {
	if t.Progress < DemolishTaskMaxProgress {
		t.Progress++
	} else {
		if t.Building != nil {
			t.Town.DestroyBuilding(t.Building, m)
		}
		if t.Road != nil {
			t.Town.DestroyRoad(t.Road, m)
		}
		t.F.Allocated = false
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

func (t *DemolishTask) Tags() Tags {
	return MakeTags(DemolishTaskTag(t.F))
}

func DemolishTaskTag(f *navigation.Field) Tag {
	return SingleTag(f.X, f.Y)
}

func (t *DemolishTask) Expired(Calendar *time.CalendarType) bool {
	return false
}

func (t *DemolishTask) Motion() uint8 {
	return navigation.MotionBuild
}

func (t *DemolishTask) Description() string {
	if t.Building != nil {
		return "Demolish buildings"
	} else if t.Road != nil {
		return "Demolish roads"
	}
	return ""
}
