package military

import (
	"medvil/model/economy"
	"medvil/model/navigation"
	"medvil/model/time"
)

type PatrolTask struct {
	economy.TaskBase
	Fields []*navigation.Field
	state  int
}

func (t *PatrolTask) Field() *navigation.Field {
	return t.Fields[t.state]
}

func (t *PatrolTask) Complete(Calendar *time.CalendarType, tool bool) bool {
	if t.state < len(t.Fields) {
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
	return false
}

func (t *PatrolTask) Motion() uint8 {
	return navigation.MotionStand
}
