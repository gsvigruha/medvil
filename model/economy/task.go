package economy

import (
	"medvil/model/navigation"
	"medvil/model/time"
)

type Task interface {
	Complete(Calendar *time.CalendarType, tool bool) bool
	Field() *navigation.Field
	Blocked() bool
	Name() string
	Tag() string
	Expired(Calendar *time.CalendarType) bool
	Pause(bool)
	IsPaused() bool
	Motion() uint8
}

type TaskBase struct {
	Paused bool
}

func (t *TaskBase) Pause(paused bool) {
	t.Paused = paused
}

func (t *TaskBase) IsPaused() bool {
	return t.Paused
}
