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
	state        int
}

func (t *PatrolTask) Destination() navigation.Destination {
<<<<<<< HEAD
	return t.Fields[t.state]
=======
	return t.Destinations[t.state]
>>>>>>> 74c77bd9cc4c12c693c54958e4b7ce9e2c5a5a54
}

func (t *PatrolTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.state < len(t.Destinations)-1 {
		t.state++
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

func (t *PatrolTask) Equipped(e economy.Equipment) bool {
	return e.Weapon()
}
