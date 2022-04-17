package economy

import (
	"medvil/model/artifacts"
	"medvil/model/navigation"
	"medvil/model/time"
)

type ManufactureTask struct {
	TaskBase
	M        *Manufacture
	F        *navigation.Field
	R        *artifacts.Resources
	Progress uint16
}

func (t *ManufactureTask) Field() *navigation.Field {
	return t.F
}

func (t *ManufactureTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.Progress < t.M.Time {
		t.Progress++
		if tool {
			t.Progress++
		}
	}
	if t.Progress >= t.M.Time {
		t.R.AddAll(t.M.Outputs)
		return true
	}
	return false
}

func (t *ManufactureTask) Blocked() bool {
	return t.R.UsedVolumeCapacity() >= 1.0
}

func (t *ManufactureTask) Name() string {
	return t.M.Name
}

func (t *ManufactureTask) Tag() string {
	return ""
}

func (t *ManufactureTask) Expired(Calendar *time.CalendarType) bool {
	return false
}
