package military

import (
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

type PatrolTask struct {
	economy.TaskBase
	Destinations []navigation.Destination
	Start        time.CalendarType
	State        int
}

func (t *PatrolTask) Destination() navigation.Destination {
	return t.Destinations[t.State]
}

func (t *PatrolTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.State < len(t.Destinations)-1 {
		t.State++
		return false
	} else {
		return true
	}
}

func (t *PatrolTask) Blocked() bool {
	return false
}

func (t *PatrolTask) Name() string {
	return "patrol"
}

func (t *PatrolTask) Tag() string {
	return ""
}

func (t *PatrolTask) Expired(Calendar *time.CalendarType) bool {
	return Calendar.DaysElapsed()-t.Start.DaysElapsed() >= 30
}

func (t *PatrolTask) Motion() uint8 {
	return navigation.MotionStand
}

func (t *PatrolTask) Equipped(e *economy.EquipmentType) bool {
	return e.Weapon
}
