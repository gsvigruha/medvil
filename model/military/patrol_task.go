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

func (t *PatrolTask) Complete(m navigation.IMap, tool bool) bool {
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

func (t *PatrolTask) Tags() economy.Tags {
	return economy.EmptyTags
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

func (t *PatrolTask) Description() string {
	return "Patrol to protect your town from outlaws"
}
